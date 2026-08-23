package evidence

import (
	"context"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func TestRepositoryCreatesAndListsScopedByInvestigation(t *testing.T) {
	db := dbtest.Connect(t)
	repository, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Microsecond)
	incidentID := uuid.New()

	investigation, err := domain.NewInvestigation(uuid.New(), incidentID, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Table("investigations").Create(investigation).Error; err != nil {
		t.Fatalf("seed investigation: %v", err)
	}
	otherInvestigation, err := domain.NewInvestigation(uuid.New(), incidentID, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Table("investigations").Create(otherInvestigation).Error; err != nil {
		t.Fatalf("seed other investigation: %v", err)
	}

	lineStart, lineEnd := 10, 20
	full, err := domain.NewEvidence(uuid.New(), investigation.ID, domain.EvidenceCodeLocation, "summary", "content", now)
	if err != nil {
		t.Fatal(err)
	}
	full.FilePath = "main.go"
	full.LineStart = &lineStart
	full.LineEnd = &lineEnd
	full.CommitSHA = "abc123"
	full.Patch = "diff --git a/main.go b/main.go"
	full.OccurredAt = &now
	if err := full.Validate(); err != nil {
		t.Fatalf("fixture is invalid: %v", err)
	}
	if err := repository.Create(ctx, full); err != nil {
		t.Fatalf("create: %v", err)
	}

	deployment, err := domain.NewEvidence(uuid.New(), investigation.ID, domain.EvidenceDeploymentMetadata, "deployment", "{}", now.Add(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.Create(ctx, deployment); err != nil {
		t.Fatalf("create deployment: %v", err)
	}

	other, err := domain.NewEvidence(uuid.New(), otherInvestigation.ID, domain.EvidenceDeploymentMetadata, "other", "{}", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.Create(ctx, other); err != nil {
		t.Fatalf("create other: %v", err)
	}

	page, err := repository.ListByInvestigation(ctx, investigation.ID, paging.Params{
		Limit: 10, Sort: []ukerpagination.SortExpression{{Field: "created_at", Direction: ukerpagination.DirectionAsc}, {Field: "id", Direction: ukerpagination.DirectionAsc}},
	})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if page.Total != 2 || len(page.Items) != 2 {
		t.Fatalf("expected 2 evidence scoped to the investigation, got total=%d items=%d", page.Total, len(page.Items))
	}
	for _, item := range page.Items {
		if item.InvestigationID != investigation.ID {
			t.Fatalf("list leaked evidence from another investigation: %+v", item)
		}
	}

	reloaded := page.Items[0]
	if reloaded.FilePath != "main.go" || reloaded.LineStart == nil || *reloaded.LineStart != lineStart ||
		reloaded.LineEnd == nil || *reloaded.LineEnd != lineEnd || reloaded.CommitSHA != "abc123" ||
		reloaded.Patch != "diff --git a/main.go b/main.go" || reloaded.OccurredAt == nil || !reloaded.Redacted {
		t.Fatalf("expected optional fields to persist correctly, got %+v", reloaded)
	}
}
