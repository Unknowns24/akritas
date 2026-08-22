package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
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
	clock              *fakeClock
}

func newValidPendingEnrollment(t *testing.T, now time.Time) *domain.PendingEnrollment {
	t.Helper()
	enrollment, err := domain.NewPendingEnrollment(
		uuid.New(), "admin@example.com", "Akritas Administrator", "argon2id-hash",
		[]byte("encrypted-secret"), now.Add(-time.Minute), now.Add(9*time.Minute),
	)
	if err != nil {
		t.Fatalf("build pending enrollment: %v", err)
	}
	return enrollment
}

func newHappyVerifyDeps(t *testing.T) (*verifyDeps, time.Time) {
	t.Helper()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	enrollment := newValidPendingEnrollment(t, now)
	return &verifyDeps{
		rateLimiter:        &fakeRateLimiter{allowed: true},
		pendingEnrollments: &fakePendingEnrollmentRepository{findResult: enrollment},
		credentialStore:    &fakeCredentialStore{decryptResult: []byte("JBSWY3DPEHPK3PXP")},
		totpVerifier:       &fakeTOTPVerifier{valid: true},
		administrators:     &fakeAdministratorRepository{exists: false},
		sessionTokens:      &fakeSessionTokenGenerator{token: "raw-token", hash: "hashed-token"},
		sessions:           &fakeAdministratorSessionRepository{},
		transactor:         &fakeTransactor{},
		clock:              &fakeClock{now: now},
	}, now
}

func (d *verifyDeps) usecase() in.VerifyAdministratorSetupUseCase {
	return NewVerifyAdministratorSetupUseCase(
		d.rateLimiter, d.pendingEnrollments, d.credentialStore, d.totpVerifier,
		d.administrators, d.sessionTokens, d.sessions, d.transactor, d.clock,
		verifySessionIdleTTL, verifySessionAbsoluteTTL,
	)
}

func validVerifyInput(enrollmentID uuid.UUID) in.VerifyAdministratorSetupInput {
	return in.VerifyAdministratorSetupInput{
		EnrollmentID: enrollmentID.String(),
		TOTPCode:     "123456",
		RateLimitKey: "203.0.113.10",
	}
}

func TestVerifyAdministratorSetupHappyPath(t *testing.T) {
	t.Parallel()

	deps, now := newHappyVerifyDeps(t)
	enrollment := deps.pendingEnrollments.findResult
	uc := deps.usecase()

	output, err := uc.Execute(context.Background(), validVerifyInput(enrollment.ID))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(deps.credentialStore.decryptedCiphertext) != string(enrollment.EncryptedTOTPSecret) {
		t.Fatal("Decrypt must be called with the enrollment's encrypted secret")
	}
	if deps.totpVerifier.secretArg != string(deps.credentialStore.decryptResult) || deps.totpVerifier.codeArg != "123456" {
		t.Fatal("TOTPVerifier must be called with the decrypted secret and the submitted code")
	}

	admin := deps.administrators.createdAdministrator
	if admin == nil {
		t.Fatal("Administrator was not created")
	}
	if admin.Email != enrollment.Email || admin.DisplayName != enrollment.DisplayName {
		t.Fatalf("created administrator does not match the enrollment: %+v", admin)
	}
	if deps.administrators.createdPasswordHash != enrollment.PasswordHash {
		t.Fatal("Create must receive the enrollment's password hash")
	}
	if string(deps.administrators.createdEncryptedTOTPSecret) != string(enrollment.EncryptedTOTPSecret) {
		t.Fatal("Create must receive the SAME ciphertext already on the enrollment, not a freshly re-encrypted one")
	}

	if !deps.transactor.called {
		t.Fatal("Transactor.WithinTransaction must be called")
	}
	if deps.sessions.savedTokenHash != "hashed-token" {
		t.Fatal("Save must receive the hash returned by the session token generator")
	}
	if deps.sessions.savedSession == nil || deps.sessions.savedSession.AdministratorID != admin.ID {
		t.Fatal("saved session must reference the newly created administrator")
	}

	if !deps.pendingEnrollments.deleteCalled || deps.pendingEnrollments.deletedID != enrollment.ID {
		t.Fatal("the consumed enrollment must be deleted")
	}

	if output.SessionToken != "raw-token" {
		t.Fatal("Output must expose the raw session token, never the hash")
	}
	if !output.AuthenticatedAt.Equal(now) {
		t.Fatalf("AuthenticatedAt = %v, want %v", output.AuthenticatedAt, now)
	}
	if !output.IdleExpiresAt.Equal(now.Add(verifySessionIdleTTL)) {
		t.Fatalf("IdleExpiresAt = %v, want %v", output.IdleExpiresAt, now.Add(verifySessionIdleTTL))
	}
	if !output.AbsoluteExpiresAt.Equal(now.Add(verifySessionAbsoluteTTL)) {
		t.Fatalf("AbsoluteExpiresAt = %v, want %v", output.AbsoluteExpiresAt, now.Add(verifySessionAbsoluteTTL))
	}
	if output.Administrator.Email != enrollment.Email {
		t.Fatal("Output.Administrator must reflect the created administrator")
	}
}

