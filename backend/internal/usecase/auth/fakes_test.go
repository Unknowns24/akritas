package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type fakeAdministratorRepository struct {
	exists bool
	err    error
	called bool
	order  *[]string

	createErr                  error
	createCalled               bool
	createdAdministrator       *domain.Administrator
	createdPasswordHash        string
	createdEncryptedTOTPSecret []byte
}

func (f *fakeAdministratorRepository) ExistsActive(ctx context.Context) (bool, error) {
	f.called = true
	if f.order != nil {
		*f.order = append(*f.order, "administrator_repository")
	}
	return f.exists, f.err
}

func (f *fakeAdministratorRepository) Create(ctx context.Context, administrator *domain.Administrator, passwordHash string, encryptedTOTPSecret []byte) error {
	f.createCalled = true
	if f.order != nil {
		*f.order = append(*f.order, "administrator_repository.create")
	}
	if f.createErr != nil {
		return f.createErr
	}
	f.createdAdministrator = administrator
	f.createdPasswordHash = passwordHash
	f.createdEncryptedTOTPSecret = encryptedTOTPSecret
	return nil
}

type fakePendingEnrollmentRepository struct {
	saved  *domain.PendingEnrollment
	err    error
	called bool

	findResult *domain.PendingEnrollment
	findErr    error
	findCalled bool

	deleteErr    error
	deleteCalled bool
	deletedID    uuid.UUID
}

func (f *fakePendingEnrollmentRepository) Save(ctx context.Context, enrollment *domain.PendingEnrollment) error {
	f.called = true
	if f.err != nil {
		return f.err
	}
	f.saved = enrollment
	return nil
}

func (f *fakePendingEnrollmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.PendingEnrollment, error) {
	f.findCalled = true
	return f.findResult, f.findErr
}

func (f *fakePendingEnrollmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	f.deleteCalled = true
	f.deletedID = id
	return f.deleteErr
}

type fakeCredentialStore struct {
	encrypted []byte
	err       error
	called    bool
	plaintext []byte

	decryptResult       []byte
	decryptErr          error
	decryptCalled       bool
	decryptedCiphertext []byte
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

func (f *fakeCredentialStore) Decrypt(ctx context.Context, ciphertext []byte) ([]byte, error) {
	f.decryptCalled = true
	f.decryptedCiphertext = ciphertext
	if f.decryptErr != nil {
		return nil, f.decryptErr
	}
	return f.decryptResult, nil
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

type fakeTOTPVerifier struct {
	valid  bool
	err    error
	called bool

	secretArg string
	codeArg   string
}

func (f *fakeTOTPVerifier) Verify(secret, code string, at time.Time) (bool, error) {
	f.called = true
	f.secretArg = secret
	f.codeArg = code
	return f.valid, f.err
}

type fakeAdministratorSessionRepository struct {
	err    error
	called bool

	savedSession   *domain.AdministratorSession
	savedTokenHash string
}

func (f *fakeAdministratorSessionRepository) Save(ctx context.Context, session *domain.AdministratorSession, tokenHash string) error {
	f.called = true
	if f.err != nil {
		return f.err
	}
	f.savedSession = session
	f.savedTokenHash = tokenHash
	return nil
}

type fakeSessionTokenGenerator struct {
	token  string
	hash   string
	err    error
	called bool
}

func (f *fakeSessionTokenGenerator) Generate() (string, string, error) {
	f.called = true
	if f.err != nil {
		return "", "", f.err
	}
	return f.token, f.hash, nil
}

// fakeTransactor runs fn directly (no real DB), matching the semantics the
// real GORM-backed Transactor exposes to callers: err (if set) simulates the
// transaction itself failing to start/commit, without running fn.
type fakeTransactor struct {
	err    error
	called bool
}

func (f *fakeTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	f.called = true
	if f.err != nil {
		return f.err
	}
	return fn(ctx)
}
