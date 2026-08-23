package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func recoveryAdministrator(t *testing.T, now time.Time) *out.AdministratorAuthentication {
	t.Helper()
	administrator, err := domain.NewAdministrator(uuid.New(), "admin@example.com", "Admin", now.Add(-time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	return &out.AdministratorAuthentication{Administrator: *administrator, PasswordHash: "old-password-hash"}
}

func TestStartAdministratorRecoveryCreatesPendingRotationWithoutChangingActiveCredentials(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	administrators := &fakeAdministratorRepository{findByEmailResult: recoveryAdministrator(t, now)}
	pending := &fakePendingEnrollmentRepository{}
	credentials := &fakeCredentialStore{}
	uc := NewStartAdministratorRecoveryUseCase(
		administrators, pending, credentials,
		&fakeTOTPSecretGenerator{secret: out.TOTPSecret{Base32Key: "NEWSEED", OtpauthURI: "otpauth://totp/Akritas"}},
		&fakePasswordHasher{hash: "new-password-hash"}, &fakeBootstrapTokenVerifier{valid: true},
		&fakeRateLimiter{allowed: true}, &fakeTransactor{}, func() time.Time { return now }, "dummy-hash",
	)

	output, err := uc.Execute(context.Background(), in.StartAdministratorRecoveryInput{Email: "admin@example.com", NewPassword: "a-new-long-password", BootstrapToken: "bootstrap", RateLimitKey: "203.0.113.1"})
	if err != nil {
		t.Fatal(err)
	}
	if output.EnrollmentID == uuid.Nil || pending.passwordHash != "new-password-hash" || string(credentials.plaintext) != "NEWSEED" {
		t.Fatalf("pending recovery not persisted safely: output=%+v hash=%q seed=%q", output, pending.passwordHash, credentials.plaintext)
	}
	if administrators.rotateCalled {
		t.Fatal("start recovery must not rotate active credentials")
	}
}

func TestStartAdministratorRecoveryCredentialFailuresAreGeneric(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	cases := map[string]struct {
		administrator *out.AdministratorAuthentication
		bootstrapOK   bool
		samePassword  bool
	}{
		"unknown email":     {bootstrapOK: true},
		"invalid bootstrap": {administrator: recoveryAdministrator(t, now)},
		"same password":     {administrator: recoveryAdministrator(t, now), bootstrapOK: true, samePassword: true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			uc := NewStartAdministratorRecoveryUseCase(
				&fakeAdministratorRepository{findByEmailResult: tc.administrator}, &fakePendingEnrollmentRepository{}, &fakeCredentialStore{},
				&fakeTOTPSecretGenerator{}, &fakePasswordHasher{verifyResult: tc.samePassword}, &fakeBootstrapTokenVerifier{valid: tc.bootstrapOK},
				&fakeRateLimiter{allowed: true}, &fakeTransactor{}, func() time.Time { return now }, "dummy-hash",
			)
			_, err := uc.Execute(context.Background(), in.StartAdministratorRecoveryInput{Email: "admin@example.com", NewPassword: "a-new-long-password", BootstrapToken: "candidate", RateLimitKey: "ip"})
			if !errors.Is(err, domain.ErrInvalidCredentials) {
				t.Fatalf("error = %v, want generic invalid credentials", err)
			}
		})
	}
}
