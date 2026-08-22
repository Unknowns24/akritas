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
	loginSessionIdleTTL     = 12 * time.Hour
	loginSessionAbsoluteTTL = 168 * time.Hour
)

type loginDeps struct {
	rateLimiter     *fakeRateLimiter
	administrators  *fakeAdministratorRepository
	passwordHasher  *fakePasswordHasher
	credentialStore *fakeCredentialStore
	totpVerifier    *fakeTOTPVerifier
	sessionTokens   *fakeSessionTokenGenerator
	sessions        *fakeAdministratorSessionRepository
	transactor      *fakeTransactor
	clock           *fakeClock
}

func newHappyLoginDeps(t *testing.T) (*loginDeps, time.Time) {
	t.Helper()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	admin, err := domain.NewAdministrator(uuid.New(), "admin@example.com", "Akritas Administrator", now.Add(-time.Hour))
	if err != nil {
		t.Fatalf("build administrator: %v", err)
	}
	return &loginDeps{
		rateLimiter: &fakeRateLimiter{allowed: true},
		administrators: &fakeAdministratorRepository{findByEmailResult: &out.AdministratorCredentials{
			Administrator:          *admin,
			PasswordHash:           "argon2id-hash",
			EncryptedTOTPSecret:    []byte("encrypted-secret"),
			LastAcceptedTOTPPeriod: 111,
		}},
		passwordHasher:  &fakePasswordHasher{verifyResult: true},
		credentialStore: &fakeCredentialStore{decryptResult: []byte("JBSWY3DPEHPK3PXP")},
		totpVerifier:    &fakeTOTPVerifier{valid: true, period: 222},
		sessionTokens:   &fakeSessionTokenGenerator{token: "raw-token", hash: "hashed-token"},
		sessions:        &fakeAdministratorSessionRepository{},
		transactor:      &fakeTransactor{},
		clock:           &fakeClock{now: now},
	}, now
}

func (d *loginDeps) usecase() in.LoginAdministratorUseCase {
	return NewLoginAdministratorUseCase(
		d.rateLimiter, d.administrators, d.passwordHasher, d.credentialStore, d.totpVerifier,
		d.sessionTokens, d.sessions, d.transactor, d.clock,
		loginSessionIdleTTL, loginSessionAbsoluteTTL,
	)
}

func validLoginInput() in.LoginAdministratorInput {
	return in.LoginAdministratorInput{
		Email:        "admin@example.com",
		Password:     "a-long-password-from-a-password-manager",
		TOTPCode:     "123456",
		RateLimitKey: "203.0.113.10",
	}
}

func TestLoginAdministratorHappyPath(t *testing.T) {
	t.Parallel()

	deps, now := newHappyLoginDeps(t)
	admin := deps.administrators.findByEmailResult.Administrator
	uc := deps.usecase()

	output, err := uc.Execute(context.Background(), validLoginInput())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Administrator.ID != admin.ID {
		t.Fatal("output must expose the authenticated administrator")
	}
	if output.SessionToken != "raw-token" {
		t.Fatal("output must expose the raw session token, never the hash")
	}
	if !output.AuthenticatedAt.Equal(now) {
		t.Fatalf("AuthenticatedAt = %v, want %v", output.AuthenticatedAt, now)
	}
	if !output.IdleExpiresAt.Equal(now.Add(loginSessionIdleTTL)) {
		t.Fatalf("IdleExpiresAt = %v, want %v", output.IdleExpiresAt, now.Add(loginSessionIdleTTL))
	}
	if !output.AbsoluteExpiresAt.Equal(now.Add(loginSessionAbsoluteTTL)) {
		t.Fatalf("AbsoluteExpiresAt = %v, want %v", output.AbsoluteExpiresAt, now.Add(loginSessionAbsoluteTTL))
	}

	if !deps.transactor.called {
		t.Fatal("Transactor.WithinTransaction must be called")
	}
	if deps.administrators.updatedPeriodID != admin.ID || deps.administrators.updatedPeriod != 222 {
		t.Fatalf("UpdateLastAcceptedTOTPPeriod called with wrong args: id=%v period=%d",
			deps.administrators.updatedPeriodID, deps.administrators.updatedPeriod)
	}
	if deps.sessions.savedTokenHash != "hashed-token" {
		t.Fatal("Save must receive the hash returned by the session token generator")
	}
	if deps.sessions.savedSession == nil || deps.sessions.savedSession.AdministratorID != admin.ID {
		t.Fatal("saved session must reference the authenticated administrator")
	}

	if string(deps.credentialStore.decryptedCiphertext) != "encrypted-secret" {
		t.Fatal("Decrypt must be called with the stored encrypted TOTP secret")
	}
	if deps.totpVerifier.secretArg != "JBSWY3DPEHPK3PXP" || deps.totpVerifier.codeArg != "123456" {
		t.Fatal("TOTPVerifier must be called with the decrypted secret and the submitted code")
	}
	if deps.passwordHasher.verifyArgs != [2]string{"a-long-password-from-a-password-manager", "argon2id-hash"} {
		t.Fatal("PasswordHasher.Verify must be called with the submitted password and the stored hash")
	}
}

