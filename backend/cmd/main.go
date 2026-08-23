package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Unknowns24/akritas/backend/config"
	"github.com/Unknowns24/akritas/backend/internal/adapter/crypto/credentialcipher"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations"
	administratorrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
	administratorsessionrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator_session"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/credentialstore"
	pendingenrollmentrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/pending_enrollment"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/adapter/security/bootstrap"
	"github.com/Unknowns24/akritas/backend/internal/adapter/security/password"
	"github.com/Unknowns24/akritas/backend/internal/adapter/security/ratelimit"
	"github.com/Unknowns24/akritas/backend/internal/adapter/security/sessiontoken"
	"github.com/Unknowns24/akritas/backend/internal/adapter/security/totp"
	integrationbootstrap "github.com/Unknowns24/akritas/backend/internal/bootstrap/integrations"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	authusecase "github.com/Unknowns24/akritas/backend/internal/usecase/auth"
)

const (
	listenAddr             = ":8080"
	setupRateLimitAttempts = 5
	setupRateLimitWindow   = 15 * time.Minute
	loginRateLimitAttempts = 5
	loginRateLimitWindow   = 15 * time.Minute
)

func main() {
	configuration, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	db, err := postgres.Open(postgres.Config{
		DSN:             configuration.DatabaseURL,
		MaxOpenConns:    configuration.DatabaseMaxOpenConnections,
		MaxIdleConns:    configuration.DatabaseMaxIdleConnections,
		ConnMaxLifetime: configuration.DatabaseConnectionMaxLifetime,
	})
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	if err := migrations.Run(db); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	cipher, err := credentialcipher.NewFromKey(configuration.MasterKey)
	if err != nil {
		log.Fatalf("initialize credential cipher: %v", err)
	}
	credentials, err := credentialstore.New(db, cipher)
	if err != nil {
		log.Fatalf("initialize credential store: %v", err)
	}

	administrators := administratorrepo.NewRepository(db)
	pendingEnrollments := pendingenrollmentrepo.NewRepository(db)
	administratorSessions := administratorsessionrepo.NewRepository(db)
	transactor := postgres.NewTransactor(db)
	now := time.Now

	setupLimiter := ratelimit.New(setupRateLimitAttempts, setupRateLimitWindow)
	verifyLimiter := ratelimit.New(setupRateLimitAttempts, setupRateLimitWindow)
	loginLimiter := ratelimit.New(loginRateLimitAttempts, loginRateLimitWindow)
	passwordHasher := password.New()
	totpGenerator := totp.NewGenerator()
	totpVerifier := totp.NewVerifier()
	sessionTokens := sessiontoken.New()

	getSetupStatus := authusecase.NewGetSetupStatusUseCase(administrators)
	startSetup := authusecase.NewStartAdministratorSetupUseCase(
		administrators, pendingEnrollments, credentials, totpGenerator, passwordHasher,
		bootstrap.New(configuration.BootstrapToken), setupLimiter, transactor, now,
	)
	verifySetup := authusecase.NewVerifyAdministratorSetupUseCase(
		verifyLimiter, pendingEnrollments, credentials, totpVerifier, administrators,
		sessionTokens, administratorSessions, transactor, now,
		configuration.SessionIdleTTL, configuration.SessionAbsoluteTTL,
	)
	loginAdministrator := authusecase.NewLoginAdministratorUseCase(
		loginLimiter, administrators, passwordHasher, credentials, totpVerifier,
		sessionTokens, administratorSessions, transactor, now,
		configuration.SessionIdleTTL, configuration.SessionAbsoluteTTL,
	)
	authenticateSession := authusecase.NewAuthenticateSessionUseCase(
		administratorSessions, sessionTokens, now, configuration.SessionIdleTTL,
	)
	useCases := &portsin.UseCases{
		GetSetupStatus:           getSetupStatus,
		StartAdministratorSetup:  startSetup,
		VerifyAdministratorSetup: verifySetup,
		LoginAdministrator:       loginAdministrator,
		AuthenticateSession:      authenticateSession,
		GetCurrentSession:        authusecase.NewGetCurrentSessionUseCase(administrators),
		LogoutAdministrator:      authusecase.NewLogoutAdministratorUseCase(administratorSessions, now),
	}

	handler, err := integrationbootstrap.Build(configuration, integrationbootstrap.Dependencies{
		DB: db, Credentials: credentials,
		Admin:    middleware.RequireSession(useCases.AuthenticateSession),
		UseCases: useCases,
	})
	if err != nil {
		log.Fatalf("build application: %v", err)
	}

	log.Printf("akritas backend listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
