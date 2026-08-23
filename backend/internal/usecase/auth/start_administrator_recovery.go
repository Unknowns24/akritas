package auth

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

const dummyTOTPSecret = "JBSWY3DPEHPK3PXP"

type startAdministratorRecoveryUseCase struct {
	administrators      out.AdministratorRepository
	pendingEnrollments  out.PendingEnrollmentRepository
	credentialStore     out.CredentialStore
	totpSecretGenerator out.TOTPSecretGenerator
	passwordHasher      out.PasswordHasher
	bootstrapTokens     out.BootstrapTokenVerifier
	rateLimiter         out.RateLimiter
	transactor          out.Transactor
	now                 func() time.Time
	dummyPasswordHash   string
}

func NewStartAdministratorRecoveryUseCase(administrators out.AdministratorRepository, pendingEnrollments out.PendingEnrollmentRepository, credentialStore out.CredentialStore, totpSecretGenerator out.TOTPSecretGenerator, passwordHasher out.PasswordHasher, bootstrapTokens out.BootstrapTokenVerifier, rateLimiter out.RateLimiter, transactor out.Transactor, now func() time.Time, dummyPasswordHash string) in.StartAdministratorRecoveryUseCase {
	return &startAdministratorRecoveryUseCase{administrators: administrators, pendingEnrollments: pendingEnrollments, credentialStore: credentialStore, totpSecretGenerator: totpSecretGenerator, passwordHasher: passwordHasher, bootstrapTokens: bootstrapTokens, rateLimiter: rateLimiter, transactor: transactor, now: now, dummyPasswordHash: dummyPasswordHash}
}

func (uc *startAdministratorRecoveryUseCase) Execute(ctx context.Context, input in.StartAdministratorRecoveryInput) (in.StartAdministratorRecoveryOutput, error) {
	for _, key := range []string{"ip:" + input.RateLimitKey, "account:" + strings.ToLower(strings.TrimSpace(input.Email))} {
		allowed, err := uc.rateLimiter.Allow(ctx, key)
		if err != nil {
			return in.StartAdministratorRecoveryOutput{}, err
		}
		if !allowed {
			return in.StartAdministratorRecoveryOutput{}, domain.ErrAuthenticationRateLimited
		}
	}

	bootstrapOK := uc.bootstrapTokens.Verify(input.BootstrapToken)
	current, err := uc.administrators.FindByEmail(ctx, input.Email)
	if err != nil {
		return in.StartAdministratorRecoveryOutput{}, err
	}
	comparisonHash := uc.dummyPasswordHash
	if current != nil {
		comparisonHash = current.PasswordHash
	}
	samePassword, err := uc.passwordHasher.Verify(input.NewPassword, comparisonHash)
	if err != nil {
		return in.StartAdministratorRecoveryOutput{}, err
	}
	if current == nil || !bootstrapOK || samePassword {
		return in.StartAdministratorRecoveryOutput{}, domain.ErrInvalidCredentials
	}

	newPasswordHash, err := uc.passwordHasher.Hash(input.NewPassword)
	if err != nil {
		return in.StartAdministratorRecoveryOutput{}, err
	}
	secret, err := uc.totpSecretGenerator.Generate("Akritas", current.Administrator.Email)
	if err != nil {
		return in.StartAdministratorRecoveryOutput{}, err
	}
	now := uc.now().UTC()
	enrollment, err := domain.NewPendingEnrollment(uuid.New(), current.Administrator.Email, current.Administrator.DisplayName, now, now.Add(pendingEnrollmentTTL))
	if err != nil {
		return in.StartAdministratorRecoveryOutput{}, err
	}

	err = uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		latest, findErr := uc.administrators.FindByEmail(txCtx, current.Administrator.Email)
		if findErr != nil {
			return findErr
		}
		if latest == nil || latest.PasswordHash != current.PasswordHash {
			return domain.ErrInvalidCredentials
		}
		previousID, replaceErr := uc.pendingEnrollments.Replace(txCtx, enrollment, newPasswordHash)
		if replaceErr != nil {
			return replaceErr
		}
		if previousID != nil {
			if deleteErr := uc.credentialStore.DeleteOwner(txCtx, out.CredentialOwnerPendingEnrollment, *previousID); deleteErr != nil {
				return deleteErr
			}
		}
		return uc.credentialStore.Put(txCtx, out.CredentialOwnerPendingEnrollment, enrollment.ID, out.SecretValue{Kind: out.SecretKindAdministratorTOTP, Plaintext: []byte(secret.Base32Key)})
	})
	if err != nil {
		return in.StartAdministratorRecoveryOutput{}, err
	}
	return in.StartAdministratorRecoveryOutput{EnrollmentID: enrollment.ID, OtpauthURI: secret.OtpauthURI, ManualEntryKey: secret.Base32Key, ExpiresAt: enrollment.ExpiresAt}, nil
}