func TestVerifyAdministratorSetupRateLimited(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyVerifyDeps(t)
	deps.rateLimiter.allowed = false
	uc := deps.usecase()

	_, err := uc.Execute(context.Background(), validVerifyInput(deps.pendingEnrollments.findResult.ID))
	if !errors.Is(err, ErrSetupRateLimited) {
		t.Fatalf("expected ErrSetupRateLimited, got %v", err)
	}
	if deps.pendingEnrollments.findCalled || deps.credentialStore.decryptCalled || deps.totpVerifier.called ||
		deps.administrators.called || deps.transactor.called {
		t.Fatal("no downstream dependency should be called when rate limited")
	}
}

func TestVerifyAdministratorSetupMalformedEnrollmentID(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyVerifyDeps(t)
	uc := deps.usecase()

	input := in.VerifyAdministratorSetupInput{EnrollmentID: "not-a-uuid", TOTPCode: "123456", RateLimitKey: "203.0.113.10"}
	_, err := uc.Execute(context.Background(), input)
	if !errors.Is(err, domain.ErrInvalidTotpEnrollmentVerification) {
		t.Fatalf("expected ErrInvalidTotpEnrollmentVerification, got %v", err)
	}
	if deps.pendingEnrollments.findCalled {
		t.Fatal("FindByID must not be called for a malformed enrollment id")
	}
}

func TestVerifyAdministratorSetupEnrollmentNotFound(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyVerifyDeps(t)
	deps.pendingEnrollments.findResult = nil
	uc := deps.usecase()

	_, err := uc.Execute(context.Background(), validVerifyInput(uuid.New()))
	if !errors.Is(err, domain.ErrInvalidTotpEnrollmentVerification) {
		t.Fatalf("expected ErrInvalidTotpEnrollmentVerification, got %v", err)
	}
	if deps.credentialStore.decryptCalled || deps.totpVerifier.called {
		t.Fatal("no verification dependency should be called when the enrollment is not found")
	}
}

func TestVerifyAdministratorSetupExpiredEnrollment(t *testing.T) {
	t.Parallel()

	deps, now := newHappyVerifyDeps(t)
	expired, err := domain.NewPendingEnrollment(
		deps.pendingEnrollments.findResult.ID, "admin@example.com", "Akritas Administrator", "argon2id-hash",
		[]byte("encrypted-secret"), now.Add(-20*time.Minute), now.Add(-10*time.Minute),
	)
	if err != nil {
		t.Fatalf("build expired enrollment: %v", err)
	}
	deps.pendingEnrollments.findResult = expired
	uc := deps.usecase()

	_, err = uc.Execute(context.Background(), validVerifyInput(expired.ID))
	if !errors.Is(err, domain.ErrInvalidTotpEnrollmentVerification) {
		t.Fatalf("expected ErrInvalidTotpEnrollmentVerification, got %v", err)
	}
	if deps.credentialStore.decryptCalled || deps.totpVerifier.called {
		t.Fatal("no verification dependency should be called for an expired enrollment")
	}
}

