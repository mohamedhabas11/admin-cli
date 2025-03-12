package internal

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"filippo.io/age"
	"github.com/klauspost/compress/zstd"
)

// Restore decompresses and extracts the tar.zst archive from srcFile into destDir.
// If a passphrase is provided, it decrypts the archive before decompressing.
func Restore(srcFile, destDir, passphrase string) error {
	in, err := openSourceFile(srcFile)
	if err != nil {
		return err
	}
	defer in.Close()

	reader, err := setupReader(in, passphrase)
	if err != nil {
		return err
	}

	decoder, err := setupDecompressor(reader)
	if err != nil {
		return err
	}
	defer decoder.Close()

	tr := setupTarReader(decoder)
	return extractTar(tr, destDir)
}

// openSourceFile opens the source backup file for reading.
func openSourceFile(srcFile string) (*os.File, error) {
	return os.Open(srcFile)
}

// setupReader returns the appropriate reader, decrypting if a passphrase is provided.
func setupReader(r io.Reader, passphrase string) (io.Reader, error) {
	if passphrase == "" {
		return r, nil
	}
	return decryptReader(r, passphrase)
}

// decryptReader decrypts the reader using the provided passphrase.
func decryptReader(r io.Reader, passphrase string) (io.Reader, error) {
	identity, err := age.NewScryptIdentity(passphrase)
	if err != nil {
		return nil, err
	}
	return age.Decrypt(r, identity)
}

// setupDecompressor initializes a zstd decompressor from the reader.
func setupDecompressor(r io.Reader) (*zstd.Decoder, error) {
	return zstd.NewReader(r)
}

// setupTarReader creates a tar reader from the reader.
func setupTarReader(r io.Reader) *tar.Reader {
	return tar.NewReader(r)
}

// extractTar processes the tar archive and extracts its contents.
func extractTar(tr *tar.Reader, destDir string) error {
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil // End of archive
		}
		if err != nil {
			return err
		}
		if err := restoreItem(tr, header, destDir); err != nil {
			return err
		}
	}
}

// restoreItem restores a single tar entry (file, directory, or symlink).
func restoreItem(tr *tar.Reader, header *tar.Header, destDir string) error {
	targetPath := filepath.Join(destDir, header.Name)
	switch header.Typeflag {
	case tar.TypeDir:
		return createDirectory(targetPath, header.Mode)
	case tar.TypeReg:
		return restoreFile(tr, targetPath, header.Mode)
	case tar.TypeSymlink:
		return createSymlink(header.Linkname, targetPath)
	default:
		return nil // Ignore unsupported types
	}
}

// createDirectory creates a directory with the specified mode.
func createDirectory(path string, mode int64) error {
	return os.MkdirAll(path, os.FileMode(mode))
}

// restoreFile extracts a regular file from the tar reader.
func restoreFile(tr *tar.Reader, path string, mode int64) error {
	if err := ensureParentDir(path); err != nil {
		return err
	}
	outFile, err := createFile(path, mode)
	if err != nil {
		return err
	}
	defer outFile.Close()
	return copyFileContents(outFile, tr)
}

// ensureParentDir creates the parent directory of the given path.
func ensureParentDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}

// createFile creates a file with the specified mode.
func createFile(path string, mode int64) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.FileMode(mode))
}

// copyFileContents copies data from the tar reader to the output file.
func copyFileContents(dst *os.File, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}

// createSymlink creates a symbolic link.
func createSymlink(target, path string) error {
	return os.Symlink(target, path)
}
