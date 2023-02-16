package enshamir

import (
	"strings"
	"testing"
)

func TestEncryptSplit(t *testing.T) {
	encryptionPassword := []byte("hello1234567890!@#$%^&*()_+~`")
	secret := []byte("this is a secret `1234567890-=~!@#$%^&*()_+")

	parts := 4
	threshold := 3

	hashedPassword, shares, err := EncryptSplit(encryptionPassword, secret, parts, threshold)
	if err != nil {
		t.Fatal(err)
	}

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
		if len(s) < threshold {
			t.FailNow()
		}

		plaintext, err := CombineDecrypt(encryptionPassword, hashedPassword, s)
		if err != nil {
			t.Fatal(err)
		}

		if string(secret) != string(plaintext) {
			t.FailNow()
		}
	}

	invalidIndexes := [][]int{{0, 1}, {1, 2}, {0, 3}, {1, 3}}
	for _, indexes := range invalidIndexes {
		var s [][]byte
		for _, i := range indexes {
			s = append(s, shares[i])
		}

		// This case is ensure we can't restore our plaintext if the number of shares is not enough
		if len(s) >= threshold {
			t.FailNow()
		}

		// It's expected to fail to decrypt the shares since we don't have enough shares.
		_, err := CombineDecrypt(encryptionPassword, hashedPassword, s)
		if !strings.Contains(err.Error(), "unable to decrypt the secret") {
			t.Fatal("unexpected err: " + err.Error())
		}
	}
}

func Test_randomBytes(t *testing.T) {
	for length := 0; length <= 256; length++ {
		b, err := randomBytes(uint32(length))
		if err != nil {
			t.Fatal(err)
		}

		if len(b) != length {
			t.Fatalf("length does not match, expected: %d, actual: %d", length, len(b))
		}
	}
}
