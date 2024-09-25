package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"admin-cli/internal"

	"github.com/spf13/cobra"
)

var (
	dirPath        string
	maxConcurrency int
)

var hashCheckCmd = &cobra.Command{
	Use:   "hashcheck",
	Short: "Compute and check file hashes for collisions in a directory using concurrency",
	Run: func(cmd *cobra.Command, args []string) {
		hashes := make(map[[32]byte][]string) // Map of hashes (SHA-256 produces a 32-byte hash) to file paths
		mu := &sync.Mutex{}                   // Mutex to protect shared state

		// Buffered channel to send file paths to workers
		fileChan := make(chan string, maxConcurrency)

		// WaitGroup to track when all workers are done
		var wg sync.WaitGroup

		// Start workers based on maxConcurrency
		for i := 0; i < maxConcurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for path := range fileChan {
					hash, err := internal.ComputeFileHash(path)
					if err != nil {
						fmt.Printf("Error computing hash for file %s: %v\n", path, err)
						continue
					}

					// Store the file path in the map by its hash
					mu.Lock() // Protect the map with a mutex since it is shared across go routines
					hashes[hash] = append(hashes[hash], path)
					mu.Unlock()
				}
			}()
		}

		// Walk the directory and send file paths to the workers
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error accessing file %s: %v", path, err)
			}
			// Ignore directories
			if info.IsDir() {
				return nil
			}
			// Send file path to the channel for processing
			fileChan <- path
			return nil
		})

		// Close the channel when all file paths have been sent
		close(fileChan)

		// Wait for all workers to finish
		wg.Wait()

		if err != nil {
			fmt.Println("Error walking the directory:", err)
			return
		}

		// Check for collisions (files with the same hash)
		internal.PrintCollisions(hashes)
	},
}

func init() {
	// Dynamically set default concurrency based on available CPU cores
	defaultConcurrency := runtime.NumCPU() // Set default to number of CPUs
	hashCheckCmd.Flags().StringVarP(&dirPath, "dir", "d", ".", "Directory to scan for file hash collisions")
	hashCheckCmd.Flags().IntVarP(&maxConcurrency, "routines", "r", defaultConcurrency, "Number of concurrent workers to process files")
	rootCmd.AddCommand(hashCheckCmd)
}
