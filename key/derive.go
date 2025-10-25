package key

import (
	"crypto/pbkdf2"
	"crypto/sha256"
)

// DeriveAESFromPassword derives a 32-byte AES key from a password and salt
func DeriveAESFromPassword(password string, salt []byte) ([]byte, error) {
	// PBKDF2 with SHA-256, 100000 iterations, 32-byte output for AES-256
	return pbkdf2.Key(sha256.New, password, salt, 100000, 32)
}
