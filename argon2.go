package enshamir

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argon2idParams struct {
	// memory 64 * 1024 sets the memory cost to ~64 MB
	memory uint32

	// The number of iterations over the memory.
	times uint32

	// the number of threads. It can be set to available cpus by `runtime.NumCPU()`.
	thread uint8

	saltLength uint32

	// AES-256 needs 32-byte key slice
	keyLength uint32
}

func (p argon2idParams) String() string {
	return fmt.Sprintf("memory: %d, times: %d, thread: %d, saltLength: %d, keyLength: %d",
		p.memory,
		p.times,
		p.thread,
		p.saltLength,
		p.keyLength,
	)
}

var defaultArgon2idParams = argon2idParams{
	memory:     2048 * 1024, // ~2GB
	times:      4,
	thread:     uint8(runtime.NumCPU()),
	saltLength: 16,
	keyLength:  32,
}

// hashPassword hashes the passwords to a 32 bytes slice by argon2id which will be used in AES-256-GCM encryption.
func hashPasswordWithSalt(password, salt []byte) []byte {
	return argon2.IDKey(
		password,
		salt,
		defaultArgon2idParams.times,
		defaultArgon2idParams.memory,
		defaultArgon2idParams.thread,
		defaultArgon2idParams.keyLength,
	)
}

// verifyPassword verifies the password against the hash. Currently it's only used in unit test.
func verifyPassword(password, hash string) error {
	p, s, k, err := decode(hash)
	if err != nil {
		return fmt.Errorf("unable to decode hash: %w", err)
	}

	newKey := argon2.IDKey(
		[]byte(password),
		s,
		p.times,
		p.memory,
		p.thread,
		p.keyLength,
	)

	if subtle.ConstantTimeCompare(k, newKey) == 0 {
		return fmt.Errorf("password does not match")
	}

	return nil
}

// The returned hash will be encoded in the format:
// $argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
func encode(p argon2idParams, salt, key []byte) string {
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.memory,
		p.times,
		p.thread,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key))
}

func decode(hash string) (argon2idParams, []byte, []byte, error) {
	var p argon2idParams
	values := strings.Split(hash, "$")
	if len(values) != 6 {
		return p, nil, nil, fmt.Errorf("invalid hash format")
	}
	if values[1] != "argon2id" {
		return p, nil, nil, fmt.Errorf("incompatible argon2 variant")
	}

	var version int
	_, err := fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return p, nil, nil, fmt.Errorf("invalid version: %w", err)
	}
	if version != argon2.Version {
		return p, nil, nil, fmt.Errorf("incompatible argon2 version: %d", version)
	}

	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &p.memory, &p.times, &p.thread)
	if err != nil {
		return p, nil, nil, fmt.Errorf("unable to parse argon2 parameters: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(values[4])
	if err != nil {
		return p, nil, nil, fmt.Errorf("unable to decode salt: %w", err)
	}

	p.saltLength = uint32(len(salt))

	key, err := base64.RawStdEncoding.DecodeString(values[5])
	if err != nil {
		return p, nil, nil, fmt.Errorf("unable to decode key: %w", err)
	}
	p.keyLength = uint32(len(key))

	return p, salt, key, nil
}
