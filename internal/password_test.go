package internal

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	length := 12
	pass, err := GeneratePassword(length)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(pass) != length {
		t.Errorf("Expected password length %d, got %d", length, len(pass))
	}

	if pass == "" {
		t.Error("Generated password should not be empty")
	}
}
