package password

import (
	"errors"
	"strings"
	"testing"
)

func TestPasswordHasherHash(t *testing.T) {
	t.Parallel()

	hasher := New()

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

	hasher := New()

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

func TestPasswordHasherVerify(t *testing.T) {
	t.Parallel()

	hasher := New()

	hash, err := hasher.Hash("a-long-password-from-a-password-manager")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ok, err := hasher.Verify("a-long-password-from-a-password-manager", hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("the correct password must verify")
	}

	ok, err = hasher.Verify("a-different-password-entirely", hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("a wrong password must not verify")
	}
}

func TestPasswordHasherVerifyRejectsMalformedHash(t *testing.T) {
	t.Parallel()

	hasher := New()

	if _, err := hasher.Verify("anything", "not-a-valid-hash"); !errors.Is(err, ErrMalformedPasswordHash) {
		t.Fatalf("expected ErrMalformedPasswordHash, got %v", err)
	}
}

func TestPasswordHasherRejectsHostileParametersBeforeDerivation(t *testing.T) {
	t.Parallel()
	hasher := New()
	for _, hash := range []string{
		"$argon2id$v=19$m=1048576,t=2,p=1$MDEyMzQ1Njc4OWFiY2RlZg$MDEyMzQ1Njc4OWFiY2RlZg",
		"$argon2id$v=19$m=19456,t=999,p=1$MDEyMzQ1Njc4OWFiY2RlZg$MDEyMzQ1Njc4OWFiY2RlZg",
		"$argon2id$v=16$m=19456,t=2,p=1$MDEyMzQ1Njc4OWFiY2RlZg$MDEyMzQ1Njc4OWFiY2RlZg",
	} {
		if _, err := hasher.Verify("password", hash); !errors.Is(err, ErrMalformedPasswordHash) {
			t.Fatalf("hostile hash was accepted: %v", err)
		}
	}
}
