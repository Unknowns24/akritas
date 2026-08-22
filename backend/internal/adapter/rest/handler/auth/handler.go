package auth

import (
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

type Handler struct {
	getSetupStatus          in.GetSetupStatusUseCase
	startAdministratorSetup in.StartAdministratorSetupUseCase
}

func NewHandler(getSetupStatus in.GetSetupStatusUseCase, startAdministratorSetup in.StartAdministratorSetupUseCase) *Handler {
	return &Handler{getSetupStatus: getSetupStatus, startAdministratorSetup: startAdministratorSetup}
}
