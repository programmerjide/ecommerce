package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	// Implement password hashing logic here
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword compares a hashed password with its possible plaintext equivalent
func VerifyPassword(hash, password string) error {
	// Implement password verification logic here
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
