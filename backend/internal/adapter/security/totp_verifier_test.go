package security

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

const testTOTPSecret = "JBSWY3DPEHPK3PXP"

func codeAt(t *testing.T, at time.Time) string {
	t.Helper()
	code, err := totp.GenerateCode(testTOTPSecret, at)
	if err != nil {
		t.Fatalf("generate code: %v", err)
	}
	return code
}

func TestTOTPVerifierAcceptsCurrentAndAdjacentPeriods(t *testing.T) {
	t.Parallel()

	verifier := NewTOTPVerifier()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)

	cases := map[string]time.Time{
		"current period":  now,
		"previous period": now.Add(-30 * time.Second),
		"next period":     now.Add(30 * time.Second),
	}

	for name, codeTime := range cases {
		codeTime := codeTime
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ok, err := verifier.Verify(testTOTPSecret, codeAt(t, codeTime), now)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatalf("code from %v must be accepted at %v", codeTime, now)
			}
		})
	}
}

func TestTOTPVerifierRejectsOutOfToleranceAndWrongCode(t *testing.T) {
	t.Parallel()

	verifier := NewTOTPVerifier()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)

	t.Run("two periods away", func(t *testing.T) {
		t.Parallel()
		ok, err := verifier.Verify(testTOTPSecret, codeAt(t, now.Add(-90*time.Second)), now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("a code two periods away must not be accepted")
		}
	})

	t.Run("wrong code", func(t *testing.T) {
		t.Parallel()
		valid := codeAt(t, now)
		wrong := "000000"
		if wrong == valid {
			wrong = "111111"
		}
		ok, err := verifier.Verify(testTOTPSecret, wrong, now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("an arbitrary wrong code must not be accepted")
		}
	})
}
