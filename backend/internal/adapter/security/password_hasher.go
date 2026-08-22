package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

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

type argon2idPasswordHasher struct{}

func NewPasswordHasher() out.PasswordHasher {
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
