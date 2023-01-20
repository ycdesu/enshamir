package enshamir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_hashPasswordWithSalt(t *testing.T) {
	passwords := []string{"1", "a", "", "0", "~", "hello1234567890!@#$%^&*()_+~`\""}

	for _, p := range passwords {
		salt, err := randomBytes(defaultArgon2idParams.saltLength)
		assert.NoError(t, err)

		hashedKey := hashPasswordWithSalt([]byte(p), salt)
		assert.Len(t, hashedKey, int(defaultArgon2idParams.keyLength))

		assert.NoError(t, verifyPassword(p, encode(defaultArgon2idParams, salt, hashedKey)))
	}
}