func TestLoginAdministratorDoesNotRevokeOtherSessions(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyLoginDeps(t)
	uc := deps.usecase()

	if _, err := uc.Execute(context.Background(), validLoginInput()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if deps.sessions.revokeCalled {
		t.Fatal("login must never revoke any existing session (ADR-008 reserves that for recovery)")
	}
}

func TestLoginAdministratorIPRateLimited(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyLoginDeps(t)
	deps.rateLimiter.allowed = false
	uc := deps.usecase()

	if _, err := uc.Execute(context.Background(), validLoginInput()); !errors.Is(err, ErrLoginRateLimited) {
		t.Fatalf("expected ErrLoginRateLimited, got %v", err)
	}
	if deps.administrators.findByEmailCalled {
		t.Fatal("no lookup should happen once the IP rate limit is hit")
	}
}

func TestLoginAdministratorAccountRateLimited(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyLoginDeps(t)
	callCount := 0
	deps.rateLimiter.allowed = true
	// First Allow (IP) must pass, second (account) must fail. Use a custom
	// fake here since fakeRateLimiter only supports one static answer.
	uc := NewLoginAdministratorUseCase(
		&sequencedRateLimiter{results: []bool{true, false}, calls: &callCount},
		deps.administrators, deps.passwordHasher, deps.credentialStore, deps.totpVerifier,
		deps.sessionTokens, deps.sessions, deps.transactor, deps.clock,
		loginSessionIdleTTL, loginSessionAbsoluteTTL,
	)

	if _, err := uc.Execute(context.Background(), validLoginInput()); !errors.Is(err, ErrLoginRateLimited) {
		t.Fatalf("expected ErrLoginRateLimited, got %v", err)
	}
	if callCount != 2 {
		t.Fatalf("expected exactly 2 rate limiter calls (ip then account), got %d", callCount)
	}
	if deps.administrators.findByEmailCalled {
		t.Fatal("no lookup should happen once the account rate limit is hit")
	}
}

type sequencedRateLimiter struct {
	results []bool
	calls   *int
}

func (l *sequencedRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	i := *l.calls
	*l.calls++
	if i >= len(l.results) {
		return false, nil
	}
	return l.results[i], nil
}

func TestLoginAdministratorUnknownEmail(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyLoginDeps(t)
	deps.administrators.findByEmailResult = nil
	uc := deps.usecase()

	if _, err := uc.Execute(context.Background(), validLoginInput()); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
	if deps.passwordHasher.verifyCalled || deps.totpVerifier.called {
		t.Fatal("no further dependency should be called for an unknown email")
	}
}

func TestLoginAdministratorWrongPassword(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyLoginDeps(t)
	deps.passwordHasher.verifyResult = false
	uc := deps.usecase()

	if _, err := uc.Execute(context.Background(), validLoginInput()); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
	if deps.credentialStore.decryptCalled || deps.totpVerifier.called {
		t.Fatal("no further dependency should be called for a wrong password")
	}
}

func TestLoginAdministratorInvalidTOTPCode(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyLoginDeps(t)
	deps.totpVerifier.valid = false
	uc := deps.usecase()

	if _, err := uc.Execute(context.Background(), validLoginInput()); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
	if deps.transactor.called {
		t.Fatal("no transaction should start for an invalid TOTP code")
	}
}

func TestLoginAdministratorRejectsReusedPeriod(t *testing.T) {
	t.Parallel()

	deps, _ := newHappyLoginDeps(t)
	deps.totpVerifier.period = deps.administrators.findByEmailResult.LastAcceptedTOTPPeriod
	uc := deps.usecase()

	if _, err := uc.Execute(context.Background(), validLoginInput()); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
	if deps.transactor.called {
		t.Fatal("no transaction should start when the period was already accepted")
	}
}

func TestLoginAdministratorInfrastructureErrorsPropagate(t *testing.T) {
	t.Parallel()

	cases := map[string]func(*loginDeps) error{
		"find by email fails": func(d *loginDeps) error {
			err := errors.New("database unavailable")
			d.administrators.findByEmailErr = err
			return err
		},
		"decrypt fails": func(d *loginDeps) error {
			err := errors.New("credential store unavailable")
			d.credentialStore.decryptErr = err
			return err
		},
		"session token generator fails": func(d *loginDeps) error {
			err := errors.New("rng unavailable")
			d.sessionTokens.err = err
			return err
		},
		"transactor fails": func(d *loginDeps) error {
			err := errors.New("transaction could not start")
			d.transactor.err = err
			return err
		},
		"update period fails inside transaction": func(d *loginDeps) error {
			err := errors.New("database unavailable")
			d.administrators.updatePeriodErr = err
			return err
		},
		"save fails inside transaction": func(d *loginDeps) error {
			err := errors.New("database unavailable")
			d.sessions.err = err
			return err
		},
	}

	for name, arrange := range cases {
		arrange := arrange
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			deps, _ := newHappyLoginDeps(t)
			wantErr := arrange(deps)
			uc := deps.usecase()

			if _, err := uc.Execute(context.Background(), validLoginInput()); !errors.Is(err, wantErr) {
				t.Fatalf("expected %v, got %v", wantErr, err)
			}
		})
	}
}
