package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

const sessionTokenLength = 32 // bytes, before base64url encoding

type sessionTokenGenerator struct{}

func NewSessionTokenGenerator() out.SessionTokenGenerator {
	return &sessionTokenGenerator{}
}

func (g *sessionTokenGenerator) Generate() (string, string, error) {
	raw := make([]byte, sessionTokenLength)
	if _, err := rand.Read(raw); err != nil {
		return "", "", err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token))
	return token, hex.EncodeToString(sum[:]), nil
}
