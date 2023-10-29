package security

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

const pwdCost = 14

// HashPassword computes the password hash.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), pwdCost)

	return string(bytes), err
}

// CreatePassword creates a new password with length.
func CreatePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'[],./!@#$%^&*()_+-="

	password := make([]byte, length)
	for i := 0; i < length; i++ {
		// #nosec G404 -- allow.
		r := rand.Intn(len(charset))
		password[i] = charset[r]
	}

	return string(password)
}
