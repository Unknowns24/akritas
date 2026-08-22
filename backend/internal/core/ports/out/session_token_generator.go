package out

// SessionTokenGenerator produces a new high-entropy opaque session token
// for the cookie, and its at-rest hash for persistence. The raw token is
// never persisted.
type SessionTokenGenerator interface {
	Generate() (token string, hash string, err error)
}
