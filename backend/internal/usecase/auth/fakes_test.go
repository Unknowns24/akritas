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

	createErr            error
	createCalled         bool
	createdAdministrator *domain.Administrator
	createdPasswordHash  string

	findByIDResult *domain.Administrator
	findByIDErr    error
	findByIDCalled bool

	findByEmailResult *out.AdministratorAuthentication
	findByEmailErr    error
	findByEmailCalled bool

	consumePeriodErr    error
	consumePeriodReject bool
	updatePeriodCalled  bool
	updatedPeriodID     uuid.UUID
	updatedPeriod       int64
	updatedPasswordHash string

	rotateResult        bool
	rotateErr           error
	rotateCalled        bool
	rotatedPasswordHash string
}

func (f *fakeAdministratorRepository) ExistsActive(ctx context.Context) (bool, error) {
	f.called = true
	if f.order != nil {
		*f.order = append(*f.order, "administrator_repository")
	}
	return f.exists, f.err
}

func (f *fakeAdministratorRepository) Create(ctx context.Context, administrator *domain.Administrator, passwordHash string) error {
	f.createCalled = true
	if f.order != nil {
		*f.order = append(*f.order, "administrator_repository.create")
	}
	if f.createErr != nil {
		return f.createErr
	}
	f.createdAdministrator = administrator
	f.createdPasswordHash = passwordHash
	return nil
}

func (f *fakeAdministratorRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Administrator, error) {
	f.findByIDCalled = true
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeAdministratorRepository) FindByEmail(ctx context.Context, email string) (*out.AdministratorAuthentication, error) {
	f.findByEmailCalled = true
	return f.findByEmailResult, f.findByEmailErr
}

func (f *fakeAdministratorRepository) ConsumeTOTPPeriod(ctx context.Context, id uuid.UUID, expectedPasswordHash string, period int64) (bool, error) {
	f.updatePeriodCalled = true
	if f.consumePeriodErr != nil {
		return false, f.consumePeriodErr
	}
	f.updatedPeriodID = id
	f.updatedPeriod = period
	f.updatedPasswordHash = expectedPasswordHash
	return !f.consumePeriodReject, nil
}

func (f *fakeAdministratorRepository) RotateCredentials(ctx context.Context, id uuid.UUID, expectedPasswordHash, newPasswordHash string, acceptedTOTPPeriod int64, updatedAt time.Time) (bool, error) {
	f.rotateCalled = true
	f.rotatedPasswordHash = newPasswordHash
	if f.rotateErr != nil {
		return false, f.rotateErr
	}
	if !f.rotateResult {
		return false, nil
	}
	return true, nil
}

type fakePendingEnrollmentRepository struct {
	saved  *domain.PendingEnrollment
	err    error
	called bool

	passwordHash string
	previousID   *uuid.UUID
	findResult   *out.PendingEnrollmentAuthentication
	findErr      error
	findCalled   bool

	deleteErr    error
	deleteCalled bool
	deletedID    uuid.UUID
}

func (f *fakePendingEnrollmentRepository) Replace(ctx context.Context, enrollment *domain.PendingEnrollment, passwordHash string) (*uuid.UUID, error) {
	f.called = true
	if f.err != nil {
		return nil, f.err
	}
	f.saved = enrollment
	f.passwordHash = passwordHash
	return f.previousID, nil
}

func (f *fakePendingEnrollmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*out.PendingEnrollmentAuthentication, error) {
	f.findCalled = true
	return f.findResult, f.findErr
}

func (f *fakePendingEnrollmentRepository) Consume(ctx context.Context, id uuid.UUID) (*out.PendingEnrollmentAuthentication, error) {
	f.deleteCalled = true
	f.deletedID = id
	if f.deleteErr != nil {
		return nil, f.deleteErr
	}
	return f.findResult, nil
}

func (f *fakePendingEnrollmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	f.deleteCalled = true
	f.deletedID = id
	return f.deleteErr
}

type fakeCredentialStore struct {
	err            error
	called         bool
	plaintext      []byte
	getResult      []byte
	getErr         error
	getCalled      bool
	ownerType      string
	ownerID        uuid.UUID
	kind           out.SecretKind
	deletedOwnerID uuid.UUID
	moved          bool
}

func (f *fakeCredentialStore) Put(ctx context.Context, ownerType string, ownerID uuid.UUID, secret out.SecretValue) error {
	f.called = true
	f.ownerType, f.ownerID, f.kind = ownerType, ownerID, secret.Kind
	f.plaintext = append([]byte(nil), secret.Plaintext...)
	return f.err
}

