package security

import (
	"strings"
	"testing"
)

func TestPasswordHasherHash(t *testing.T) {
	t.Parallel()

	hasher := NewPasswordHasher()

	hash, err := hasher.Hash("a-long-password-from-a-password-manager")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Fatalf("hash does not use the argon2id prefix: %q", hash)
	}
	if !strings.Contains(hash, "m=19456,t=2,p=1") {
		t.Fatalf("hash does not embed the ADR-008 parameters: %q", hash)
	}
}

func TestPasswordHasherProducesDistinctHashesForSamePassword(t *testing.T) {
	t.Parallel()

	hasher := NewPasswordHasher()

	first, err := hasher.Hash("a-long-password-from-a-password-manager")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := hasher.Hash("a-long-password-from-a-password-manager")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first == second {
		t.Fatal("two hashes of the same password must differ (random salt)")
	}
}
