package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var dirPath string

var hashCheckCmd = &cobra.Command{
	Use:   "hashcheck",
	Short: "Compute and check file hashes for collisions in a directory",
	Run: func(cmd *cobra.Command, args []string) {
		// Leverage a map to store the hashes of the files in the directory
		hashes := make(map[[32]byte][]string) // Map of hashes (SHA-256 produces a 32-byte hash) to file paths
		// Leverage filepath.Walk to traverse the directory and compute the hash of each file, storing it in the map
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Ignore directories
			if info.IsDir() {
				return nil
			}

			// Compute the hash of the file
			hash, err := computeFileHash(path)
			if err != nil {
				return fmt.Errorf("error computing hash for file %s: %v", path, err)
			}

			// Store the file path in the map by its hash
			hashes[hash] = append(hashes[hash], path)

			return nil
		})

		if err != nil {
			fmt.Println("Error walking the directory:", err)
			return
		}

		// Check for collisions (files with the same hash)
		printCollisions(hashes)
	},
}

func init() {
	hashCheckCmd.Flags().StringVarP(&dirPath, "dir", "d", ".", "Directory to scan for file hash collisions")
	rootCmd.AddCommand(hashCheckCmd)
}

// computeFileHash computes the SHA-256 hash of a file at the given path
func computeFileHash(path string) ([32]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return [32]byte{}, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return [32]byte{}, err
	}

	return sha256.Sum256(hash.Sum(nil)), nil
}

// printCollisions prints the files that have the same hash
func printCollisions(hashes map[[32]byte][]string) {
	hasCollisions := false // Default to no collisions

	for hash, files := range hashes {
		// If a hash has more than one file, it means there is a collision.
		if len(files) > 1 {
			hasCollisions = true // Update outer variable
			fmt.Printf("Hash collision detected for hash %x:\n", hash)
			for _, file := range files {
				fmt.Println(" -", file)
			}
			fmt.Println() // Add a newline for readability
		}
	}

	if !hasCollisions {
		fmt.Println("No hash collisions found.")
	}
}
