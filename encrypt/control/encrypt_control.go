package control

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	_ "image/jpeg"
	"io"

	"github.com/audryus/steganocc/key"
)

func Encrypt(password, message string) (string, error) {
	// Generate salt and derive key
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %v", err)
	}
	myKey, err := key.DeriveAESFromPassword(password, salt)
	if err != nil {
		return "", fmt.Errorf("derive AES Key From Password: %v", err)
	}
	// Encrypt the message (returns nonce+ciphertext)
	encrypted, err := encryptMessage([]byte(message), myKey)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %v", err)
	}
	// Append the salt so it can be extracted later for key derivation
	encrypted = append(encrypted, salt...)
	encoded := base64.StdEncoding.EncodeToString(encrypted)

	return encoded, nil
}

// encryptMessage encrypts the plaintext using AES-GCM and returns nonce+ciphertext
func encryptMessage(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}