func TestVerifyAdministratorSetupWrongCode(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyVerifyDeps(t)
	deps.totpVerifier.valid = false
	uc := deps.usecase()

	_, err := uc.Execute(context.Background(), validVerifyInput(deps.pendingEnrollments.findResult.ID))
	if !errors.Is(err, domain.ErrInvalidTotpEnrollmentVerification) {
		t.Fatalf("expected ErrInvalidTotpEnrollmentVerification, got %v", err)
	}
	if deps.administrators.called || deps.transactor.called {
		t.Fatal("no dependency past TOTP verification should be called on a wrong code")
	}
}

func TestVerifyAdministratorSetupAdministratorAlreadyExistsEarly(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyVerifyDeps(t)
	deps.administrators.exists = true
	uc := deps.usecase()

	_, err := uc.Execute(context.Background(), validVerifyInput(deps.pendingEnrollments.findResult.ID))
	if !errors.Is(err, domain.ErrAdministratorAlreadyExists) {
		t.Fatalf("expected ErrAdministratorAlreadyExists, got %v", err)
	}
	if deps.transactor.called {
		t.Fatal("Transactor must not be called once ExistsActive already reports true")
	}
}

func TestVerifyAdministratorSetupAdministratorAlreadyExistsRaceInsideTransaction(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyVerifyDeps(t)
	deps.administrators.createErr = domain.ErrAdministratorAlreadyExists
	uc := deps.usecase()

	_, err := uc.Execute(context.Background(), validVerifyInput(deps.pendingEnrollments.findResult.ID))
	if !errors.Is(err, domain.ErrAdministratorAlreadyExists) {
		t.Fatalf("expected ErrAdministratorAlreadyExists, got %v", err)
	}
	if deps.sessions.called {
		t.Fatal("Save must not be called when Create fails inside the transaction")
	}
	if deps.pendingEnrollments.deleteCalled {
		t.Fatal("the enrollment must not be consumed when activation failed")
	}
}

func TestVerifyAdministratorSetupInfrastructureErrorsPropagate(t *testing.T) {
	t.Parallel()

	cases := map[string]func(*verifyDeps) error{
		"decrypt fails": func(d *verifyDeps) error {
			err := errors.New("credential store unavailable")
			d.credentialStore.decryptErr = err
			return err
		},
		"totp verifier fails": func(d *verifyDeps) error {
			err := errors.New("totp verifier unavailable")
			d.totpVerifier.err = err
			return err
		},
		"exists active fails": func(d *verifyDeps) error {
			err := errors.New("database unavailable")
			d.administrators.err = err
			return err
		},
		"session token generator fails": func(d *verifyDeps) error {
			err := errors.New("rng unavailable")
			d.sessionTokens.err = err
			return err
		},
		"transactor fails": func(d *verifyDeps) error {
			err := errors.New("transaction could not start")
			d.transactor.err = err
			return err
		},
		"save fails inside transaction": func(d *verifyDeps) error {
			err := errors.New("database unavailable")
			d.sessions.err = err
			return err
		},
		"delete fails after commit": func(d *verifyDeps) error {
			err := errors.New("database unavailable")
			d.pendingEnrollments.deleteErr = err
			return err
		},
	}

	for name, arrange := range cases {
		arrange := arrange
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			deps, _ := newHappyVerifyDeps(t)
			wantErr := arrange(deps)
			uc := deps.usecase()

			_, err := uc.Execute(context.Background(), validVerifyInput(deps.pendingEnrollments.findResult.ID))
			if !errors.Is(err, wantErr) {
				t.Fatalf("expected %v, got %v", wantErr, err)
			}
		})
	}
}
