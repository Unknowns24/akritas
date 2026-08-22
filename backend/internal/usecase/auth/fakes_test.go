package auth

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type fakeAdministratorRepository struct {
	exists bool
	err    error
	called bool
	order  *[]string
}

func (f *fakeAdministratorRepository) ExistsActive(ctx context.Context) (bool, error) {
	f.called = true
	if f.order != nil {
		*f.order = append(*f.order, "administrator_repository")
	}
	return f.exists, f.err
}

type fakePendingEnrollmentRepository struct {
	saved  *domain.PendingEnrollment
	err    error
	called bool
}

func (f *fakePendingEnrollmentRepository) Save(ctx context.Context, enrollment *domain.PendingEnrollment) error {
	f.called = true
	if f.err != nil {
		return f.err
	}
	f.saved = enrollment
	return nil
}

type fakeCredentialStore struct {
	encrypted []byte
	err       error
	called    bool
	plaintext []byte
}

func (f *fakeCredentialStore) Encrypt(ctx context.Context, plaintext []byte) ([]byte, error) {
	f.called = true
	f.plaintext = plaintext
	if f.err != nil {
		return nil, f.err
	}
	if f.encrypted != nil {
		return f.encrypted, nil
	}
	return []byte("encrypted:" + string(plaintext)), nil
}

type fakeTOTPSecretGenerator struct {
	secret out.TOTPSecret
	err    error
	called bool
}

func (f *fakeTOTPSecretGenerator) Generate(issuer, accountEmail string) (out.TOTPSecret, error) {
	f.called = true
	if f.err != nil {
		return out.TOTPSecret{}, f.err
	}
	return f.secret, nil
}

type fakePasswordHasher struct {
	hash   string
	err    error
	called bool
}

func (f *fakePasswordHasher) Hash(password string) (string, error) {
	f.called = true
	if f.err != nil {
		return "", f.err
	}
	return f.hash, nil
}

type fakeBootstrapTokenVerifier struct {
	valid  bool
	called bool
	order  *[]string
}

func (f *fakeBootstrapTokenVerifier) Verify(candidate string) bool {
	f.called = true
	if f.order != nil {
		*f.order = append(*f.order, "bootstrap_token_verifier")
	}
	return f.valid
}

type fakeRateLimiter struct {
	allowed bool
	err     error
	called  bool
	order   *[]string
}

func (f *fakeRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	f.called = true
	if f.order != nil {
		*f.order = append(*f.order, "rate_limiter")
	}
	return f.allowed, f.err
}

type fakeClock struct {
	now time.Time
}

func (f *fakeClock) Now() time.Time {
	return f.now
}