func (f *fakeCredentialStore) Get(ctx context.Context, ownerType string, ownerID uuid.UUID, kind out.SecretKind) ([]byte, error) {
	f.getCalled = true
	f.ownerType, f.ownerID, f.kind = ownerType, ownerID, kind
	return append([]byte(nil), f.getResult...), f.getErr
}

func (f *fakeCredentialStore) DeleteOwner(ctx context.Context, ownerType string, ownerID uuid.UUID) error {
	f.deletedOwnerID = ownerID
	return f.err
}

func (f *fakeCredentialStore) MoveOwner(ctx context.Context, fromType string, fromID uuid.UUID, toType string, toID uuid.UUID) error {
	f.moved = true
	return f.err
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

	verifyResult bool
	verifyErr    error
	verifyCalled bool
	verifyArgs   [2]string
}

func (f *fakePasswordHasher) Hash(password string) (string, error) {
	f.called = true
	if f.err != nil {
		return "", f.err
	}
	return f.hash, nil
}

func (f *fakePasswordHasher) Verify(password, hash string) (bool, error) {
	f.verifyCalled = true
	f.verifyArgs = [2]string{password, hash}
	return f.verifyResult, f.verifyErr
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
	period int64
	err    error
	called bool

	secretArg string
	codeArg   string
}

func (f *fakeTOTPVerifier) Verify(secret, code string, at time.Time) (bool, int64, error) {
	f.called = true
	f.secretArg = secret
	f.codeArg = code
	return f.valid, f.period, f.err
}

type fakeAdministratorSessionRepository struct {
	err    error
	called bool

	savedSession   *domain.AdministratorSession
	savedTokenHash string

	findByTokenHashResult *domain.AdministratorSession
	findByTokenHashErr    error
	findByTokenHashCalled bool

	updateIdleErr        error
	updateIdleCalled     bool
	updatedIdleSessionID uuid.UUID
	updatedIdleExpiresAt time.Time

	revokeErr        error
	revokeCalled     bool
	revokedSessionID uuid.UUID
	revokedAt        time.Time
	revokeAllErr     error
	revokeAllCalled  bool
	revokedAdminID   uuid.UUID
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

func (f *fakeAdministratorSessionRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.AdministratorSession, error) {
	f.findByTokenHashCalled = true
	return f.findByTokenHashResult, f.findByTokenHashErr
}

func (f *fakeAdministratorSessionRepository) UpdateIdleExpiry(ctx context.Context, id uuid.UUID, idleExpiresAt time.Time) error {
	f.updateIdleCalled = true
	if f.updateIdleErr != nil {
		return f.updateIdleErr
	}
	f.updatedIdleSessionID = id
	f.updatedIdleExpiresAt = idleExpiresAt
	return nil
}

func (f *fakeAdministratorSessionRepository) RefreshActive(ctx context.Context, tokenHash string, now, requestedIdleExpiry time.Time) (*domain.AdministratorSession, error) {
	f.findByTokenHashCalled = true
	if f.findByTokenHashErr != nil {
		return nil, f.findByTokenHashErr
	}
	if f.updateIdleErr != nil {
		return nil, f.updateIdleErr
	}
	if f.findByTokenHashResult == nil || !f.findByTokenHashResult.IsActive(now) {
		return nil, nil
	}
	session := *f.findByTokenHashResult
	if requestedIdleExpiry.After(session.AbsoluteExpiresAt) {
		requestedIdleExpiry = session.AbsoluteExpiresAt
	}
	session.IdleExpiresAt = requestedIdleExpiry
	f.updateIdleCalled = true
	f.updatedIdleSessionID = session.ID
	f.updatedIdleExpiresAt = requestedIdleExpiry
	return &session, nil
}

func (f *fakeAdministratorSessionRepository) Revoke(ctx context.Context, id uuid.UUID, revokedAt time.Time) error {
	f.revokeCalled = true
	if f.revokeErr != nil {
		return f.revokeErr
	}
	f.revokedSessionID = id
	f.revokedAt = revokedAt
	return nil
}

func (f *fakeAdministratorSessionRepository) RevokeAll(ctx context.Context, administratorID uuid.UUID, revokedAt time.Time) error {
	f.revokeAllCalled = true
	f.revokedAdminID = administratorID
	return f.revokeAllErr
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

func (f *fakeSessionTokenGenerator) Hash(token string) string {
	return "hash-of:" + token
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
