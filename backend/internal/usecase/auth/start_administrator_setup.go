package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

const pendingEnrollmentTTL = 10 * time.Minute

type startAdministratorSetupUseCase struct {
	administrators      out.AdministratorRepository
	pendingEnrollments  out.PendingEnrollmentRepository
	credentialStore     out.CredentialStore
	totpSecretGenerator out.TOTPSecretGenerator
	passwordHasher      out.PasswordHasher
	bootstrapTokens     out.BootstrapTokenVerifier
	rateLimiter         out.RateLimiter
	transactor          out.Transactor
	now                 func() time.Time
}

func NewStartAdministratorSetupUseCase(
	administrators out.AdministratorRepository,
	pendingEnrollments out.PendingEnrollmentRepository,
	credentialStore out.CredentialStore,
	totpSecretGenerator out.TOTPSecretGenerator,
	passwordHasher out.PasswordHasher,
	bootstrapTokens out.BootstrapTokenVerifier,
	rateLimiter out.RateLimiter,
	transactor out.Transactor,
	now func() time.Time,
) in.StartAdministratorSetupUseCase {
	return &startAdministratorSetupUseCase{
		administrators:      administrators,
		pendingEnrollments:  pendingEnrollments,
		credentialStore:     credentialStore,
		totpSecretGenerator: totpSecretGenerator,
		passwordHasher:      passwordHasher,
		bootstrapTokens:     bootstrapTokens,
		rateLimiter:         rateLimiter,
		transactor:          transactor,
		now:                 now,
	}
}

func (uc *startAdministratorSetupUseCase) Execute(ctx context.Context, input in.StartAdministratorSetupInput) (in.StartAdministratorSetupOutput, error) {
	allowed, err := uc.rateLimiter.Allow(ctx, input.RateLimitKey)
	if err != nil {
		return in.StartAdministratorSetupOutput{}, err
	}
	if !allowed {
		return in.StartAdministratorSetupOutput{}, domain.ErrAuthenticationRateLimited
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

	now := uc.now().UTC()
	enrollment, err := domain.NewPendingEnrollment(uuid.New(), input.Email, input.DisplayName, now, now.Add(pendingEnrollmentTTL))
	if err != nil {
		return in.StartAdministratorSetupOutput{}, err
	}

	err = uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		previousID, replaceErr := uc.pendingEnrollments.Replace(txCtx, enrollment, passwordHash)
		if replaceErr != nil {
			return replaceErr
		}
		if previousID != nil {
			if deleteErr := uc.credentialStore.DeleteOwner(txCtx, out.CredentialOwnerPendingEnrollment, *previousID); deleteErr != nil {
				return deleteErr
			}
		}
		return uc.credentialStore.Put(txCtx, out.CredentialOwnerPendingEnrollment, enrollment.ID, out.SecretValue{
			Kind: out.SecretKindAdministratorTOTP, Plaintext: []byte(secret.Base32Key),
		})
	})
	if err != nil {
		return in.StartAdministratorSetupOutput{}, err
	}

	return in.StartAdministratorSetupOutput{
		EnrollmentID:   enrollment.ID,
		OtpauthURI:     secret.OtpauthURI,
		ManualEntryKey: secret.Base32Key,
		ExpiresAt:      enrollment.ExpiresAt,
	}, nil
}
