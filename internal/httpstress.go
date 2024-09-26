package internal

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// PerformStressTest performs the HTTP stress testing
func PerformStressTest(url string, numRequests, concurrency int) {
	var wg sync.WaitGroup
	start := time.Now()
	successCount := 0
	var successMutex sync.Mutex

	// Channel to collect errors
	errors := make(chan error, numRequests)

	// Create a channel for request IDs
	requests := make(chan int, numRequests)

	// Populate the requests channel
	for i := 0; i < numRequests; i++ {
		requests <- i
	}
	close(requests) // Close the channel once populated

	// Function for each worker
	stress := func() {
		defer wg.Done()
		for range requests {
			resp, err := http.Get(url)
			if err != nil {
				errors <- err
				continue
			}
			if resp.StatusCode == http.StatusOK {
				successMutex.Lock()
				successCount++
				successMutex.Unlock()
			} else {
				errors <- fmt.Errorf("non-200 response: %d", resp.StatusCode)
			}
			resp.Body.Close()
		}
	}

	// Launch workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go stress()
	}

	wg.Wait()     // Wait for all workers to finish
	close(errors) // Close the errors channel

	// Report results
	elapsed := time.Since(start)
	fmt.Printf("Completed %d requests in %s\n", numRequests, elapsed)
	fmt.Printf("Successful requests: %d\n", successCount)

	// Report any errors
	for err := range errors {
		fmt.Printf("Error: %v\n", err)
	}
}
