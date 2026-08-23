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

type setupDeps struct {
	administrators     *fakeAdministratorRepository
	pendingEnrollments *fakePendingEnrollmentRepository
	credentialStore    *fakeCredentialStore
	totpGenerator      *fakeTOTPSecretGenerator
	passwordHasher     *fakePasswordHasher
	bootstrapTokens    *fakeBootstrapTokenVerifier
	rateLimiter        *fakeRateLimiter
	transactor         *fakeTransactor
	clock              *fakeClock
}

func newHappyDeps() *setupDeps {
	return &setupDeps{
		administrators:     &fakeAdministratorRepository{exists: false},
		pendingEnrollments: &fakePendingEnrollmentRepository{},
		credentialStore:    &fakeCredentialStore{},
		totpGenerator: &fakeTOTPSecretGenerator{secret: out.TOTPSecret{
			Base32Key:  "JBSWY3DPEHPK3PXP",
			OtpauthURI: "otpauth://totp/Akritas:admin@example.com?secret=JBSWY3DPEHPK3PXP&issuer=Akritas",
		}},
		passwordHasher:  &fakePasswordHasher{hash: "$argon2id$v=19$m=19456,t=2,p=1$salt$hash"},
		bootstrapTokens: &fakeBootstrapTokenVerifier{valid: true},
		rateLimiter:     &fakeRateLimiter{allowed: true},
		transactor:      &fakeTransactor{},
		clock:           &fakeClock{now: time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)},
	}
}

func (d *setupDeps) usecase() in.StartAdministratorSetupUseCase {
	return NewStartAdministratorSetupUseCase(
		d.administrators, d.pendingEnrollments, d.credentialStore, d.totpGenerator,
		d.passwordHasher, d.bootstrapTokens, d.rateLimiter, d.transactor, d.clock.Now,
	)
}

func validSetupInput() in.StartAdministratorSetupInput {
	return in.StartAdministratorSetupInput{
		Email:          "admin@example.com",
		DisplayName:    "Akritas Administrator",
		Password:       "a-long-password-from-a-password-manager",
		BootstrapToken: "deployment-bootstrap-secret-not-a-totp-seed",
		RateLimitKey:   "203.0.113.10",
	}
}

func TestStartAdministratorSetupHappyPath(t *testing.T) {
	t.Parallel()

	deps := newHappyDeps()
	uc := deps.usecase()

	output, err := uc.Execute(context.Background(), validSetupInput())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.OtpauthURI != deps.totpGenerator.secret.OtpauthURI {
		t.Fatalf("otpauth_uri mismatch: got %q", output.OtpauthURI)
	}
	if output.ManualEntryKey != deps.totpGenerator.secret.Base32Key {
		t.Fatalf("manual_entry_key mismatch: got %q", output.ManualEntryKey)
	}
	wantExpiresAt := deps.clock.now.Add(pendingEnrollmentTTL)
	if !output.ExpiresAt.Equal(wantExpiresAt) {
		t.Fatalf("expires_at = %v, want %v", output.ExpiresAt, wantExpiresAt)
	}
	if output.EnrollmentID == uuid.Nil {
		t.Fatal("enrollment_id must not be zero")
	}

	saved := deps.pendingEnrollments.saved
	if saved == nil {
		t.Fatal("pending enrollment was not saved")
	}
	if saved.ID != output.EnrollmentID {
		t.Fatal("saved enrollment id must match returned enrollment_id")
	}
	if deps.pendingEnrollments.passwordHash != deps.passwordHasher.hash {
		t.Fatal("saved enrollment must carry the hashed password")
	}

	if string(deps.credentialStore.plaintext) != deps.totpGenerator.secret.Base32Key {
		t.Fatalf("credential store must encrypt the base32 secret, got %q", deps.credentialStore.plaintext)
	}
}

