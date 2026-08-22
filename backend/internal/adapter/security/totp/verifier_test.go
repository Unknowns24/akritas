package totp

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

func periodOf(at time.Time) int64 {
	return at.Unix() / totpPeriodSeconds
}

func TestTOTPVerifierAcceptsCurrentAndAdjacentPeriods(t *testing.T) {
	t.Parallel()

	verifier := NewVerifier()
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
			ok, period, err := verifier.Verify(testTOTPSecret, codeAt(t, codeTime), now)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatalf("code from %v must be accepted at %v", codeTime, now)
			}
			if want := periodOf(codeTime); period != want {
				t.Fatalf("period = %d, want %d (the counter the code actually belongs to)", period, want)
			}
		})
	}
}

func TestTOTPVerifierRejectsOutOfToleranceAndWrongCode(t *testing.T) {
	t.Parallel()

	verifier := NewVerifier()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)

	t.Run("two periods away", func(t *testing.T) {
		t.Parallel()
		ok, period, err := verifier.Verify(testTOTPSecret, codeAt(t, now.Add(-90*time.Second)), now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("a code two periods away must not be accepted")
		}
		if period != 0 {
			t.Fatalf("period = %d, want 0 on rejection", period)
		}
	})

	t.Run("wrong code", func(t *testing.T) {
		t.Parallel()
		valid := codeAt(t, now)
		wrong := "000000"
		if wrong == valid {
			wrong = "111111"
		}
		ok, period, err := verifier.Verify(testTOTPSecret, wrong, now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatal("an arbitrary wrong code must not be accepted")
		}
		if period != 0 {
			t.Fatalf("period = %d, want 0 on rejection", period)
		}
	})
}
