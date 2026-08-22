package auth

import (
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

const sessionCookieName = "akritas_session"

type Handler struct {
	getSetupStatus           in.GetSetupStatusUseCase
	startAdministratorSetup  in.StartAdministratorSetupUseCase
	verifyAdministratorSetup in.VerifyAdministratorSetupUseCase
	sessionCookieSecure      bool
}

func NewHandler(
	getSetupStatus in.GetSetupStatusUseCase,
	startAdministratorSetup in.StartAdministratorSetupUseCase,
	verifyAdministratorSetup in.VerifyAdministratorSetupUseCase,
	sessionCookieSecure bool,
) *Handler {
	return &Handler{
		getSetupStatus:           getSetupStatus,
		startAdministratorSetup:  startAdministratorSetup,
		verifyAdministratorSetup: verifyAdministratorSetup,
		sessionCookieSecure:      sessionCookieSecure,
	}
}
