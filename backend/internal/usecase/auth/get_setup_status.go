package auth

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type getSetupStatusUseCase struct {
	administrators out.AdministratorRepository
}

func NewGetSetupStatusUseCase(administrators out.AdministratorRepository) in.GetSetupStatusUseCase {
	return &getSetupStatusUseCase{administrators: administrators}
}

func (uc *getSetupStatusUseCase) Execute(ctx context.Context) (in.SetupStatus, error) {
	exists, err := uc.administrators.ExistsActive(ctx)
	if err != nil {
		return in.SetupStatus{}, err
	}
	return in.SetupStatus{SetupRequired: !exists, RegistrationOpen: !exists}, nil
}
