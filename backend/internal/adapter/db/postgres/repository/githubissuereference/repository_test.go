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
	"gorm.io/gorm"
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

func insertReferenceFixture(t *testing.T, db *gorm.DB, now time.Time) (uuid.UUID, uuid.UUID, uuid.UUID) {
	t.Helper()
	accountID := uuid.New()
	serverID := uuid.New()
	projectID := uuid.New()
	incidentID := uuid.New()
	firstInvestigationID := uuid.New()
	secondInvestigationID := uuid.New()

	statements := []struct {
		sql  string
		args []any
	}{
		{
			sql: `INSERT INTO github_accounts (
				id, display_name, account_type, authentication_method, account_identifier,
				authentication_status, credential_configured, repository_count, manage_url,
				created_at, updated_at
			) VALUES (?, ?, 'organization', 'personal_access_token', ?, 'connected', true, 1, '', ?, ?)`,
			args: []any{accountID, "Acme", "acme", now, now},
		},
		{
			sql: `INSERT INTO dokploy_servers (
				id, name, base_url, server_identifier, connection_status, credential_configured,
				application_count, created_at, updated_at
			) VALUES (?, ?, ?, ?, 'connected', true, 1, ?, ?)`,
			args: []any{serverID, "Dokploy", "https://dokploy.example.test", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", now, now},
		},
		{
			sql: `INSERT INTO projects (
				id, name, description, monitoring_status, health_status, github_account_id,
				repository_identifier, repository_owner, repository_name, repository_full_name,
				default_branch, repository_private, repository_html_url, dokploy_server_id,
				application_identifier, instance_identifier, application_display_name,
				application_environment, application_status, monitoring_enabled,
				grouping_window_ns, context_before, context_after, created_at, updated_at
			) VALUES (?, ?, '', 'disabled', 'unknown', ?, 'repo-1', 'acme', 'service',
				'acme/service', 'main', true, 'https://github.com/acme/service', ?,
				'app-1', 'instance-1', 'Service', 'prod', 'running', false,
				300000000000, 10, 10, ?, ?)`,
			args: []any{projectID, "Project " + projectID.String(), accountID, serverID, now, now},
		},
		{
			sql: `INSERT INTO incidents (
				id, key, project_id, fingerprint, severity, phase, first_seen_at, last_seen_at,
				occurrence_count, title, summary
			) VALUES (?, 'INC-1', ?, 'fingerprint', 'error', 'publishing_issue', ?, ?, 2, 'Incident', '')`,
			args: []any{incidentID, projectID, now.Add(-time.Hour), now},
		},
		{
			sql: `INSERT INTO investigations (
				id, incident_id, status, created_at, started_at, finished_at, summary,
				root_cause, root_cause_status, resolution_status, confidence, evidence_count
			) VALUES (?, ?, 'completed', ?, ?, ?, 'summary', 'cause', 'identified', 'fixable', 0.90, 0)`,
			args: []any{firstInvestigationID, incidentID, now.Add(-30 * time.Minute), now.Add(-29 * time.Minute), now.Add(-20 * time.Minute)},
		},
		{
			sql: `INSERT INTO investigations (
				id, incident_id, status, created_at, started_at, finished_at, summary,
				root_cause, root_cause_status, resolution_status, confidence, evidence_count
			) VALUES (?, ?, 'completed', ?, ?, ?, 'summary', 'cause', 'identified', 'requires_human', 0.75, 0)`,
			args: []any{secondInvestigationID, incidentID, now.Add(-10 * time.Minute), now.Add(-9 * time.Minute), now.Add(-5 * time.Minute)},
		},
	}
	for _, statement := range statements {
		if err := db.Exec(statement.sql, statement.args...).Error; err != nil {
			t.Fatal(err)
		}
	}
	return incidentID, firstInvestigationID, secondInvestigationID
}
