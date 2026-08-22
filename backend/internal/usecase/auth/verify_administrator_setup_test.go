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

const (
	verifySessionIdleTTL     = 12 * time.Hour
	verifySessionAbsoluteTTL = 168 * time.Hour
)

type verifyDeps struct {
	rateLimiter        *fakeRateLimiter
	pendingEnrollments *fakePendingEnrollmentRepository
	credentialStore    *fakeCredentialStore
	totpVerifier       *fakeTOTPVerifier
	administrators     *fakeAdministratorRepository
	sessionTokens      *fakeSessionTokenGenerator
	sessions           *fakeAdministratorSessionRepository
	transactor         *fakeTransactor
	now                time.Time
}

func newHappyVerifyDeps(t *testing.T) *verifyDeps {
	t.Helper()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	enrollment, err := domain.NewPendingEnrollment(uuid.New(), "admin@example.com", "Akritas Administrator", now.Add(-time.Minute), now.Add(9*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	return &verifyDeps{
		rateLimiter:        &fakeRateLimiter{allowed: true},
		pendingEnrollments: &fakePendingEnrollmentRepository{findResult: &out.PendingEnrollmentAuthentication{Enrollment: *enrollment, PasswordHash: "argon2id-hash"}},
		credentialStore:    &fakeCredentialStore{getResult: []byte("JBSWY3DPEHPK3PXP")},
		totpVerifier:       &fakeTOTPVerifier{valid: true},
		administrators:     &fakeAdministratorRepository{},
		sessionTokens:      &fakeSessionTokenGenerator{token: "raw-token", hash: "hashed-token"},
		sessions:           &fakeAdministratorSessionRepository{},
		transactor:         &fakeTransactor{},
		now:                now,
	}
}

func (d *verifyDeps) usecase() in.VerifyAdministratorSetupUseCase {
	return NewVerifyAdministratorSetupUseCase(
		d.rateLimiter, d.pendingEnrollments, d.credentialStore, d.totpVerifier,
		d.administrators, d.sessionTokens, d.sessions, d.transactor, func() time.Time { return d.now },
		verifySessionIdleTTL, verifySessionAbsoluteTTL,
	)
}

func (d *verifyDeps) input() in.VerifyAdministratorSetupInput {
	return in.VerifyAdministratorSetupInput{EnrollmentID: d.pendingEnrollments.findResult.Enrollment.ID.String(), TOTPCode: "123456", RateLimitKey: "203.0.113.10"}
}

func TestVerifyAdministratorSetupHappyPathIsAtomic(t *testing.T) {
	t.Parallel()
	deps := newHappyVerifyDeps(t)
	output, err := deps.usecase().Execute(context.Background(), deps.input())
	if err != nil {
		t.Fatal(err)
	}
	if !deps.transactor.called || deps.administrators.createdAdministrator == nil || deps.sessions.savedSession == nil || !deps.credentialStore.moved || !deps.pendingEnrollments.deleteCalled {
		t.Fatal("activation must create administrator/session, move TOTP and consume enrollment in the transaction")
	}
	if deps.administrators.createdPasswordHash != "argon2id-hash" || output.SessionToken != "raw-token" {
		t.Fatal("activation must transfer the password hash and return only the raw cookie token")
	}
}

func TestVerifyAdministratorSetupRejectsInvalidEnrollmentAndCodeGenerically(t *testing.T) {
	t.Parallel()
	deps := newHappyVerifyDeps(t)
	deps.totpVerifier.valid = false
	_, err := deps.usecase().Execute(context.Background(), deps.input())
	if !errors.Is(err, domain.ErrInvalidTotpEnrollmentVerification) || deps.transactor.called {
		t.Fatalf("wrong code must fail before transaction with generic error: %v", err)
	}

	deps = newHappyVerifyDeps(t)
	deps.pendingEnrollments.findResult = nil
	_, err = deps.usecase().Execute(context.Background(), in.VerifyAdministratorSetupInput{EnrollmentID: uuid.NewString(), TOTPCode: "123456", RateLimitKey: "ip"})
	if !errors.Is(err, domain.ErrInvalidTotpEnrollmentVerification) {
		t.Fatalf("missing enrollment must use same generic error: %v", err)
	}
}

func TestVerifyAdministratorSetupRateLimited(t *testing.T) {
	t.Parallel()
	deps := newHappyVerifyDeps(t)
	deps.rateLimiter.allowed = false
	_, err := deps.usecase().Execute(context.Background(), deps.input())
	if !errors.Is(err, domain.ErrAuthenticationRateLimited) || deps.pendingEnrollments.findCalled {
		t.Fatalf("rate limit must stop all downstream work: %v", err)
	}
}

func TestVerifyAdministratorSetupRollsErrorFromAtomicSteps(t *testing.T) {
	t.Parallel()
	deps := newHappyVerifyDeps(t)
	want := errors.New("session unavailable")
	deps.sessions.err = want
	_, err := deps.usecase().Execute(context.Background(), deps.input())
	if !errors.Is(err, want) || deps.pendingEnrollments.deleteCalled || deps.credentialStore.moved {
		t.Fatalf("failed session write must stop remaining atomic steps: %v", err)
	}
}
