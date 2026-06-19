package utils

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptAES_Success(t *testing.T) {
	key := []byte("MySuperSecretKeyForEncryption123")
	plaintext := []byte("Hello World! This is a sensitive payload.")

	ciphertext, err := EncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	if len(ciphertext) == 0 {
		t.Fatal("expected non-empty ciphertext")
	}

	decrypted, err := DecryptAES(ciphertext, key)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("decrypted payload does not match plaintext: expected %s, got %s", plaintext, decrypted)
	}
}

func TestDecryptAES_InvalidKey(t *testing.T) {
	key1 := []byte("MySuperSecretKeyForEncryption123")
	key2 := []byte("IncorrectKeyWhichShouldNotDecode")
	plaintext := []byte("Hello World!")

	ciphertext, _ := EncryptAES(plaintext, key1)

	_, err := DecryptAES(ciphertext, key2)
	if err == nil {
		t.Fatal("expected decryption to fail with incorrect key, but it succeeded")
	}
}

func TestDecryptAES_InvalidBase64(t *testing.T) {
	key := []byte("MySuperSecretKeyForEncryption123")
	_, err := DecryptAES("invalid-base64-string", key)
	if err == nil {
		t.Fatal("expected base64 decoding to fail, got nil error")
	}
}
