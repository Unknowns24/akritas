package auth

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type getCurrentSessionUseCase struct {
	administrators out.AdministratorRepository
}

func NewGetCurrentSessionUseCase(administrators out.AdministratorRepository) in.GetCurrentSessionUseCase {
	return &getCurrentSessionUseCase{administrators: administrators}
}

func (uc *getCurrentSessionUseCase) Execute(ctx context.Context, session domain.AdministratorSession) (in.GetCurrentSessionOutput, error) {
	administrator, err := uc.administrators.FindByID(ctx, session.AdministratorID)
	if err != nil {
		return in.GetCurrentSessionOutput{}, err
	}
	if administrator == nil {
		return in.GetCurrentSessionOutput{}, domain.ErrInactiveAdministratorSession
	}

	return in.GetCurrentSessionOutput{
		Administrator:     *administrator,
		AuthenticatedAt:   session.AuthenticatedAt,
		IdleExpiresAt:     session.IdleExpiresAt,
		AbsoluteExpiresAt: session.AbsoluteExpiresAt,
	}, nil
}
