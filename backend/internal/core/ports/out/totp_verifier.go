package out

import "time"

// TOTPVerifier checks a submitted code against a decrypted secret, per
// RFC 6238 with the tolerance ADR-008 requires. Separate from
// TOTPSecretGenerator (single responsibility: generation vs. verification).
//
// period is the RFC 6238 time-step counter that matched, valid only when
// valid is true. Callers that need to reject reuse of an already-accepted
// period (ADR-008) compare it against the last one they persisted.
type TOTPVerifier interface {
	Verify(secret, code string, at time.Time) (valid bool, period int64, err error)
}
