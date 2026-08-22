package security

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestSessionTokenGeneratorGenerate(t *testing.T) {
	t.Parallel()

	generator := NewSessionTokenGenerator()

	token, hash, err := generator.Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		t.Fatalf("token is not valid base64url: %v", err)
	}
	if len(decoded) < 32 {
		t.Fatalf("token has insufficient entropy: %d decoded bytes", len(decoded))
	}

	sum := sha256.Sum256([]byte(token))
	if hash != hex.EncodeToString(sum[:]) {
		t.Fatal("hash must be the SHA-256 hex digest of the returned token")
	}
}

func TestSessionTokenGeneratorProducesDistinctTokens(t *testing.T) {
	t.Parallel()

	generator := NewSessionTokenGenerator()

	firstToken, firstHash, err := generator.Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secondToken, secondHash, err := generator.Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if firstToken == secondToken {
		t.Fatal("two calls must not generate the same token")
	}
	if firstHash == secondHash {
		t.Fatal("two calls must not generate the same hash")
	}
}
