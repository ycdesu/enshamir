package enshamir

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_encrypt(t *testing.T) {
	key, err := randomBytes(32)
	assert.NoError(t, err)

	plaintext := []byte("hello this is plaintext`1234567890-=~!@#$%^&*()_+")
	ciphertext, err := encrypt(key, plaintext)
	assert.NoError(t, err)

	decrypted, err := decrypt(key, ciphertext)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(plaintext, decrypted))
}
