package enshamir

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptSplit(t *testing.T) {
	encryptionPassword := []byte("hello1234567890!@#$%^&*()_+~`")
	secret := []byte("this is a secret `1234567890-=~!@#$%^&*()_+")

	parts := 4
	threshold := 3

	hashedPassword, shares, err := EncryptSplit(encryptionPassword, secret, parts, threshold)
	assert.NoError(t, err)

	// Verify if we can really combine and decrypt the shares.
	validIndexes := [][]int{{
		0, 1, 2, 3, // all parts
	}, {
		0, 1, 2,
	}, {
		0, 1, 3,
	}, {
		0, 2, 3,
	}, {
		1, 2, 3,
	}}

	for _, indexes := range validIndexes {
		var s [][]byte
		for _, i := range indexes {
			s = append(s, shares[i])
		}
		assert.True(t, len(s) >= threshold)

		plaintext, err := CombineDecrypt(encryptionPassword, hashedPassword, s)
		if assert.NoError(t, err) {
			assert.Equal(t, string(secret), string(plaintext))
		}
	}

	invalidIndexes := [][]int{{0, 1}, {1, 2}, {0, 3}, {1, 3}}
	for _, indexes := range invalidIndexes {
		var s [][]byte
		for _, i := range indexes {
			s = append(s, shares[i])
		}
		assert.True(t, len(s) < threshold)

		// It's expected to fail to decrypt the shares since we don't have enough shares.
		_, err := CombineDecrypt(encryptionPassword, hashedPassword, s)
		assert.True(t, strings.Contains(err.Error(), "unable to decrypt the secret"))
	}
}

func Test_randomBytes(t *testing.T) {
	for length := 0; length <= 256; length++ {
		b, err := randomBytes(uint32(length))
		assert.NoError(t, err)
		assert.Len(t, b, length)
	}
}
