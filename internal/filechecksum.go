package internal

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// ComputeFileHash computes the SHA-256 hash of a file at the given path
func ComputeFileHash(path string) ([32]byte, error) {
	file, err := os.Open(path) // Open the file at the given path
	if err != nil {
		return [32]byte{}, err // Return an empty array and the error if opening fails
	}
	defer file.Close() // Ensure the file is closed when the function exits

	hash := sha256.New()                           // Create a new SHA-256 hash instance
	if _, err := io.Copy(hash, file); err != nil { // Copy file content into the hash
		return [32]byte{}, err // Return an empty array and the error if copying fails
	}

	var result [32]byte            // Declare a variable to hold the resulting hash
	copy(result[:], hash.Sum(nil)) // Copy the computed hash into the result variable

	return result, nil // Return the hash and no error
}

// PrintCollisions prints the files that have the same hash
func PrintCollisions(hashes map[[32]byte][]string) {
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
