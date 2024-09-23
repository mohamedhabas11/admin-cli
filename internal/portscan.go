package internal

import (
	"fmt"
	"net"
	"sync"
)

// ScanPorts scans a range of ports on a given IP address
// and returns a slice of open ports
func ScanPorts(ip string, ports []int) []int {
	var openPorts []int
	var wg sync.WaitGroup
	results := make(chan int)

	// Function to scan a single port
	scan := func(port int) {
		defer wg.Done()
		target := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.Dial("tcp", target)
		if err != nil {
			return // Port is closed
		}
		conn.Close()
		results <- port // Send open port to channel
	}

	// Start goroutines for each port
	for _, port := range ports {
		wg.Add(1)
		go scan(port)
	}

	// Close results channel once all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results from the channel
	for port := range results {
		openPorts = append(openPorts, port)
	}

	return openPorts
}
