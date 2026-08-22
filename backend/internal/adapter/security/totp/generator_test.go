package totp

import (
	"regexp"
	"strings"
	"testing"
)

var manualEntryKeyPattern = regexp.MustCompile(`^[A-Z2-7]{16,128}$`)

func TestTOTPSecretGeneratorGenerate(t *testing.T) {
	t.Parallel()

	generator := NewGenerator()

	secret, err := generator.Generate("Akritas", "admin@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !manualEntryKeyPattern.MatchString(secret.Base32Key) {
		t.Fatalf("manual_entry_key %q does not match OpenAPI pattern", secret.Base32Key)
	}
	if !strings.HasPrefix(secret.OtpauthURI, "otpauth://totp/") {
		t.Fatalf("otpauth_uri does not start with otpauth://totp/: %q", secret.OtpauthURI)
	}
	if !strings.Contains(secret.OtpauthURI, "admin@example.com") {
		t.Fatalf("otpauth_uri does not reference the account email: %q", secret.OtpauthURI)
	}
	if !strings.Contains(secret.OtpauthURI, secret.Base32Key) {
		t.Fatalf("otpauth_uri does not embed the generated secret: %q", secret.OtpauthURI)
	}
}

func TestTOTPSecretGeneratorProducesDistinctSecrets(t *testing.T) {
	t.Parallel()

	generator := NewGenerator()

	first, err := generator.Generate("Akritas", "admin@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := generator.Generate("Akritas", "admin@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first.Base32Key == second.Base32Key {
		t.Fatal("two calls must not generate the same secret")
	}
}
