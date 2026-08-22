package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Unknowns24/akritas/backend/config"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations"
	administratorrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
	pendingenrollmentrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/pending_enrollment"
	authhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/router"
	"github.com/Unknowns24/akritas/backend/internal/adapter/security"
	authusecase "github.com/Unknowns24/akritas/backend/internal/usecase/auth"
)

const (
	listenAddr             = ":8080"
	setupRateLimitAttempts = 5
	setupRateLimitWindow   = 15 * time.Minute
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	db, err := postgres.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	if err := migrations.Run(db); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	administrators := administratorrepo.NewRepository(db)
	pendingEnrollments := pendingenrollmentrepo.NewRepository(db)

	credentialStore, err := security.NewCredentialStore(cfg.MasterKey)
	if err != nil {
		log.Fatalf("initialize credential store: %v", err)
	}
	totpGenerator := security.NewTOTPSecretGenerator()
	passwordHasher := security.NewPasswordHasher()
	bootstrapTokens := security.NewBootstrapTokenVerifier(cfg.BootstrapToken)
	rateLimiter := security.NewRateLimiter(setupRateLimitAttempts, setupRateLimitWindow)
	clock := security.NewClock()

	getSetupStatus := authusecase.NewGetSetupStatusUseCase(administrators)
	startAdministratorSetup := authusecase.NewStartAdministratorSetupUseCase(
		administrators, pendingEnrollments, credentialStore, totpGenerator,
		passwordHasher, bootstrapTokens, rateLimiter, clock,
	)

	handler := router.New(router.Dependencies{
		Auth: authhandler.NewHandler(getSetupStatus, startAdministratorSetup),
	})

	log.Printf("akritas backend listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
