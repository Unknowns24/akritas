package issuecontent

import (
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func TestBuilderProducesDeterministicAuditableAndSafeContent(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	accountID, incidentID, investigationID := uuid.New(), uuid.New(), uuid.New()
	repository, _ := domain.NewGitHubRepository(accountID, "42", "Unknowns24", "akritas", "main", true, "https://github.com/Unknowns24/akritas")
	application, _ := domain.NewDokployApplication(uuid.New(), "app-1", "instance-1", "Akritas API", "production", domain.DokployApplicationRunning)
	source, _ := domain.SourceFromApplication(application)
	project, err := domain.NewProject(uuid.New(), "Akritas", "", repository, source, domain.DefaultMonitoringConfiguration(), now)
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
			t.Fatal("content leaked sensitive fixture data")
		}
	}
	if len(first.Title) > maximumTitleBytes || len(first.Body) > maximumBodyBytes {
		t.Fatalf("content limits exceeded: title=%d body=%d", len(first.Title), len(first.Body))
	}
}

func TestBuilderRedactsSecretsAcrossPublishedFieldsAndKeepsAuditContent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		mutate func(*Input) []string
	}{
		{
			name: "incident title",
			mutate: func(input *Input) []string {
				input.Incident.Title = `Database failure Authorization: Basic incident-title-secret`
				return []string{"incident-title-secret"}
			},
		},
		{
			name: "project application and repository context",
			mutate: func(input *Input) []string {
				input.Project.Name = `Akritas PASSWORD="project secret"`
				input.Project.DokploySource.DisplayName = `API TOKEN='application secret'`
				input.Project.DokploySource.Environment = `production SECRET="environment secret"`
				input.Project.GitHubRepository.Name = `backend-token=repository-secret`
				input.Project.GitHubRepository.FullName = input.Project.GitHubRepository.Owner + "/" + input.Project.GitHubRepository.Name
				input.Project.GitHubRepository.DefaultBranch = `main API_KEY="branch secret"`
				return []string{"project secret", "application secret", "environment secret", "repository-secret", "branch secret"}
			},
		},
		{
			name: "evidence summary content patch file and commit",
			mutate: func(input *Input) []string {
				input.Evidence[0].Summary = `Observed refusal {"password":"evidence-summary-secret"}`
				input.Evidence[0].Content = `Authorization: Bearer evidence-content-secret`
				input.Evidence[0].Patch = `+ TOKEN='evidence patch secret'`
				input.Evidence[0].FilePath = `internal/password=file-secret/connect.go`
				input.Evidence[0].CommitSHA = `ghs_ABCDEFGHIJKLMNOPQRSTUVWXYZ123456`
				return []string{"evidence-summary-secret", "evidence-content-secret", "evidence patch secret", "file-secret", "ghs_ABCDEFGHIJKLMNOPQRSTUVWXYZ123456"}
			},
		},
		{
			name: "root cause and summary",
			mutate: func(input *Input) []string {
				input.Investigation.RootCause = `DATABASE_PASSWORD="root cause secret"`
				input.Investigation.Summary = `Incident reproduced with token=summary-secret`
				return []string{"root cause secret", "summary-secret"}
			},
		},
		{
			name: "hypotheses files commits and actions",
			mutate: func(input *Input) []string {
				input.Investigation.Hypotheses = append(input.Investigation.Hypotheses, `Hypothesis uses API_KEY="hypothesis secret"`)
				input.Investigation.RelevantFiles = append(input.Investigation.RelevantFiles, `internal/secret=relevant-file-secret/config.go`)
				input.Investigation.RelevantCommits = append(input.Investigation.RelevantCommits, `github_pat_11ABCDEFGHijklmnopQRSTUVwxYZ1234567890`)
				input.Investigation.RecommendedActions = append(input.Investigation.RecommendedActions, `Rotate TOKEN='action secret'`)
				return []string{"hypothesis secret", "relevant-file-secret", "github_pat_11ABCDEFGHijklmnopQRSTUVwxYZ1234567890", "action secret"}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input, incidentID, investigationID := issueContentFixture(t)
			leaks := tt.mutate(&input)
			content, err := New().Build(input)
			if err != nil {
				t.Fatal(err)
			}
			combined := content.Title + "\n" + content.Body
			for _, leak := range leaks {
				if strings.Contains(combined, leak) {
					t.Fatalf("case %q leaked sensitive fixture data", tt.name)
				}
			}
			assertAuditableIssueContent(t, content, incidentID, investigationID)
		})
	}
}

