package integrations

import (
	"errors"
	"testing"

	"github.com/Unknowns24/akritas/backend/config"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/router"
)

func TestBuildFailsBeforeOpeningPostgreSQLWithoutAdministratorMiddleware(t *testing.T) {
	t.Parallel()

	handler, err := Build(config.Config{DatabaseURL: "postgres://must-not-be-opened.invalid/akritas"}, Dependencies{})
	if handler != nil || !errors.Is(err, router.ErrAdminMiddlewareUnavailable) {
		t.Fatalf("Build() = (%v, %v), want fail-closed middleware error", handler, err)
	}
}
