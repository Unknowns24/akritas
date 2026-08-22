package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func validPendingEnrollmentArgs() (uuid.UUID, string, string, time.Time, time.Time) {
	createdAt := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	return uuid.New(), "admin@example.com", "Akritas Administrator", createdAt, createdAt.Add(10 * time.Minute)
}

func TestNewPendingEnrollmentValid(t *testing.T) {
	t.Parallel()

	id, email, displayName, createdAt, expiresAt := validPendingEnrollmentArgs()
	enrollment, err := NewPendingEnrollment(id, email, displayName, createdAt, expiresAt)
	if err != nil {
		t.Fatalf("valid pending enrollment rejected: %v", err)
	}
	if enrollment.ID != id || enrollment.Email != email || enrollment.DisplayName != displayName ||
		!enrollment.CreatedAt.Equal(createdAt) || !enrollment.ExpiresAt.Equal(expiresAt) {
		t.Fatal("constructed enrollment does not match input")
	}
}

func TestNewPendingEnrollmentRejectsInvalidFields(t *testing.T) {
	t.Parallel()

	id, email, displayName, createdAt, expiresAt := validPendingEnrollmentArgs()
	cases := map[string]struct {
		id          uuid.UUID
		email       string
		displayName string
		createdAt   time.Time
		expiresAt   time.Time
	}{
		"zero id":                {uuid.Nil, email, displayName, createdAt, expiresAt},
		"blank email":            {id, "", displayName, createdAt, expiresAt},
		"malformed email":        {id, "not-an-email", displayName, createdAt, expiresAt},
		"blank display name":     {id, email, "", createdAt, expiresAt},
		"zero created at":        {id, email, displayName, time.Time{}, expiresAt},
		"expires equal created":  {id, email, displayName, createdAt, createdAt},
		"expires before created": {id, email, displayName, createdAt, createdAt.Add(-time.Minute)},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if _, err := NewPendingEnrollment(tc.id, tc.email, tc.displayName, tc.createdAt, tc.expiresAt); !errors.Is(err, ErrInvalidPendingEnrollment) {
				t.Fatalf("expected ErrInvalidPendingEnrollment, got %v", err)
			}
		})
	}
}

func TestPendingEnrollmentIsExpired(t *testing.T) {
	t.Parallel()

	id, email, displayName, createdAt, expiresAt := validPendingEnrollmentArgs()
	enrollment, err := NewPendingEnrollment(id, email, displayName, createdAt, expiresAt)
	if err != nil {
		t.Fatalf("valid pending enrollment rejected: %v", err)
	}

	if enrollment.IsExpired(expiresAt.Add(-time.Second)) {
		t.Fatal("enrollment must not be expired strictly before expires_at")
	}
	if !enrollment.IsExpired(expiresAt) || !enrollment.IsExpired(expiresAt.Add(time.Second)) {
		t.Fatal("enrollment must be expired at and after expires_at")
	}
}
