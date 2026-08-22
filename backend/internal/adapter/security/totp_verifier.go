package security

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// RFC 6238 defaults matching TOTPSecretGenerator and ADR-008: 6 digits,
// SHA1, 30s period, one period of tolerance before/after.
type totpVerifier struct{}

func NewTOTPVerifier() out.TOTPVerifier {
	return &totpVerifier{}
}

func (v *totpVerifier) Verify(secret, code string, at time.Time) (bool, error) {
	return totp.ValidateCustom(code, secret, at, totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
}
