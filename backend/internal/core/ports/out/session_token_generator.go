package out

// SessionTokenGenerator produces a new high-entropy opaque session token
// for the cookie, and its at-rest hash for persistence. The raw token is
// never persisted.
type SessionTokenGenerator interface {
	Generate() (token string, hash string, err error)
	// Hash deterministically computes the same at-rest representation
	// Generate returns for a token it created, so callers (e.g. the auth
	// middleware) can look a token up without generating a new one.
	Hash(token string) string
}
