package investigation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func TestGetInvestigationHappyPath(t *testing.T) {
	t.Parallel()
	deps := newStartDeps()
	deps.investigations.findByIDResult = &domain.Investigation{ID: uuid.New(), CreatedAt: time.Now()}
	value, err := deps.usecase().GetInvestigation(context.Background(), deps.investigations.findByIDResult.ID)
	if err != nil {
		t.Fatal(err)
	}
	if value.ID != deps.investigations.findByIDResult.ID {
		t.Fatal("expected the store's investigation to be returned unchanged")
	}
}

func TestGetInvestigationNotFound(t *testing.T) {
	t.Parallel()
	deps := newStartDeps()
	deps.investigations.findByIDErr = domain.ErrInvestigationNotFound
	if _, err := deps.usecase().GetInvestigation(context.Background(), uuid.New()); !errors.Is(err, domain.ErrInvestigationNotFound) {
		t.Fatalf("expected ErrInvestigationNotFound, got %v", err)
	}
}
