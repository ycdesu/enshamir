package enshamir

import (
	"testing"
)

func Test_hashPasswordWithSalt(t *testing.T) {
	passwords := []string{"112345678912345678912345678923456789112345678912345678912345678923456789"}

	for _, p := range passwords {
		salt, err := randomBytes(defaultArgon2idParams.saltLength)
		if err != nil {
			t.Fatal(err)
		}

		hashedKey := hashPasswordWithSalt([]byte(p), salt)
		if len(hashedKey) != int(defaultArgon2idParams.keyLength) {
			t.Fatal("hashed key length is not correct")
		}

		if err := verifyPassword(p, encode(defaultArgon2idParams, salt, hashedKey)); err != nil {
			t.Fatal(err)
		}
	}
}
