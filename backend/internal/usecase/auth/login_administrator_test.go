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

type loginDeps struct {
	rateLimiter    *fakeRateLimiter
	administrators *fakeAdministratorRepository
	passwordHasher *fakePasswordHasher
	credentials    *fakeCredentialStore
	totpVerifier   *fakeTOTPVerifier
	tokens         *fakeSessionTokenGenerator
	sessions       *fakeAdministratorSessionRepository
	transactor     *fakeTransactor
	now            time.Time
}

func newLoginDeps(t *testing.T) loginDeps {
	t.Helper()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	administrator, err := domain.NewAdministrator(uuid.New(), "admin@example.com", "Admin", now.Add(-time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	return loginDeps{
		rateLimiter: &fakeRateLimiter{allowed: true},
		administrators: &fakeAdministratorRepository{findByEmailResult: &out.AdministratorAuthentication{
			Administrator: *administrator,
			PasswordHash:  "password-hash",
		}},
		passwordHasher: &fakePasswordHasher{verifyResult: true},
		credentials:    &fakeCredentialStore{getResult: []byte("totp-secret")},
		totpVerifier:   &fakeTOTPVerifier{valid: true, period: 100},
		tokens:         &fakeSessionTokenGenerator{token: "opaque-token", hash: "token-hash"},
		sessions:       &fakeAdministratorSessionRepository{},
		transactor:     &fakeTransactor{},
		now:            now,
	}
}

func (d loginDeps) useCase() in.LoginAdministratorUseCase {
	return NewLoginAdministratorUseCase(
		d.rateLimiter, d.administrators, d.passwordHasher, d.credentials,
		d.totpVerifier, d.tokens, d.sessions, d.transactor, func() time.Time { return d.now },
		12*time.Hour, 7*24*time.Hour,
		"dummy-password-hash",
	)
}

func validLoginInput() in.LoginAdministratorInput {
	return in.LoginAdministratorInput{
		Email: "admin@example.com", Password: "correct-password", TOTPCode: "123456", RateLimitKey: "203.0.113.1",
	}
}

func TestLoginAdministratorCreatesSessionAndConsumesTOTPPeriodAtomically(t *testing.T) {
	t.Parallel()
	d := newLoginDeps(t)

	output, err := d.useCase().Execute(context.Background(), validLoginInput())
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if output.SessionToken != "opaque-token" || !d.transactor.called || !d.sessions.called {
		t.Fatal("expected token and transactional session creation")
	}
	if !d.administrators.updatePeriodCalled || d.administrators.updatedPeriod != 100 {
		t.Fatal("expected TOTP period to be consumed")
	}
	if d.credentials.ownerType != out.CredentialOwnerAdministrator || d.credentials.kind != out.SecretKindAdministratorTOTP {
		t.Fatal("expected administrator TOTP credential lookup")
	}
}

func TestLoginAdministratorRejectsReplayedOrConcurrentlyConsumedPeriod(t *testing.T) {
	t.Parallel()

	t.Run("already accepted before transaction", func(t *testing.T) {
		d := newLoginDeps(t)
		d.administrators.findByEmailResult.Administrator.LastAcceptedTOTPPeriod = 100
		_, err := d.useCase().Execute(context.Background(), validLoginInput())
		if !errors.Is(err, domain.ErrInvalidCredentials) || d.transactor.called {
			t.Fatalf("expected pre-transaction replay rejection, got %v", err)
		}
	})

	t.Run("lost compare-and-set race", func(t *testing.T) {
		d := newLoginDeps(t)
		d.administrators.consumePeriodReject = true
		_, err := d.useCase().Execute(context.Background(), validLoginInput())
		if !errors.Is(err, domain.ErrInvalidCredentials) || d.sessions.called {
			t.Fatalf("expected concurrent replay rejection without session, got %v", err)
		}
	})
}

func TestLoginAdministratorUsesGenericCredentialFailures(t *testing.T) {
	t.Parallel()

	cases := map[string]func(*loginDeps){
		"unknown account": func(d *loginDeps) { d.administrators.findByEmailResult = nil },
		"bad password":    func(d *loginDeps) { d.passwordHasher.verifyResult = false },
		"missing totp":    func(d *loginDeps) { d.credentials.getErr = errors.New("missing") },
		"bad totp":        func(d *loginDeps) { d.totpVerifier.valid = false },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			d := newLoginDeps(t)
			mutate(&d)
			_, err := d.useCase().Execute(context.Background(), validLoginInput())
			if !errors.Is(err, domain.ErrInvalidCredentials) {
				t.Fatalf("expected generic invalid credentials, got %v", err)
			}
		})
	}
}

func TestLoginAdministratorHasIndependentRateLimitScopes(t *testing.T) {
	t.Parallel()
	d := newLoginDeps(t)
	d.rateLimiter.allowed = false

	_, err := d.useCase().Execute(context.Background(), validLoginInput())
	if !errors.Is(err, domain.ErrAuthenticationRateLimited) || d.administrators.findByEmailCalled {
		t.Fatalf("expected rate limit before credential lookup, got %v", err)
	}
}

func TestLoginAdministratorDoesNotPersistSessionWhenTransactionFails(t *testing.T) {
	t.Parallel()
	d := newLoginDeps(t)
	d.transactor.err = errors.New("transaction unavailable")

	_, err := d.useCase().Execute(context.Background(), validLoginInput())
	if err == nil || d.sessions.called {
		t.Fatalf("expected failed transaction without session, got %v", err)
	}
}
