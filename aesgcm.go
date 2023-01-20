package enshamir

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

const (
	nonceLength = 12
	keyLength   = 32
)

// https://pkg.go.dev/crypto/cipher#example-NewGCM-Encrypt
func encrypt(key []byte, plaintext []byte) ([]byte, error) {
	if len(key) != keyLength {
		return nil, fmt.Errorf("AES-256 key should be 32 bytes")
	}

	aesgcm, err := newAESGCM(key)
	if err != nil {
		return nil, err
	}

	nonce, err := randomBytes(nonceLength)
	if err != nil {
		return nil, err
	}

	return aesgcm.Seal(nonce, nonce, plaintext, nil), nil
}

// https://pkg.go.dev/crypto/cipher#example-NewGCM-Decrypt
func decrypt(key []byte, cipherText []byte) ([]byte, error) {
	if len(key) != keyLength {
		return nil, fmt.Errorf("AES-256 key should be 32 bytes")
	}
	aesgcm, err := newAESGCM(key)
	if err != nil {
		return nil, err
	}

	nonce, data, err := extractNonce(cipherText)
	if err != nil {
		return nil, err
	}
	if len(nonce) != nonceLength {
		return nil, fmt.Errorf("nonce should be 12 bytes")
	}

	return aesgcm.Open(nil, nonce, data, nil)
}

func extractNonce(cipherText []byte) ([]byte, []byte, error) {
	if len(cipherText) < nonceLength {
		return nil, nil, fmt.Errorf("invalid data length")
	}
	nonce, data := cipherText[:nonceLength], cipherText[nonceLength:]

	return nonce, data, nil
}

func newAESGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm, nil
}
