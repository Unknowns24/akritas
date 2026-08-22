package out

// PasswordHasher hashes a plaintext password with Argon2id per ADR-008.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) (bool, error)
}
