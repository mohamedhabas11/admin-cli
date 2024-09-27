package internal

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

// ScanPorts scans a range of ports on a given IP address
// and returns a slice of open ports with optional progress bar feedback.
func ScanPorts(ip string, ports []int, showProgress bool) []int {
	var openPorts []int
	var mu sync.Mutex // Mutex to protect shared resources
	var wg sync.WaitGroup
	results := make(chan int, len(ports)) // Buffered channel to avoid blocking
	semaphore := make(chan struct{}, 100) // Limit to 100 concurrent scans
	var bar *progressbar.ProgressBar

	// Initialize progress bar if needed
	if showProgress {
		bar = progressbar.Default(int64(len(ports))) // Initialize progress bar
	}

	// Function to scan a single port
	scan := func(port int) {
		defer wg.Done()
		target := fmt.Sprintf("%s:%d", ip, port)

		// Add a timeout for connection attempts
		conn, err := net.DialTimeout("tcp", target, 1*time.Second)
		if err == nil {
			// Port is open
			conn.Close()
			mu.Lock()       // Lock before appending to shared resource
			results <- port // Send the open port to the results channel
			mu.Unlock()
		}
		if showProgress {
			bar.Add(1) // Update the progress bar
		}
	}

	// Start goroutines for each port
	for _, port := range ports {
		wg.Add(1)
		semaphore <- struct{}{} // Limit concurrency
		go func(p int) {
			defer func() { <-semaphore }()
			scan(p)
		}(port)
	}

	go func() {
		// Close results channel once all goroutines are done
		wg.Wait()
		close(results)
	}()

	// Collect results and display open ports immediately
	for port := range results {
		mu.Lock() // Ensure safe access to openPorts
		openPorts = append(openPorts, port)
		mu.Unlock()

		// Print open port immediately
		fmt.Printf("\r\033[KPort %d is open\n", port) // Clear the line before printing
	}

	if showProgress {
		bar.Finish() // Finish the progress bar
	}

	// Sort the open ports for consistent output
	sort.Ints(openPorts)

	return openPorts
}
