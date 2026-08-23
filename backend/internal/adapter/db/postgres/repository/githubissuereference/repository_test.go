//go:build integration

package githubissuereference

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func TestRepositoryPersistsIdempotentInvestigationReferenceAndFindsLatestIncidentIssue(t *testing.T) {
	db := dbtest.ConnectContainer(t)
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)
	incidentID, firstInvestigationID, secondInvestigationID := insertReferenceFixture(t, db, now)
	repository, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	first, _ := domain.NewGitHubIssueReference(incidentID, firstInvestigationID, 7, "https://github.com/acme/service/issues/7", "acme/service", now)
	second, _ := domain.NewGitHubIssueReference(incidentID, secondInvestigationID, 8, "https://github.com/acme/service/issues/8", "acme/service", now.Add(time.Minute))
	if err := repository.Create(ctx, &first); err != nil {
		t.Fatal(err)
	}
	if err := repository.Create(ctx, &second); err != nil {
		t.Fatal(err)
	}
	if err := repository.Create(ctx, &second); !errors.Is(err, domain.ErrGitHubIssueAlreadyPublished) {
		t.Fatalf("expected duplicate Investigation conflict, got %v", err)
	}
	stored, err := repository.FindByInvestigation(ctx, firstInvestigationID)
	if err != nil || stored == nil || stored.Number != 7 || stored.IncidentID != incidentID {
		t.Fatalf("stored=%+v err=%v", stored, err)
	}
	latest, err := repository.FindLatestByIncident(ctx, incidentID)
	if err != nil || latest == nil || latest.Number != 8 || latest.InvestigationID != secondInvestigationID {
		t.Fatalf("latest=%+v err=%v", latest, err)
	}
	missing, err := repository.FindByInvestigation(ctx, uuid.New())
	if err != nil || missing != nil {
		t.Fatalf("missing=%+v err=%v", missing, err)
	}
}

func insertReferenceFixture(t *testing.T, db interface {
	Exec(string, ...any) interface{ GetError() error }
}, now time.Time) (uuid.UUID, uuid.UUID, uuid.UUID) {
	t.Helper()
	panic("implemented with the repository")
}
