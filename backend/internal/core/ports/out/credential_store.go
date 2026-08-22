package out

import "context"

// CredentialStore encrypts and decrypts secret material at rest with the
// infrastructure master key (ADR-005). Decrypt must only be called at the
// point of use.
type CredentialStore interface {
	Encrypt(ctx context.Context, plaintext []byte) ([]byte, error)
	Decrypt(ctx context.Context, ciphertext []byte) ([]byte, error)
}
