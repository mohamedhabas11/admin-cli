package internal

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"
	"testing"
)

// TestComputeFileHash tests the ComputeFileHash function
func TestComputeFileHash(t *testing.T) {
	// create a temporary file with known content
	tempFile, err := os.CreateTemp("", "checksum_test")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	// clean up the temporary file
	defer os.Remove(tempFile.Name())

	// known content
	content := []byte("Hello, World!")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	// close file before computing the hash
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// compute the hash of the temporary file manually
	expectedHash := sha256.Sum256(content)

	// Call the ComputeFileHash function under test
	actualHash, err := ComputeFileHash(tempFile.Name())
	if err != nil {
		t.Fatalf("ComputeFileHash() err = %v; want nil", err)
	}

	// compare the actual hash with the expected hash
	if actualHash != expectedHash {
		t.Errorf("ComputeFileHash() = %x; want %x", actualHash, expectedHash)
	}
}

// TestPrintCollisions checks if collisions are correctly detected and printed
func TestPrintCollisions(t *testing.T) {
	hash1 := sha256.Sum256([]byte("file1_content"))
	hash2 := sha256.Sum256([]byte("file2_content"))

	// Simulate a collision by having multiple files with the same hash
	hashes := map[[32]byte][]string{
		hash1: {"file1", "file3"},
		hash2: {"file2"},
	}

	// Create a pipe to capture os.Stdout output
	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w

	// Call the function to be tested
	PrintCollisions(hashes)

	// Restore original os.Stdout
	w.Close()
	os.Stdout = origStdout

	// Read the captured output from the pipe
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Check if the output contains the expected message
	expected := "Hash collision detected for hash"
	if !bytes.Contains(buf.Bytes(), []byte(expected)) {
		t.Errorf("Expected collision message in output, got: %s", buf.String())
	}
}
