package issuecontent

import (
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func TestBuilderProducesDeterministicAuditableAndSafeContent(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	accountID, incidentID, investigationID := uuid.New(), uuid.New(), uuid.New()
	repository, _ := domain.NewGitHubRepository(accountID, "42", "Unknowns24", "akritas", "main", true, "https://github.com/Unknowns24/akritas")
	application, _ := domain.NewDokployApplication(uuid.New(), "app-1", "instance-1", "Akritas API", "production", domain.DokployApplicationRunning)
	project, err := domain.NewProject(uuid.New(), "Akritas", "", repository, application, domain.DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatal(err)
	}
	incident, err := domain.NewIncident(incidentID, "AKR-45", project.ID, "sha256:0123456789abcdef", domain.SeverityError,
		"Database failure github_pat_ABCDEFGHIJKLMNOPQRSTUVWXYZ123456", now)
	if err != nil {
		t.Fatal(err)
	}
	incident.OccurrenceCount = 3
	incident.LastSeenAt = now.Add(2 * time.Minute)
	investigation, _ := domain.NewInvestigation(investigationID, incidentID, now)
	_ = investigation.Start(now.Add(time.Second))
	evidenceID := uuid.New()
	investigation.EvidenceCount = 1
	_ = investigation.Complete(now.Add(3*time.Minute), "Connection failed", "DATABASE_PASSWORD=supersecret",
		domain.RootCauseIdentified, domain.ResolutionFixable, 0.84,
		[]string{"The database is unavailable"}, []string{"internal/db/connect.go"}, []string{"deadbeef"}, []string{"Retry safely"}, []uuid.UUID{evidenceID})
	evidence, err := domain.NewEvidence(evidenceID, investigationID, domain.EvidenceLogExcerpt, "Observed refusal",
		"Authorization: Bearer abc.def.ghi\nconnection refused", now)
	if err != nil {
		t.Fatal(err)
	}

	builder := New()
	first, err := builder.Build(Input{Project: *project, Incident: *incident, Investigation: *investigation, Evidence: []domain.Evidence{*evidence}})
	if err != nil {
		t.Fatal(err)
	}
	second, err := builder.Build(Input{Project: *project, Incident: *incident, Investigation: *investigation, Evidence: []domain.Evidence{*evidence}})
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatal("same input must produce identical content")
	}
	for _, want := range []string{
		"[AKR-45] Database failure", "## Project Context", "Akritas API", "Unknowns24/akritas",
		"## Incident — Observed", incidentID.String(), "Occurrences", now.Format(time.RFC3339),
		"## Observed Evidence", "Observed refusal", "connection refused",
		"## Investigation — QVAC Analysis", "model-generated conclusions", "identified", "0.8400", "fixable",
		"internal/db/connect.go", "deadbeef", "Retry safely", investigationID.String(),
	} {
		if !strings.Contains(first.Title+"\n"+first.Body, want) {
			t.Fatalf("content missing %q:\n%s\n%s", want, first.Title, first.Body)
		}
	}
	for _, secret := range []string{"github_pat_", "supersecret", "abc.def.ghi"} {
		if strings.Contains(first.Title+first.Body, secret) {
			t.Fatalf("secret %q leaked", secret)
		}
	}
	if len(first.Title) > maximumTitleBytes || len(first.Body) > maximumBodyBytes {
		t.Fatalf("content limits exceeded: title=%d body=%d", len(first.Title), len(first.Body))
	}
}

func TestBuilderRejectsNonCompletedInvestigation(t *testing.T) {
	t.Parallel()
	investigation, _ := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	if _, err := New().Build(Input{Investigation: *investigation}); err == nil {
		t.Fatal("pending Investigation must not produce Issue content")
	}
}
