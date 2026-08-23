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

func pendingRecovery(t *testing.T, now time.Time) *out.PendingEnrollmentAuthentication {
	t.Helper()
	enrollment, err := domain.NewPendingEnrollment(uuid.New(), "admin@example.com", "Admin", now.Add(-time.Minute), now.Add(9*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	return &out.PendingEnrollmentAuthentication{Enrollment: *enrollment, PasswordHash: "new-password-hash"}
}

func TestVerifyAdministratorRecoveryAtomicallyRotatesAndCreatesFreshSession(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	enrollment := pendingRecovery(t, now)
	administrators := &fakeAdministratorRepository{findByEmailResult: recoveryAdministrator(t, now), rotateResult: true}
	pending := &fakePendingEnrollmentRepository{findResult: enrollment}
	sessions := &fakeAdministratorSessionRepository{}
	credentials := &fakeCredentialStore{getResult: []byte("new-totp-secret")}
	uc := NewVerifyAdministratorRecoveryUseCase(
		&fakeRateLimiter{allowed: true}, pending, credentials, &fakeTOTPVerifier{valid: true, period: 123},
		administrators, &fakeSessionTokenGenerator{token: "fresh-token", hash: "fresh-hash"}, sessions,
		&fakeTransactor{}, func() time.Time { return now }, 12*time.Hour, 7*24*time.Hour,
	)

	output, err := uc.Execute(context.Background(), in.VerifyAdministratorRecoveryInput{EnrollmentID: enrollment.Enrollment.ID.String(), TOTPCode: "123456", RateLimitKey: "203.0.113.1"})
	if err != nil {
		t.Fatal(err)
	}
	if !administrators.rotateCalled || administrators.rotatedPasswordHash != "new-password-hash" || !credentials.moved || !sessions.revokeAllCalled || !sessions.called {
		t.Fatalf("rotation not completed: rotate=%v move=%v revokeAll=%v save=%v", administrators.rotateCalled, credentials.moved, sessions.revokeAllCalled, sessions.called)
	}
	if output.SessionToken != "fresh-token" || output.Administrator.LastAcceptedTOTPPeriod != 123 {
		t.Fatalf("unexpected output: %+v", output)
	}
}

func TestVerifyAdministratorRecoveryRejectsInvalidEnrollmentOrCodeGenerically(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	for name, pending := range map[string]*out.PendingEnrollmentAuthentication{"unknown": nil, "wrong code": pendingRecovery(t, now)} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			uc := NewVerifyAdministratorRecoveryUseCase(
				&fakeRateLimiter{allowed: true}, &fakePendingEnrollmentRepository{findResult: pending}, &fakeCredentialStore{getResult: []byte("seed")},
				&fakeTOTPVerifier{valid: false}, &fakeAdministratorRepository{}, &fakeSessionTokenGenerator{}, &fakeAdministratorSessionRepository{},
				&fakeTransactor{}, func() time.Time { return now }, time.Hour, 2*time.Hour,
			)
			id := uuid.NewString()
			if pending != nil {
				id = pending.Enrollment.ID.String()
			}
			_, err := uc.Execute(context.Background(), in.VerifyAdministratorRecoveryInput{EnrollmentID: id, TOTPCode: "000000", RateLimitKey: "ip"})
			if !errors.Is(err, domain.ErrInvalidCredentials) {
				t.Fatalf("error = %v, want generic invalid credentials", err)
			}
		})
	}
}

func TestVerifyAdministratorRecoveryPropagatesEveryTransactionalFailure(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	failure := errors.New("injected transactional failure")
	tests := map[string]func(*fakePendingEnrollmentRepository, *fakeAdministratorRepository, *fakeCredentialStore, *fakeAdministratorSessionRepository){
		"consume enrollment": func(pending *fakePendingEnrollmentRepository, _ *fakeAdministratorRepository, _ *fakeCredentialStore, _ *fakeAdministratorSessionRepository) {
			pending.deleteErr = failure
		},
		"rotate password": func(_ *fakePendingEnrollmentRepository, administrators *fakeAdministratorRepository, _ *fakeCredentialStore, _ *fakeAdministratorSessionRepository) {
			administrators.rotateErr = failure
		},
		"replace TOTP": func(_ *fakePendingEnrollmentRepository, _ *fakeAdministratorRepository, credentials *fakeCredentialStore, _ *fakeAdministratorSessionRepository) {
			credentials.err = failure
		},
		"revoke sessions": func(_ *fakePendingEnrollmentRepository, _ *fakeAdministratorRepository, _ *fakeCredentialStore, sessions *fakeAdministratorSessionRepository) {
			sessions.revokeAllErr = failure
		},
		"save fresh session": func(_ *fakePendingEnrollmentRepository, _ *fakeAdministratorRepository, _ *fakeCredentialStore, sessions *fakeAdministratorSessionRepository) {
			sessions.err = failure
		},
	}
	for name, inject := range tests {
		name, inject := name, inject
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			enrollment := pendingRecovery(t, now)
			pending := &fakePendingEnrollmentRepository{findResult: enrollment}
			administrators := &fakeAdministratorRepository{findByEmailResult: recoveryAdministrator(t, now), rotateResult: true}
			credentials := &fakeCredentialStore{getResult: []byte("new-totp-secret")}
			sessions := &fakeAdministratorSessionRepository{}
			inject(pending, administrators, credentials, sessions)
			uc := NewVerifyAdministratorRecoveryUseCase(
				&fakeRateLimiter{allowed: true}, pending, credentials, &fakeTOTPVerifier{valid: true, period: 123},
				administrators, &fakeSessionTokenGenerator{token: "fresh-token", hash: "fresh-hash"}, sessions,
				&fakeTransactor{}, func() time.Time { return now }, 12*time.Hour, 7*24*time.Hour,
			)
			_, err := uc.Execute(context.Background(), in.VerifyAdministratorRecoveryInput{EnrollmentID: enrollment.Enrollment.ID.String(), TOTPCode: "123456", RateLimitKey: "ip"})
			if !errors.Is(err, failure) {
				t.Fatalf("error=%v want injected failure", err)
			}
		})
	}
}
