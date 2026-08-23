//go:build integration

package monitoring_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	dbadapter "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	incidentrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/incident"
	monitoringrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/monitoring"
	projectrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/project"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
	"gorm.io/gorm"
)

func TestMonitoringPersistenceIsDurableTransactionalAndSerialized(t *testing.T) {
	db := dbtest.ConnectContainer(t)
	ctx := context.Background()
	project := insertProject(t, db, time.Now().UTC())
	repository, err := monitoringrepo.New(db)
	if err != nil {
		t.Fatal(err)
	}
	transactor := dbadapter.NewTransactor(db)
	checkpoint, err := domain.NewMonitoringCheckpoint(uuid.New(), project, domain.InitialLogIngestionLast10000, project.CreatedAt)
	if err != nil || repository.CreateCheckpoint(ctx, checkpoint) != nil {
		t.Fatalf("create checkpoint: %v", err)
	}

	rollbackCause := errors.New("rollback monitoring unit")
	err = transactor.WithinTransaction(ctx, func(txctx context.Context) error {
		if _, lockErr := repository.LockProject(txctx, project.ID); lockErr != nil {
			return lockErr
		}
		locked, getErr := repository.GetCurrentCheckpoint(txctx, project.ID, true)
		if getErr != nil {
			return getErr
		}
		version := locked.Version
		locked.Advance(domain.MonitoringCursor{Timestamp: project.CreatedAt.Add(time.Second), ContentHash: "hash"}, project.CreatedAt.Add(time.Second))
		if updateErr := repository.UpdateCheckpoint(txctx, locked, version); updateErr != nil {
			return updateErr
		}
		incident := newIncident(t, project.ID, "AKR-1", project.CreatedAt)
		if createErr := repository.CreateIncident(txctx, incident); createErr != nil {
			return createErr
		}
		event := newEvent(t, incident.ID, project, "same-occurrence", project.CreatedAt)
		if createErr := repository.CreateLogEvent(txctx, event); createErr != nil {
			return createErr
		}
		return rollbackCause
	})
	if !errors.Is(err, rollbackCause) {
		t.Fatalf("transaction error = %v", err)
	}
	reloaded, err := repository.GetCurrentCheckpoint(ctx, project.ID, false)
	if err != nil || reloaded.Version != 1 || reloaded.CursorTimestamp != nil || !reloaded.InitialBackfillPending {
		t.Fatalf("checkpoint survived rollback: %+v, %v", reloaded, err)
	}
	var count int64
	if err := db.Table("incidents").Count(&count).Error; err != nil || count != 0 {
		t.Fatalf("incidents survived rollback: %d, %v", count, err)
	}

	baseline := newIncident(t, project.ID, "AKR-2", project.CreatedAt)
	if err := repository.CreateIncident(ctx, baseline); err != nil {
		t.Fatal(err)
	}
	first := newEvent(t, baseline.ID, project, "same-occurrence", project.CreatedAt)
	if err := repository.CreateLogEvent(ctx, first); err != nil {
		t.Fatal(err)
	}
	queries, _ := incidentrepo.New(db)
	gotIncident, err := queries.Get(ctx, baseline.ID)
	if err != nil || gotIncident.Project == nil || gotIncident.Project.Name != project.Name {
		t.Fatalf("incident detail projection = %+v, %v", gotIncident, err)
	}
	incidentPage, err := queries.List(ctx, paging.Params{Limit: 25, Filters: map[string]string{"project_id_eq": project.ID.String()}, Sort: []ukerpagination.SortExpression{{Field: "last_seen_at", Direction: ukerpagination.DirectionDesc}, {Field: "id", Direction: ukerpagination.DirectionDesc}}})
	if err != nil || incidentPage.Total != 1 || len(incidentPage.Items) != 1 {
		t.Fatalf("incident list = %+v, %v", incidentPage, err)
	}
	eventPage, err := queries.ListLogEvents(ctx, baseline.ID, paging.Params{Limit: 25, Sort: []ukerpagination.SortExpression{{Field: "timestamp", Direction: ukerpagination.DirectionDesc}, {Field: "id", Direction: ukerpagination.DirectionDesc}}})
	if err != nil || eventPage.Total != 1 || len(eventPage.Items) != 1 || !eventPage.Items[0].RawContextRedacted {
		t.Fatalf("log event list = %+v, %v", eventPage, err)
	}
	duplicate := newEvent(t, baseline.ID, project, "same-occurrence", project.CreatedAt)
	if err := repository.CreateLogEvent(ctx, duplicate); err == nil {
		t.Fatal("duplicate occurrence key was accepted")
	}

	occurredAt := project.CreatedAt.Add(time.Minute)
	var group sync.WaitGroup
	errorsCh := make(chan error, 2)
	for range 2 {
		group.Add(1)
		go func() {
			defer group.Done()
			errorsCh <- transactor.WithinTransaction(ctx, func(txctx context.Context) error {
				if _, lockErr := repository.LockProject(txctx, project.ID); lockErr != nil {
					return lockErr
				}
				incident, findErr := repository.FindGroupableIncident(txctx, project.ID, baseline.Fingerprint, occurredAt, 30*time.Minute)
				if findErr != nil {
					return findErr
				}
				if recordErr := incident.RecordOccurrence(project.ID, baseline.Fingerprint, occurredAt, 30*time.Minute); recordErr != nil {
					return recordErr
				}
				return repository.UpdateIncident(txctx, incident)
			})
		}()
	}
	group.Wait()
	close(errorsCh)
	for err := range errorsCh {
		if err != nil {
			t.Fatal(err)
		}
	}
	var saved domain.Incident
	if err := db.Table("incidents").Where("id = ?", baseline.ID).Take(&saved).Error; err != nil || saved.OccurrenceCount != 3 || !saved.LastSeenAt.Equal(occurredAt) {
		t.Fatalf("concurrent grouping = %+v, %v", saved, err)
	}

	projects, _ := projectrepo.New(db)
	expectedUpdatedAt := project.UpdatedAt
	disabled := project.MonitoringConfiguration.Clone()
	disabled.Enabled = false
	if err := project.ReplaceMonitoringConfiguration(disabled, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	if err := projects.Update(ctx, &project, expectedUpdatedAt); err != nil {
		t.Fatal(err)
	}
	if err := projects.Delete(ctx, project.ID); !errors.Is(err, domain.ErrProjectHasDependencies) {
		t.Fatalf("Project delete with history = %v", err)
	}
}

func insertProject(t *testing.T, db *gorm.DB, now time.Time) domain.Project {
	t.Helper()
	account, _ := domain.NewGitHubAccount(uuid.New(), "GitHub", domain.GitHubAccountOrganization, domain.GitHubAuthenticationPersonalAccessToken, "acme", domain.IntegrationStatusConnected, now)
	server, _ := domain.NewDokployServer(uuid.New(), "Dokploy", "https://dokploy.example.com", "server", domain.IntegrationStatusConnected, now)
	if err := db.Table("github_accounts").Create(account).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Table("dokploy_servers").Create(server).Error; err != nil {
		t.Fatal(err)
	}
	repository, _ := domain.NewGitHubRepository(account.ID, "42", "acme", "service", "main", true, "https://github.com/acme/service")
	application, _ := domain.NewDokployApplication(server.ID, "app", "instance", "service", "production", domain.DokployApplicationRunning)
	configuration, _ := domain.NewMonitoringConfiguration(true, []string{}, []string{}, 30*time.Minute, 2, 2)
	project, err := domain.NewProject(uuid.New(), "Service", "", repository, application, configuration, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Table("projects").Create(project).Error; err != nil {
		t.Fatal(err)
	}
	return *project
}

func newIncident(t *testing.T, projectID uuid.UUID, key string, timestamp time.Time) *domain.Incident {
	t.Helper()
	value, err := domain.NewIncident(uuid.New(), key, projectID, "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", domain.SeverityError, "failure", timestamp)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func newEvent(t *testing.T, incidentID uuid.UUID, project domain.Project, occurrenceKey string, timestamp time.Time) *domain.LogEvent {
	t.Helper()
	value, err := domain.NewLogEvent(uuid.New(), project.ID, timestamp, domain.SeverityError, "ERROR failure", "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", []string{"error_level"}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := value.AssociateOccurrence(incidentID, "app", "instance", occurrenceKey); err != nil {
		t.Fatal(err)
	}
	return value
}
