package integrations

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Unknowns24/akritas/backend/config"
	resthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/router"
)

func TestBuildFailsBeforeOpeningPostgreSQLWithoutAdministratorMiddleware(t *testing.T) {
	t.Parallel()

	handler, err := Build(config.Config{DatabaseURL: "postgres://must-not-be-opened.invalid/akritas"}, Dependencies{})
	if handler != nil || !errors.Is(err, router.ErrAdminMiddlewareUnavailable) {
		t.Fatalf("Build() = (%v, %v), want fail-closed middleware error", handler, err)
	}
}

func TestBuildFailsBeforeOpeningPostgreSQLWithoutUseCases(t *testing.T) {
	t.Parallel()

	passThrough := func(next http.Handler) http.Handler { return next }
	handler, err := Build(
		config.Config{DatabaseURL: "postgres://must-not-be-opened.invalid/akritas"},
		Dependencies{Admin: passThrough},
	)
	if handler != nil || !errors.Is(err, resthandler.ErrInvalidHandlersConfiguration) {
		t.Fatalf("Build() = (%v, %v), want fail-closed handlers error", handler, err)
	}
}
