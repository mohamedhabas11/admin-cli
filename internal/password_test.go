package internal

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	password, err := GeneratePassword(10)
	if err != nil {
		t.Fatalf("GeneratePassword() err = %v; want nil", err)
	}
	if len(password) != 10 {
		t.Errorf("len(password) = %d; want 10", len(password))
	}
}
