package security

import (
	"crypto/subtle"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type constantTimeBootstrapTokenVerifier struct {
	token []byte
}

func NewBootstrapTokenVerifier(token string) out.BootstrapTokenVerifier {
	return &constantTimeBootstrapTokenVerifier{token: []byte(token)}
}

func (v *constantTimeBootstrapTokenVerifier) Verify(candidate string) bool {
	if len(v.token) == 0 || len(candidate) == 0 {
		return false
	}
	return subtle.ConstantTimeCompare(v.token, []byte(candidate)) == 1
}
