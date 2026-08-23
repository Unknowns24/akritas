package auth

import (
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

type Handler struct {
	getSetupStatus              in.GetSetupStatusUseCase
	startAdministratorSetup     in.StartAdministratorSetupUseCase
	verifyAdministratorSetup    in.VerifyAdministratorSetupUseCase
	loginAdministrator          in.LoginAdministratorUseCase
	startAdministratorRecovery  in.StartAdministratorRecoveryUseCase
	verifyAdministratorRecovery in.VerifyAdministratorRecoveryUseCase
	getCurrentSession           in.GetCurrentSessionUseCase
	logoutAdministrator         in.LogoutAdministratorUseCase
	sessionCookieSecure         bool
}

func NewHandlerWithRecovery(
	getSetupStatus in.GetSetupStatusUseCase,
	startAdministratorSetup in.StartAdministratorSetupUseCase,
	verifyAdministratorSetup in.VerifyAdministratorSetupUseCase,
	loginAdministrator in.LoginAdministratorUseCase,
	startAdministratorRecovery in.StartAdministratorRecoveryUseCase,
	verifyAdministratorRecovery in.VerifyAdministratorRecoveryUseCase,
	getCurrentSession in.GetCurrentSessionUseCase,
	logoutAdministrator in.LogoutAdministratorUseCase,
	sessionCookieSecure bool,
) *Handler {
	handler := NewHandler(getSetupStatus, startAdministratorSetup, verifyAdministratorSetup, loginAdministrator, getCurrentSession, logoutAdministrator, sessionCookieSecure)
	handler.startAdministratorRecovery = startAdministratorRecovery
	handler.verifyAdministratorRecovery = verifyAdministratorRecovery
	return handler
}

func NewHandler(
	getSetupStatus in.GetSetupStatusUseCase,
	startAdministratorSetup in.StartAdministratorSetupUseCase,
	verifyAdministratorSetup in.VerifyAdministratorSetupUseCase,
	loginAdministrator in.LoginAdministratorUseCase,
	getCurrentSession in.GetCurrentSessionUseCase,
	logoutAdministrator in.LogoutAdministratorUseCase,
	sessionCookieSecure bool,
) *Handler {
	return &Handler{
		getSetupStatus:           getSetupStatus,
		startAdministratorSetup:  startAdministratorSetup,
		verifyAdministratorSetup: verifyAdministratorSetup,
		loginAdministrator:       loginAdministrator,
		getCurrentSession:        getCurrentSession,
		logoutAdministrator:      logoutAdministrator,
		sessionCookieSecure:      sessionCookieSecure,
	}
}
