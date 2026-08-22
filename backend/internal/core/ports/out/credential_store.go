package out

import "context"

// CredentialStore encrypts secret material at rest with the infrastructure
// master key (ADR-005). Decrypt is intentionally absent: this task never
// reads a TOTP secret back, only writes it.
type CredentialStore interface {
	Encrypt(ctx context.Context, plaintext []byte) ([]byte, error)
}
