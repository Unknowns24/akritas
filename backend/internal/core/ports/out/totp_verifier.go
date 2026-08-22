package out

import "time"

// TOTPVerifier checks a submitted code against a decrypted secret, per
// RFC 6238 with the tolerance ADR-008 requires. Separate from
// TOTPSecretGenerator (single responsibility: generation vs. verification).
type TOTPVerifier interface {
	Verify(secret, code string, at time.Time) (bool, error)
}
