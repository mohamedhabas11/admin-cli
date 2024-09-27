package internal

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

// ScanPorts scans a range of ports on a given IP address
// and returns a slice of open ports with progress bar feedback.
func ScanPorts(ip string, ports []int) []int {
	var openPorts []int
	var wg sync.WaitGroup
	results := make(chan int, len(ports))         // Buffered channel to avoid blocking
	semaphore := make(chan struct{}, 100)         // Limit to 100 concurrent scans
	bar := progressbar.Default(int64(len(ports))) // Initialize progress bar

	// Function to scan a single port
	scan := func(port int) {
		defer wg.Done()
		target := fmt.Sprintf("%s:%d", ip, port)

		// Add a timeout for connection attempts
		conn, err := net.DialTimeout("tcp", target, 1*time.Second)
		if err == nil {
			// Port is open
			conn.Close()
			results <- port
		}
		// Update the progress bar after each scan
		bar.Add(1)
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

	// Close results channel once all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results and display open ports immediately
	for port := range results {
		openPorts = append(openPorts, port)
	}

	// Finish the progress bar
	bar.Finish()

	return openPorts
}
