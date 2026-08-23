package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// Parameters follow the OWASP minimum pinned by ADR-008.
const (
	argon2Memory      = 19456 // KiB
	argon2Iterations  = 2
	argon2Parallelism = 1
	argon2SaltLength  = 16
	argon2KeyLength   = 32
)

var ErrMalformedPasswordHash = errors.New("malformed password hash")

const (
	maximumArgon2Memory      = 64 * 1024
	maximumArgon2Iterations  = 10
	maximumArgon2Parallelism = 4
)

type argon2idPasswordHasher struct{}

func New() out.PasswordHasher {
	return &argon2idPasswordHasher{}
}

func (h *argon2idPasswordHasher) Hash(password string) (string, error) {
	salt := make([]byte, argon2SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, argon2Iterations, argon2Memory, argon2Parallelism, argon2KeyLength)
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argon2Memory, argon2Iterations, argon2Parallelism,
		base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

// Verify re-derives the key using the parameters embedded in hash itself
// (not the current argon2* constants) so that a future parameter change
// doesn't break verification of hashes written under older parameters --
// ADR-008's "los parámetros quedan versionados junto al hash".
func (h *argon2idPasswordHasher) Verify(password, hash string) (bool, error) {
	memory, iterations, parallelism, salt, key, err := parseArgon2idHash(hash)
	if err != nil {
		return false, err
	}
	candidate := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(key)))
	return subtle.ConstantTimeCompare(candidate, key) == 1, nil
}

func parseArgon2idHash(encoded string) (memory, iterations uint32, parallelism uint8, salt, key []byte, err error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return 0, 0, 0, nil, nil, ErrMalformedPasswordHash
	}
	var version int
	if _, scanErr := fmt.Sscanf(parts[2], "v=%d", &version); scanErr != nil {
		return 0, 0, 0, nil, nil, ErrMalformedPasswordHash
	}
	if _, scanErr := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); scanErr != nil {
		return 0, 0, 0, nil, nil, ErrMalformedPasswordHash
	}
	if version != argon2.Version || memory < argon2Memory || memory > maximumArgon2Memory || iterations < argon2Iterations || iterations > maximumArgon2Iterations || parallelism < argon2Parallelism || parallelism > maximumArgon2Parallelism {
		return 0, 0, 0, nil, nil, ErrMalformedPasswordHash
	}
	salt, decodeErr := base64.RawStdEncoding.DecodeString(parts[4])
	if decodeErr != nil {
		return 0, 0, 0, nil, nil, ErrMalformedPasswordHash
	}
	key, decodeErr = base64.RawStdEncoding.DecodeString(parts[5])
	if decodeErr != nil || len(salt) < 16 || len(key) < 16 || len(key) > 64 {
		return 0, 0, 0, nil, nil, ErrMalformedPasswordHash
	}
	return memory, iterations, parallelism, salt, key, nil
}
