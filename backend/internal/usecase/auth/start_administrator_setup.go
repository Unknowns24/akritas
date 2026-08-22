package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// ErrSetupRateLimited signals the rate limit was exceeded for this request.
// It intentionally lives outside the domain error catalog: the OpenAPI
// ErrorCode pattern has no type letter for HTTP 429.
var ErrSetupRateLimited = errors.New("administrator setup rate limit exceeded")

const pendingEnrollmentTTL = 10 * time.Minute

type startAdministratorSetupUseCase struct {
	administrators      out.AdministratorRepository
	pendingEnrollments  out.PendingEnrollmentRepository
	credentialStore     out.CredentialStore
	totpSecretGenerator out.TOTPSecretGenerator
	passwordHasher      out.PasswordHasher
	bootstrapTokens     out.BootstrapTokenVerifier
	rateLimiter         out.RateLimiter
	clock               out.Clock
}

func NewStartAdministratorSetupUseCase(
	administrators out.AdministratorRepository,
	pendingEnrollments out.PendingEnrollmentRepository,
	credentialStore out.CredentialStore,
	totpSecretGenerator out.TOTPSecretGenerator,
	passwordHasher out.PasswordHasher,
	bootstrapTokens out.BootstrapTokenVerifier,
	rateLimiter out.RateLimiter,
	clock out.Clock,
) in.StartAdministratorSetupUseCase {
	return &startAdministratorSetupUseCase{
		administrators:      administrators,
		pendingEnrollments:  pendingEnrollments,
		credentialStore:     credentialStore,
		totpSecretGenerator: totpSecretGenerator,
		passwordHasher:      passwordHasher,
		bootstrapTokens:     bootstrapTokens,
		rateLimiter:         rateLimiter,
		clock:               clock,
	}
}

func (uc *startAdministratorSetupUseCase) Execute(ctx context.Context, input in.StartAdministratorSetupInput) (in.StartAdministratorSetupOutput, error) {
	allowed, err := uc.rateLimiter.Allow(ctx, input.RateLimitKey)
	if err != nil {
		return in.StartAdministratorSetupOutput{}, err
	}
	if !allowed {
		return in.StartAdministratorSetupOutput{}, ErrSetupRateLimited
	}

	if !uc.bootstrapTokens.Verify(input.BootstrapToken) {
		return in.StartAdministratorSetupOutput{}, domain.ErrInvalidBootstrapToken
	}

	exists, err := uc.administrators.ExistsActive(ctx)
	if err != nil {
		return in.StartAdministratorSetupOutput{}, err
	}
	if exists {
		return in.StartAdministratorSetupOutput{}, domain.ErrAdministratorAlreadyExists
	}

	passwordHash, err := uc.passwordHasher.Hash(input.Password)
	if err != nil {
		return in.StartAdministratorSetupOutput{}, err
	}

	secret, err := uc.totpSecretGenerator.Generate("Akritas", input.Email)
	if err != nil {
		return in.StartAdministratorSetupOutput{}, err
	}

	encryptedSecret, err := uc.credentialStore.Encrypt(ctx, []byte(secret.Base32Key))
	if err != nil {
		return in.StartAdministratorSetupOutput{}, err
	}

	now := uc.clock.Now()
	enrollment, err := domain.NewPendingEnrollment(uuid.New(), input.Email, input.DisplayName, passwordHash, encryptedSecret, now, now.Add(pendingEnrollmentTTL))
	if err != nil {
		return in.StartAdministratorSetupOutput{}, err
	}

	if err := uc.pendingEnrollments.Save(ctx, enrollment); err != nil {
		return in.StartAdministratorSetupOutput{}, err
	}

	return in.StartAdministratorSetupOutput{
		EnrollmentID:   enrollment.ID,
		OtpauthURI:     secret.OtpauthURI,
		ManualEntryKey: secret.Base32Key,
		ExpiresAt:      enrollment.ExpiresAt,
	}, nil
}
