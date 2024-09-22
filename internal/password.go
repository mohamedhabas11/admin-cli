package internal

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"

func GeneratePassword(length int) (string, error) {
	password := make([]byte, length)
	for char := range password {
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[char] = charset[randIndex.Int64()]
	}
	return string(password), nil
}
