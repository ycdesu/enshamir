package enshamir

import (
	"crypto/rand"
	"fmt"

	"github.com/hashicorp/vault/shamir"
)

// EncryptSplit encrypts the secret with password and generate a `parts` number of shares.
// `threshold` is the minimum number of shares required to reconstruct the secret.
//
// We use argon2id to hash the password and generate a 32-byte key. You must save the salt which will be used to decrypt
// your secret.
func EncryptSplit(password, secret []byte, parts, threshold int) ([]byte, [][]byte, error) {
	salt, err := randomBytes(defaultArgon2idParams.saltLength)
	if err != nil {
		return nil, nil, err
	}

	key := hashPasswordWithSalt(password, salt)

	encryptedSecret, err := encrypt(key, secret)
	if err != nil {
		return nil, nil, err
	}

	shares, err := shamir.Split(encryptedSecret, parts, threshold)
	if err != nil {
		return nil, nil, err
	}

	return salt, shares, nil
}

func CombineDecrypt(password, salt []byte, shares [][]byte) ([]byte, error) {
	encryptedSecret, err := shamir.Combine(shares)
	if err != nil {
		return nil, fmt.Errorf("unable to combine shares: %w", err)
	}

	key := hashPasswordWithSalt(password, salt)

	decrypted, err := decrypt(key, encryptedSecret)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt the secret: %w", err)
	}
	return decrypted, nil
}

func randomBytes(len uint32) ([]byte, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
