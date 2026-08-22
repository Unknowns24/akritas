package security

import (
	"crypto/subtle"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// RFC 6238 defaults matching TOTPSecretGenerator and ADR-008: 6 digits,
// SHA1, 30s period, one period of tolerance before/after.
const totpPeriodSeconds = 30

type totpVerifier struct{}

func NewTOTPVerifier() out.TOTPVerifier {
	return &totpVerifier{}
}

// Verify does not delegate to totp.ValidateCustom, which reports only a
// bool: ADR-008 also requires rejecting reuse of an already-accepted
// period, so callers need to know which of the tolerated periods matched.
// It tries the current, previous and next period, comparing each candidate
// code in constant time.
func (v *totpVerifier) Verify(secret, code string, at time.Time) (bool, int64, error) {
	counter := at.Unix() / totpPeriodSeconds
	for _, offset := range []int64{0, -1, 1} {
		candidateCounter := counter + offset
		candidateTime := time.Unix(candidateCounter*totpPeriodSeconds, 0)
		candidateCode, err := totp.GenerateCodeCustom(secret, candidateTime, totp.ValidateOpts{
			Period:    totpPeriodSeconds,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
		if err != nil {
			return false, 0, err
		}
		if subtle.ConstantTimeCompare([]byte(candidateCode), []byte(code)) == 1 {
			return true, candidateCounter, nil
		}
	}
	return false, 0, nil
}