func TestStartAdministratorSetupRateLimited(t *testing.T) {
	t.Parallel()

	deps := newHappyDeps()
	deps.rateLimiter.allowed = false
	uc := deps.usecase()

	if _, err := uc.Execute(context.Background(), validSetupInput()); !errors.Is(err, domain.ErrAuthenticationRateLimited) {
		t.Fatalf("expected ErrAuthenticationRateLimited, got %v", err)
	}
	if deps.bootstrapTokens.called || deps.administrators.called || deps.passwordHasher.called ||
		deps.totpGenerator.called || deps.credentialStore.called || deps.pendingEnrollments.called {
		t.Fatal("no downstream dependency should be called when rate limited")
	}
}

func TestStartAdministratorSetupInvalidBootstrapToken(t *testing.T) {
	t.Parallel()

	deps := newHappyDeps()
	deps.bootstrapTokens.valid = false
	uc := deps.usecase()

	if _, err := uc.Execute(context.Background(), validSetupInput()); !errors.Is(err, domain.ErrInvalidBootstrapToken) {
		t.Fatalf("expected ErrInvalidBootstrapToken, got %v", err)
	}
	if deps.administrators.called || deps.passwordHasher.called || deps.totpGenerator.called ||
		deps.credentialStore.called || deps.pendingEnrollments.called {
		t.Fatal("no dependency past bootstrap token verification should be called")
	}
}

func TestStartAdministratorSetupAdministratorAlreadyExists(t *testing.T) {
	t.Parallel()

	deps := newHappyDeps()
	deps.administrators.exists = true
	uc := deps.usecase()

	if _, err := uc.Execute(context.Background(), validSetupInput()); !errors.Is(err, domain.ErrAdministratorAlreadyExists) {
		t.Fatalf("expected ErrAdministratorAlreadyExists, got %v", err)
	}
	if deps.passwordHasher.called || deps.totpGenerator.called || deps.credentialStore.called || deps.pendingEnrollments.called {
		t.Fatal("no secret-generating dependency should be called once an administrator already exists")
	}
}

func TestStartAdministratorSetupVerificationOrder(t *testing.T) {
	t.Parallel()

	var order []string
	deps := newHappyDeps()
	deps.rateLimiter.order = &order
	deps.bootstrapTokens.order = &order
	deps.administrators.order = &order
	uc := deps.usecase()

	if _, err := uc.Execute(context.Background(), validSetupInput()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"rate_limiter", "bootstrap_token_verifier", "administrator_repository"}
	if len(order) != len(want) {
		t.Fatalf("call order = %v, want %v", order, want)
	}
	for i, step := range want {
		if order[i] != step {
			t.Fatalf("call order = %v, want %v", order, want)
		}
	}
}

func TestStartAdministratorSetupDownstreamErrorsPropagateAndDoNotPersist(t *testing.T) {
	t.Parallel()

	cases := map[string]func(*setupDeps) error{
		"password hasher fails": func(d *setupDeps) error {
			err := errors.New("hasher unavailable")
			d.passwordHasher.err = err
			return err
		},
		"totp generator fails": func(d *setupDeps) error {
			err := errors.New("totp generator unavailable")
			d.totpGenerator.err = err
			return err
		},
		"credential store fails": func(d *setupDeps) error {
			err := errors.New("credential store unavailable")
			d.credentialStore.err = err
			return err
		},
		"pending enrollment save fails": func(d *setupDeps) error {
			err := errors.New("database unavailable")
			d.pendingEnrollments.err = err
			return err
		},
	}

	for name, arrange := range cases {
		arrange := arrange
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			deps := newHappyDeps()
			wantErr := arrange(deps)
			uc := deps.usecase()

			if _, err := uc.Execute(context.Background(), validSetupInput()); !errors.Is(err, wantErr) {
				t.Fatalf("expected %v, got %v", wantErr, err)
			}
			if deps.pendingEnrollments.saved != nil && !deps.transactor.called {
				t.Fatal("no pending enrollment must be persisted when a dependency fails")
			}
		})
	}
}