func TestBuilderRejectsNonCompletedInvestigation(t *testing.T) {
	t.Parallel()
	investigation, _ := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	if _, err := New().Build(Input{Investigation: *investigation}); err == nil {
		t.Fatal("pending Investigation must not produce Issue content")
	}
}

func issueContentFixture(t *testing.T) (Input, uuid.UUID, uuid.UUID) {
	t.Helper()
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	accountID, incidentID, investigationID := uuid.New(), uuid.New(), uuid.New()
	repository, err := domain.NewGitHubRepository(accountID, "42", "Unknowns24", "akritas", "main", true, "https://github.com/Unknowns24/akritas")
	if err != nil {
		t.Fatal(err)
	}
	application, err := domain.NewDokployApplication(uuid.New(), "app-1", "instance-1", "Akritas API", "production", domain.DokployApplicationRunning)
	if err != nil {
		t.Fatal(err)
	}
	source, err := domain.SourceFromApplication(application)
	if err != nil {
		t.Fatal(err)
	}
	project, err := domain.NewProject(uuid.New(), "Akritas", "", repository, source, domain.DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatal(err)
	}
	incident, err := domain.NewIncident(incidentID, "AKR-47", project.ID, "sha256:0123456789abcdef", domain.SeverityError,
		"Database failure", now)
	if err != nil {
		t.Fatal(err)
	}
	incident.OccurrenceCount = 3
	incident.LastSeenAt = now.Add(2 * time.Minute)
	investigation, err := domain.NewInvestigation(investigationID, incidentID, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := investigation.Start(now.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	evidenceID := uuid.New()
	investigation.EvidenceCount = 1
	if err := investigation.Complete(now.Add(3*time.Minute), "Connection failed", "Database outage",
		domain.RootCauseIdentified, domain.ResolutionFixable, 0.84,
		[]string{"The database is unavailable"}, []string{"internal/db/connect.go"}, []string{"deadbeef"}, []string{"Retry safely"}, []uuid.UUID{evidenceID}); err != nil {
		t.Fatal(err)
	}
	evidence, err := domain.NewEvidence(evidenceID, investigationID, domain.EvidenceLogExcerpt, "Observed refusal",
		"connection refused", now)
	if err != nil {
		t.Fatal(err)
	}
	evidence.Patch = "safe patch"
	return Input{Project: *project, Incident: *incident, Investigation: *investigation, Evidence: []domain.Evidence{*evidence}}, incidentID, investigationID
}

func assertAuditableIssueContent(t *testing.T, content portsout.IssueContent, incidentID, investigationID uuid.UUID) {
	t.Helper()
	if len(content.Title) > maximumTitleBytes || len(content.Body) > maximumBodyBytes {
		t.Fatalf("content limits exceeded: title=%d body=%d", len(content.Title), len(content.Body))
	}
	combined := content.Title + "\n" + content.Body
	for _, want := range []string{
		"<!-- akritas:investigation_id=" + investigationID.String() + " -->",
		"Project", "Application", "Environment", "Repository", "Default branch",
		"Incident ID", incidentID.String(), "Fingerprint", "Severity", "Occurrences", "First seen", "Last seen",
		"## Observed Evidence", "### Evidence", "Observed excerpt",
		"## Investigation — QVAC Analysis", "Root Cause Status", "Root Cause / Hypothesis", "Confidence",
		"Resolution Status", "Investigation Summary", "Hypotheses", "Relevant Files", "Relevant Commits", "Recommended Actions",
	} {
		if !strings.Contains(combined, want) {
			t.Fatalf("auditable content missing %q", want)
		}
	}
	if strings.Index(content.Body, "## Observed Evidence") >= strings.Index(content.Body, "## Investigation — QVAC Analysis") {
		t.Fatal("observed Evidence must remain visibly separated before QVAC conclusions")
	}
}
