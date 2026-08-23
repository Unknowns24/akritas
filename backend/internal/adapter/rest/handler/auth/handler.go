package auth

import (
	"net/http"

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
	sessionCookieSameSite       http.SameSite
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
	sessionCookieSameSite ...http.SameSite,
) *Handler {
	handler := NewHandler(getSetupStatus, startAdministratorSetup, verifyAdministratorSetup, loginAdministrator, getCurrentSession, logoutAdministrator, sessionCookieSecure, sessionCookieSameSite...)
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
	sessionCookieSameSite ...http.SameSite,
) *Handler {
	sameSite := http.SameSiteLaxMode
	if len(sessionCookieSameSite) > 0 {
		sameSite = sessionCookieSameSite[0]
	}
	return &Handler{
		getSetupStatus:           getSetupStatus,
		startAdministratorSetup:  startAdministratorSetup,
		verifyAdministratorSetup: verifyAdministratorSetup,
		loginAdministrator:       loginAdministrator,
		getCurrentSession:        getCurrentSession,
		logoutAdministrator:      logoutAdministrator,
		sessionCookieSecure:      sessionCookieSecure,
		sessionCookieSameSite:    sameSite,
	}
}
