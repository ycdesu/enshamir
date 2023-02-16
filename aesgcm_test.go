package enshamir

import (
	"bytes"
	"testing"
)

func Test_encrypt(t *testing.T) {
	key, err := randomBytes(32)
	if err != nil {
		t.Fatal(err)
	}

	plaintext := []byte("hello this is plaintext`1234567890-=~!@#$%^&*()_+")
	ciphertext, err := encrypt(key, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := decrypt(key, ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(plaintext, decrypted) {
		t.FailNow()
	}
}
