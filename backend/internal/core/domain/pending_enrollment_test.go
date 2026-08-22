package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func validPendingEnrollmentArgs() (uuid.UUID, string, string, string, []byte, time.Time, time.Time) {
	createdAt := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	return uuid.New(), "admin@example.com", "Akritas Administrator", "argon2id-hash",
		[]byte("ciphertext"), createdAt, createdAt.Add(10 * time.Minute)
}

func TestNewPendingEnrollmentValid(t *testing.T) {
	t.Parallel()

	id, email, displayName, passwordHash, secret, createdAt, expiresAt := validPendingEnrollmentArgs()
	enrollment, err := NewPendingEnrollment(id, email, displayName, passwordHash, secret, createdAt, expiresAt)
	if err != nil {
		t.Fatalf("valid pending enrollment rejected: %v", err)
	}
	if enrollment.ID != id || enrollment.Email != email || enrollment.DisplayName != displayName ||
		enrollment.PasswordHash != passwordHash || string(enrollment.EncryptedTOTPSecret) != string(secret) ||
		!enrollment.CreatedAt.Equal(createdAt) || !enrollment.ExpiresAt.Equal(expiresAt) {
		t.Fatal("constructed enrollment does not match input")
	}
}

func TestNewPendingEnrollmentRejectsInvalidFields(t *testing.T) {
	t.Parallel()

	id, email, displayName, passwordHash, secret, createdAt, expiresAt := validPendingEnrollmentArgs()

	cases := map[string]struct {
		id           uuid.UUID
		email        string
		displayName  string
		passwordHash string
		secret       []byte
		createdAt    time.Time
		expiresAt    time.Time
	}{
		"zero id":                {uuid.Nil, email, displayName, passwordHash, secret, createdAt, expiresAt},
		"blank email":            {id, "", displayName, passwordHash, secret, createdAt, expiresAt},
		"malformed email":        {id, "not-an-email", displayName, passwordHash, secret, createdAt, expiresAt},
		"blank display name":     {id, email, "", passwordHash, secret, createdAt, expiresAt},
		"blank password hash":    {id, email, displayName, "", secret, createdAt, expiresAt},
		"empty encrypted totp":   {id, email, displayName, passwordHash, nil, createdAt, expiresAt},
		"zero created at":        {id, email, displayName, passwordHash, secret, time.Time{}, expiresAt},
		"expires equal created":  {id, email, displayName, passwordHash, secret, createdAt, createdAt},
		"expires before created": {id, email, displayName, passwordHash, secret, createdAt, createdAt.Add(-time.Minute)},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if _, err := NewPendingEnrollment(tc.id, tc.email, tc.displayName, tc.passwordHash, tc.secret, tc.createdAt, tc.expiresAt); !errors.Is(err, ErrInvalidPendingEnrollment) {
				t.Fatalf("expected ErrInvalidPendingEnrollment, got %v", err)
			}
		})
	}
}

func TestPendingEnrollmentIsExpired(t *testing.T) {
	t.Parallel()

	id, email, displayName, passwordHash, secret, createdAt, expiresAt := validPendingEnrollmentArgs()
	enrollment, err := NewPendingEnrollment(id, email, displayName, passwordHash, secret, createdAt, expiresAt)
	if err != nil {
		t.Fatalf("valid pending enrollment rejected: %v", err)
	}

	if enrollment.IsExpired(expiresAt.Add(-time.Second)) {
		t.Fatal("enrollment must not be expired strictly before expires_at")
	}
	if !enrollment.IsExpired(expiresAt) {
		t.Fatal("enrollment must be expired at expires_at")
	}
	if !enrollment.IsExpired(expiresAt.Add(time.Second)) {
		t.Fatal("enrollment must be expired after expires_at")
	}
}
