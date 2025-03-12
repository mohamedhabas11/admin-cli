package internal

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"filippo.io/age"
	"github.com/klauspost/compress/zstd"
)

// Backup creates a compressed and optionally encrypted tar archive of the source path.
func Backup(srcPath, destFile string, compLevel int, followSymlinks bool, passphrase string) error {

	// Create the destination file
	out, err := CreateDestinationFile(destFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Setup the writer for the compressor based on encryption
	var compressorWriter io.Writer
	if passphrase != "" {
		// chain compressorWriter to the encoder io.Writer to create encrypted destination file
		encryptor, err := setupEncryptor(out, passphrase)
		if err != nil {
			return err
		}
		defer encryptor.Close()
		compressorWriter = encryptor
	} else {
		// output encoder to create destination file
		compressorWriter = out
	}

	// Setup the zstd compressor
	encoder, err := setupCompressor(compressorWriter, compLevel)
	if err != nil {
		return err
	}
	defer encoder.Close()

	// Setup the tar writer
	tw := setupTarWriter(encoder)
	defer tw.Close()

	// Decide which stat function to use
	statFunc := getStatFunc(followSymlinks)

	// Get file info for the source path
	fi, err := statFunc(srcPath)
	if err != nil {
		return err
	}

	// Recursively walk directories and archive files
	if fi.IsDir() {
		return filepath.Walk(srcPath, func(file string, _ os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			return archiveFile(tw, file, statFunc, srcPath)
		})
	}

	return archiveFile(tw, srcPath, statFunc, filepath.Dir(srcPath))
}

// CreateDestinationFile creates the destination file for the backup.
func CreateDestinationFile(destFile string) (*os.File, error) {
	return os.Create(destFile)
}

// setupCompressor creates a zstd compressor that writes to the given writer.
func setupCompressor(w io.Writer, compLevel int) (*zstd.Encoder, error) {
	return zstd.NewWriter(w, zstd.WithEncoderLevel(zstd.EncoderLevel(compLevel)))
}

// setupTarWriter creates a tar writer that writes to the given writer.
func setupTarWriter(w io.Writer) *tar.Writer {
	return tar.NewWriter(w)
}

// getStatFunc returns the apporiate stat function on whether to follow symlinks.
func getStatFunc(followSymlinks bool) func(name string) (os.FileInfo, error) {
	if followSymlinks {
		return os.Stat
	}
	return os.Lstat
}

// archiveFile adds a file to the tar archive.
func archiveFile(tw *tar.Writer, file string, statFunc func(name string) (os.FileInfo, error), basePath string) error {
	// Get file info to store in the tar header Metadata(permissions, filetype, filename, relative path)
	fileInfo, err := statFunc(file)
	if err != nil {
		return err
	}

	// Create tar header
	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return err
	}

	// Compute relative path for the header name
	relPath, err := filepath.Rel(basePath, file)
	if err != nil {
		return err
	}
	header.Name = relPath

	// Write the header to the tar archive
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	// Skip copying content from non-regular files (e.g, directories, symlinks)
	if !fileInfo.Mode().IsRegular() {
		return nil
	}

	// Open the file and copy its contents to the tar writer
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(tw, f)
	return err
}

// setupEncryptor creates an Age encryptor with a given passphrase.
func setupEncryptor(w io.Writer, passphrase string) (io.WriteCloser, error) {
	recipient, err := age.NewScryptRecipient(passphrase)
	if err != nil {
		return nil, err
	}
	return age.Encrypt(w, recipient)
}
