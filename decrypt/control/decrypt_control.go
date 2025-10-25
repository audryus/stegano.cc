package control

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/audryus/steganocc/key"
)

// Decrypt reads an image from file, extracts the embedded encrypted payload,
// derives AES key from password+salt and returns decrypted plaintext.
func Decrypt(password, message string, file io.Reader) (string, string, error) {
	var messageDecrypted []byte
	var fileDecrypted []byte
	var err error

	if file != nil {
		img, err := decodeImage(file)
		if err != nil {
			return "", "", err
		}
		fileDecrypted, err = extractEmbeddedPayload(img)
		if err != nil {
			return "", "", err
		}
		fileDecrypted, err = base64.StdEncoding.DecodeString(string(fileDecrypted))
		if err != nil {
			return "", "", err
		}
		salt := fileDecrypted[len(fileDecrypted)-16:]
		derivedKey, err := key.DeriveAESFromPassword(password, salt)
		if err != nil {
			return "", "", fmt.Errorf("derive AES Key From Password: %v", err)
		}
		fileDecrypted, err = decryptMessage(fileDecrypted, derivedKey)
		if err != nil {
			return "", "", fmt.Errorf("decryption failed: %v", err)
		}
	}

	if len(message) > 0 {
		messageDecrypted, err = base64.StdEncoding.DecodeString(message)
		if err != nil {
			return "", "", err
		}
		salt := messageDecrypted[len(messageDecrypted)-16:]
		derivedKey, err := key.DeriveAESFromPassword(password, salt)
		if err != nil {
			return "", "", fmt.Errorf("derive AES Key From Password: %v", err)
		}
		messageDecrypted, err = decryptMessage(messageDecrypted, derivedKey)
		if err != nil {
			return "", "", fmt.Errorf("decryption failed: %v", err)
		}
	}

	return string(messageDecrypted), string(fileDecrypted), nil
}

func decodeImage(file io.Reader) (image.Image, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}
	return img, nil
}

// extractEmbeddedPayload extracts the embedded bytes written LSB-first across
// RGB channels. Returns the payload (excluding the 4-byte length header).
func extractEmbeddedPayload(img image.Image) ([]byte, error) {
	bounds := img.Bounds()

	var (
		currByte    byte
		bitsFilled  int
		collected   []byte
		targetTotal int = -1 // unknown until first 4 bytes read
		cbErr       error
	)

	// callback invoked for each LSB bit; returning true stops iteration.
	cb := func(bit uint8) bool {
		stop, err := accumulateBit(bit, &currByte, &bitsFilled, &collected, &targetTotal, bounds)
		if err != nil {
			cbErr = err
			return true
		}
		return stop
	}

	iterateLSBBits(img, cb)

	if cbErr != nil {
		return nil, cbErr
	}
	if targetTotal == -1 || len(collected) < 4 {
		return nil, fmt.Errorf("failed to extract embedded payload")
	}
	if len(collected) < targetTotal {
		return nil, fmt.Errorf("extracted bytes shorter than expected")
	}
	// collected contains: [4 bytes length][payload...]
	return collected[4:targetTotal], nil
}

// accumulateBit handles assembling bits into bytes and checks payload length/sanity.
// Returns (stopIteration, error).
func accumulateBit(bit uint8, currByte *byte, bitsFilled *int, collected *[]byte, targetTotal *int, bounds image.Rectangle) (bool, error) {
	// assemble byte
	*currByte = (*currByte << 1) | bit
	*bitsFilled++
	if *bitsFilled != 8 {
		return false, nil
	}

	// full byte collected
	*collected = append(*collected, *currByte)
	*currByte = 0
	*bitsFilled = 0

	// determine expected total after first 4 bytes
	if *targetTotal == -1 && len(*collected) >= 4 {
		payloadLen := int(binary.BigEndian.Uint32((*collected)[:4]))
		*targetTotal = 4 + payloadLen
		// sanity check
		maxPayload := maxEmbeddedPayload(bounds)
		if payloadLen <= 0 || payloadLen > maxPayload {
			return true, fmt.Errorf("invalid embedded payload length: %d", payloadLen)
		}
	}

	// stop when we've collected the expected total
	if *targetTotal != -1 && len(*collected) >= *targetTotal {
		return true, nil
	}
	return false, nil
}

// iterateLSBBits walks image pixels and calls cb for each LSB bit (R,G,B order).
// If cb returns true iteration stops early.
func iterateLSBBits(img image.Image, cb func(bit uint8) bool) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			for _, bit := range []uint8{r8 & 1, g8 & 1, b8 & 1} {
				if cb(bit) {
					return
				}
			}
		}
	}
}

func maxEmbeddedPayload(bounds image.Rectangle) int {
	return (bounds.Dx() * bounds.Dy() * 3) / 8
}

// decryptMessage decrypts the nonce+ciphertext+salt using AES-GCM
func decryptMessage(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize+16 {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, rest := ciphertext[:nonceSize], ciphertext[nonceSize:]
	ciphertext, salt := rest[:len(rest)-16], rest[len(rest)-16:]
	_ = salt // Salt is already used to derive key
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}
	return plaintext, nil
}
