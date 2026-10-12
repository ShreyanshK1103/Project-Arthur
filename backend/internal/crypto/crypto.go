package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
)

var key []byte

// Init loads ENV_ENCRYPTION_KEY (base64-encoded, 32 bytes).
func Init() error {
	raw := os.Getenv("ENV_ENCRYPTION_KEY")
	if raw == "" {
		return errors.New("ENV_ENCRYPTION_KEY is not set")
	}

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return fmt.Errorf("ENV_ENCRYPTION_KEY is not valid base64: %w", err)
	}

	if len(decoded) != 32 {
		return fmt.Errorf("ENV_ENCRYPTION_KEY must decode to 32 bytes, got %d", len(decoded))
	}

	key = decoded
	return nil
}

func newGCM() (cipher.AEAD, error) {
	if key == nil {
		return nil, errors.New("crypto not initialised: call crypto.Init() at startup")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

// Encrypt returns nonce || ciphertext. A fresh random nonce is used every call
func Encrypt(plaintext string) ([]byte, error) {
	gcm, err := newGCM()
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

// Decrypt reverses Encrypt. It fails if the data was tampered with or the key is wrong.

func Decrypt(data []byte) (string, error) {
	gcm, err := newGCM()
	if err != nil {
		return "", err
	}

	if len(data) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
