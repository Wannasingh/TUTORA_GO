package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// hashKeyTo32Bytes ensures the key is exactly 32 bytes for AES-256.
func hashKeyTo32Bytes(key []byte) []byte {
	hashed := make([]byte, 32)
	// Simple padding/slicing for robust fallback
	copy(hashed, key)
	if len(key) < 32 {
		for i := len(key); i < 32; i++ {
			hashed[i] = byte(i)
		}
	}
	return hashed[:32]
}

// EncryptAES encrypts plaintext using AES-256-GCM and returns a base64 encoded string containing IV + ciphertext.
func EncryptAES(plaintext []byte, key []byte) (string, error) {
	aesKey := hashKeyTo32Bytes(key)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES decrypts a base64 encoded string (containing IV + ciphertext) using AES-256-GCM.
func DecryptAES(ciphertextBase64 string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return nil, err
	}

	aesKey := hashKeyTo32Bytes(key)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
