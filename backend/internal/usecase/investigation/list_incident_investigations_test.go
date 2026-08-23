package investigation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func TestListIncidentInvestigationsRejectsMissingIncident(t *testing.T) {
	t.Parallel()
	deps := newStartDeps()
	deps.incidents.exists = false
	_, err := deps.usecase().ListIncidentInvestigations(context.Background(), uuid.New(), paging.Params{})
	if !errors.Is(err, domain.ErrIncidentNotFound) {
		t.Fatalf("expected ErrIncidentNotFound, got %v", err)
	}
}

func TestListIncidentInvestigationsReturnsStorePage(t *testing.T) {
	t.Parallel()
	deps := newStartDeps()
	deps.investigations.listResult = paging.Slice[domain.Investigation]{Items: []domain.Investigation{{ID: uuid.New(), CreatedAt: time.Now()}}, Total: 1}
	page, err := deps.usecase().ListIncidentInvestigations(context.Background(), uuid.New(), paging.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Items) != 1 {
		t.Fatal("expected the store's page to be returned unchanged")
	}
}
