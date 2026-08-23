package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestProjectMonitoringTransitionsAndActiveAssociationGuard(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	repository := testRepository(t, uuid.New(), "1")
	application := testApplication(t, uuid.New(), "app-1")
	project, err := NewProject(uuid.New(), " Akritas ", " demo ", repository, application, DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatal(err)
	}
	if project.Name != "Akritas" || project.MonitoringStatus != MonitoringStatusDisabled {
		t.Fatalf("unexpected project: %+v", project)
	}

	enabled, err := NewMonitoringConfiguration(true, []string{"panic"}, []string{"healthcheck"}, 15*time.Minute, 8, 12)
	if err != nil {
		t.Fatal(err)
	}
	if err := project.ReplaceMonitoringConfiguration(enabled, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if project.MonitoringStatus != MonitoringStatusStarting || !project.MonitoringConfiguration.Enabled {
		t.Fatalf("not activated: %+v", project)
	}
	if err := project.ReadyForMonitoringEngine(); err != nil {
		t.Fatalf("not ready: %v", err)
	}
	if err := project.ReplaceIntegrations(repository, application, now.Add(2*time.Minute)); !errors.Is(err, ErrProjectMustBeDisabled) {
		t.Fatalf("active association change = %v", err)
	}

	disabled := enabled.Clone()
	disabled.Enabled = false
	if err := project.ReplaceMonitoringConfiguration(disabled, now.Add(3*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if project.MonitoringStatus != MonitoringStatusDisabled || len(project.MonitoringConfiguration.ErrorPatterns) != 1 || len(project.MonitoringConfiguration.IgnoredPatterns) != 1 {
		t.Fatalf("disable lost configuration: %+v", project)
	}
}

func TestProjectRejectsInconsistentMonitoringStateAndSnapshotIdentityChanges(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	repository := testRepository(t, uuid.New(), "1")
	application := testApplication(t, uuid.New(), "app-1")
	project, err := NewProject(uuid.New(), "Akritas", "", repository, application, DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatal(err)
	}
	project.MonitoringStatus = MonitoringStatusStarting
	if err := project.Validate(); !errors.Is(err, ErrInvalidProject) {
		t.Fatalf("inconsistent state = %v", err)
	}
	project.MonitoringStatus = MonitoringStatusDisabled
	other := testRepository(t, uuid.New(), "2")
	if err := project.RefreshIntegrationSnapshots(other, application, now.Add(time.Minute)); !errors.Is(err, ErrInvalidProject) {
		t.Fatalf("identity refresh = %v", err)
	}
}

func testRepository(t *testing.T, accountID uuid.UUID, identifier string) GitHubRepository {
	t.Helper()
	value, err := NewGitHubRepository(accountID, identifier, "Unknowns24", "akritas", "main", true, "https://github.com/Unknowns24/akritas")
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func testApplication(t *testing.T, serverID uuid.UUID, identifier string) DokploySource {
	t.Helper()
	value, err := NewDokployApplication(serverID, identifier, "instance-1", "Akritas", "production", DokployApplicationRunning)
	if err != nil {
		t.Fatal(err)
	}
	source, err := SourceFromApplication(value)
	if err != nil {
		t.Fatal(err)
	}
	return source
}
