package investigation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func TestRepositoryPersistsCreatesFindsAndUpdates(t *testing.T) {
	db := dbtest.Connect(t)
	repository, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Microsecond)
	incidentID := uuid.New()

	value, err := domain.NewInvestigation(uuid.New(), incidentID, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.Create(ctx, value); err != nil {
		t.Fatalf("create: %v", err)
	}

	found, err := repository.FindByID(ctx, value.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if found.Status != domain.InvestigationStatusPending || len(found.Hypotheses) != 0 {
		t.Fatalf("expected a pending investigation with empty slices, got %+v", found)
	}

	if err := found.Start(now.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := repository.Update(ctx, found); err != nil {
		t.Fatalf("update to running: %v", err)
	}

	rootCause := domain.RootCauseIdentified
	resolution := domain.ResolutionFixable
	if err := found.Complete(now.Add(2*time.Second), "summary", "root cause", rootCause, resolution, 0.75,
		[]string{"h1", "h2"}, []string{"main.go"}, []string{"abc123"}, []string{"patch"}); err != nil {
		t.Fatal(err)
	}
	if err := repository.Update(ctx, found); err != nil {
		t.Fatalf("update to completed: %v", err)
	}

	reloaded, err := repository.FindByID(ctx, value.ID)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Status != domain.InvestigationStatusCompleted || reloaded.Summary != "summary" ||
		len(reloaded.Hypotheses) != 2 || reloaded.Hypotheses[1] != "h2" || reloaded.IncidentID != incidentID {
		t.Fatalf("expected persisted completed investigation, got %+v", reloaded)
	}
}

func TestRepositoryFindByIDNotFound(t *testing.T) {
	db := dbtest.Connect(t)
	repository, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repository.FindByID(context.Background(), uuid.New()); !errors.Is(err, domain.ErrInvestigationNotFound) {
		t.Fatalf("expected ErrInvestigationNotFound, got %v", err)
	}
}

func TestRepositoryExistsActiveForIncident(t *testing.T) {
	db := dbtest.Connect(t)
	repository, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Microsecond)
	incidentID := uuid.New()

	active, err := repository.ExistsActiveForIncident(ctx, incidentID)
	if err != nil || active {
		t.Fatalf("expected no active investigation yet, got active=%v err=%v", active, err)
	}

	pending, err := domain.NewInvestigation(uuid.New(), incidentID, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.Create(ctx, pending); err != nil {
		t.Fatal(err)
	}
	active, err = repository.ExistsActiveForIncident(ctx, incidentID)
	if err != nil || !active {
		t.Fatalf("expected a pending investigation to count as active, got active=%v err=%v", active, err)
	}

	if err := pending.Start(now.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	rootCause := domain.RootCauseIdentified
	resolution := domain.ResolutionFixable
	if err := pending.Complete(now.Add(2*time.Second), "summary", "cause", rootCause, resolution, 0.5, nil, nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if err := repository.Update(ctx, pending); err != nil {
		t.Fatal(err)
	}
	active, err = repository.ExistsActiveForIncident(ctx, incidentID)
	if err != nil || active {
		t.Fatalf("expected no active investigation once completed, got active=%v err=%v", active, err)
	}
}

func TestRepositoryListByIncidentScopesAndPaginates(t *testing.T) {
	db := dbtest.Connect(t)
	repository, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Microsecond)
	incidentID := uuid.New()
	otherIncidentID := uuid.New()

	for i := 0; i < 3; i++ {
		value, err := domain.NewInvestigation(uuid.New(), incidentID, now.Add(time.Duration(i)*time.Second))
		if err != nil {
			t.Fatal(err)
		}
		if err := repository.Create(ctx, value); err != nil {
			t.Fatal(err)
		}
	}
	other, err := domain.NewInvestigation(uuid.New(), otherIncidentID, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.Create(ctx, other); err != nil {
		t.Fatal(err)
	}

	page, err := repository.ListByIncident(ctx, incidentID, paging.Params{
		Limit: 10, Sort: []ukerpagination.SortExpression{{Field: "created_at", Direction: ukerpagination.DirectionAsc}, {Field: "id", Direction: ukerpagination.DirectionAsc}},
	})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if page.Total != 3 || len(page.Items) != 3 {
		t.Fatalf("expected 3 investigations scoped to the incident, got total=%d items=%d", page.Total, len(page.Items))
	}
	for _, item := range page.Items {
		if item.IncidentID != incidentID {
			t.Fatalf("list leaked an investigation from another incident: %+v", item)
		}
	}
}
