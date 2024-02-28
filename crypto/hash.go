package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

/*
Create a hash of the secret. Uses Argon2id. Salt is generated randomly.

Args:

	secret: secret to hash

Returns:

	string: base64 representation hash
	error: error
*/
func Hash(secret []byte) (string, error) {
	salt := make([]byte, SALT_SIZE)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash_bytes := argon2.IDKey(secret, salt, ARGON_TIME, ARGON_MEMORY, ARGON_THREADS, ARGON_KEY_LEN)
	hash_b64 := base64.RawStdEncoding.EncodeToString(hash_bytes)
	salt_b64 := base64.RawStdEncoding.EncodeToString(salt)

	hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, ARGON_MEMORY, ARGON_TIME, ARGON_THREADS, salt_b64, hash_b64)

	return hash, nil
}
