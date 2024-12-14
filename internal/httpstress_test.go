package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestPerformStressTest tests the PerformStressTest function with a successful response
func TestPerformStressTest(t *testing.T) {
	// Create a test HTTP server that returns a 200 OK response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a 200 OK status
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Define the number of requests and concurrency
	numRequests := 3
	concurrency := 1

	// Run the stress test with the test server URL
	PerformStressTest(server.URL, numRequests, concurrency)

	// You can assert or check for side effects based on your own logic
	// In this case, the test should pass if no errors are printed
}

// TestPerformStressTest_WithError tests the PerformStressTest function with a server error
func TestPerformStressTest_WithError(t *testing.T) {
	// Create a test HTTP server that returns a 500 Internal Server Error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate an internal server error (500)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Define the number of requests and concurrency
	numRequests := 3
	concurrency := 1

	// Run the stress test with the test server URL
	PerformStressTest(server.URL, numRequests, concurrency)

	// This test is focused on ensuring the error handling works correctly,
	// so checking for output or checking the log could be beneficial
}
