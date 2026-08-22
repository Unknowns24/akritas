package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestProjectRenameReplaceIntegrationsAndMonitoring(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	project := newTestProject(t, now)
	later := now.Add(time.Minute)

	if err := project.Rename("renamed", "updated description", later); err != nil {
		t.Fatalf("rename rejected: %v", err)
	}
	if project.Name != "renamed" || project.Description != "updated description" || project.UpdatedAt != later {
		t.Fatalf("rename did not persist identity: %+v", project)
	}

	accountID := uuid.New()
	repository, err := NewGitHubRepository(accountID, "repo-2", "Unknowns24", "sentinel", "develop", true, "https://github.com/Unknowns24/sentinel")
	if err != nil {
		t.Fatal(err)
	}
	serverID := uuid.New()
	application, err := NewDokployApplication(serverID, "app-2", "instance-2", "worker", "staging", DokployApplicationStopped)
	if err != nil {
		t.Fatal(err)
	}
	if err := project.ReplaceIntegrations(repository, application, later.Add(time.Second)); err != nil {
		t.Fatalf("replace integrations rejected: %v", err)
	}
	if project.GitHubRepository.RepositoryIdentifier != "repo-2" || project.DokployApplication.ApplicationIdentifier != "app-2" {
		t.Fatalf("integrations were not replaced: %+v", project)
	}

	enabled, err := NewMonitoringConfiguration(true, []string{`database .* down`}, []string{`healthcheck`}, time.Hour, 5, 8)
	if err != nil {
		t.Fatal(err)
	}
	if err := project.ReplaceMonitoringConfiguration(enabled, later.Add(2*time.Second)); err != nil {
		t.Fatalf("enable monitoring rejected: %v", err)
	}
	if project.MonitoringStatus != MonitoringStatusStarting || !project.MonitoringConfiguration.Enabled {
		t.Fatalf("disabled project must start monitoring: %+v", project)
	}

	patterns := []string{`original`}
	config := DefaultMonitoringConfiguration()
	config.ErrorPatterns = patterns
	if err := project.ReplaceMonitoringConfiguration(config, later.Add(3*time.Second)); err != nil {
		t.Fatal(err)
	}
	if project.MonitoringStatus != MonitoringStatusDisabled {
		t.Fatalf("disabled configuration must disable monitoring, got %s", project.MonitoringStatus)
	}
	patterns[0] = "mutated"
	if project.MonitoringConfiguration.ErrorPatterns[0] == "mutated" {
		t.Fatal("replace monitoring retained caller-owned slices")
	}

	project.MonitoringStatus = MonitoringStatusDegraded
	reenable, err := NewMonitoringConfiguration(true, nil, nil, time.Minute, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := project.ReplaceMonitoringConfiguration(reenable, later.Add(4*time.Second)); err != nil {
		t.Fatal(err)
	}
	if project.MonitoringStatus != MonitoringStatusDegraded {
		t.Fatalf("enabled degraded project must keep status, got %s", project.MonitoringStatus)
	}

	if err := project.Rename("", "", later); !errors.Is(err, ErrInvalidProject) {
		t.Fatalf("expected invalid rename, got %v", err)
	}
}

func TestProjectEnableDisableKeepsPatternsAndRequiresValidRefs(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	project := newTestProject(t, now)
	later := now.Add(time.Minute)
	if err := project.ReadyForMonitoringEngine(); !errors.Is(err, ErrInvalidMonitoringConfiguration) {
		t.Fatalf("disabled project must not be ready, got %v", err)
	}

	enabled, err := NewMonitoringConfiguration(true, []string{`panic`}, []string{`healthcheck`}, time.Hour, 5, 8)
	if err != nil {
		t.Fatal(err)
	}
	if err := project.ReplaceMonitoringConfiguration(enabled, later); err != nil {
		t.Fatalf("enable rejected: %v", err)
	}
	if project.MonitoringStatus != MonitoringStatusStarting || !project.MonitoringConfiguration.Enabled {
		t.Fatalf("disabled project must become starting: %+v", project)
	}
	if err := project.ReadyForMonitoringEngine(); err != nil {
		t.Fatalf("enabled project with valid refs must be ready: %v", err)
	}

	disabled, err := NewMonitoringConfiguration(false, []string{`panic`}, []string{`healthcheck`}, time.Hour, 5, 8)
	if err != nil {
		t.Fatal(err)
	}
	if err := project.ReplaceMonitoringConfiguration(disabled, later.Add(time.Second)); err != nil {
		t.Fatalf("disable rejected: %v", err)
	}
	if project.MonitoringStatus != MonitoringStatusDisabled || project.MonitoringConfiguration.Enabled {
		t.Fatalf("disable must set disabled without wiping config: %+v", project)
	}
	if len(project.MonitoringConfiguration.ErrorPatterns) != 1 || project.MonitoringConfiguration.ErrorPatterns[0] != `panic` {
		t.Fatalf("disable must keep error patterns: %+v", project.MonitoringConfiguration.ErrorPatterns)
	}
	if len(project.MonitoringConfiguration.IgnoredPatterns) != 1 || project.MonitoringConfiguration.IgnoredPatterns[0] != `healthcheck` {
		t.Fatalf("disable must keep ignored patterns: %+v", project.MonitoringConfiguration.IgnoredPatterns)
	}

	if err := project.ReplaceMonitoringConfiguration(enabled, later.Add(2*time.Second)); err != nil {
		t.Fatalf("re-enable rejected: %v", err)
	}
	if project.MonitoringStatus != MonitoringStatusStarting {
		t.Fatalf("disabled→enabled must become starting, got %s", project.MonitoringStatus)
	}

	project.MonitoringStatus = MonitoringStatusDegraded
	if err := project.ReplaceMonitoringConfiguration(enabled, later.Add(3*time.Second)); err != nil {
		t.Fatal(err)
	}
	if project.MonitoringStatus != MonitoringStatusDegraded {
		t.Fatalf("re-enable must keep degraded, got %s", project.MonitoringStatus)
	}

	if err := project.ReplaceMonitoringConfiguration(disabled, later.Add(4*time.Second)); err != nil {
		t.Fatal(err)
	}
	brokenRepo := *project
	brokenRepo.GitHubRepository.GitHubAccountID = uuid.Nil
	if err := brokenRepo.ReplaceMonitoringConfiguration(enabled, later.Add(5*time.Second)); !errors.Is(err, ErrInvalidGitHubRepository) {
		t.Fatalf("enable with invalid repo must fail, got %v", err)
	}
	if brokenRepo.MonitoringStatus != MonitoringStatusDisabled || brokenRepo.MonitoringConfiguration.Enabled {
		t.Fatalf("failed enable must leave project disabled: %+v", brokenRepo)
	}

	brokenApp := *project
	brokenApp.DokployApplication.ApplicationIdentifier = ""
	if err := brokenApp.ReplaceMonitoringConfiguration(enabled, later.Add(6*time.Second)); !errors.Is(err, ErrInvalidDokployApplication) {
		t.Fatalf("enable with invalid app must fail, got %v", err)
	}
	if brokenApp.MonitoringStatus != MonitoringStatusDisabled || brokenApp.MonitoringConfiguration.Enabled {
		t.Fatalf("failed enable must leave project disabled: %+v", brokenApp)
	}
}

func TestAllBuiltInDetectionRulesAreEnabledCatalog(t *testing.T) {
	t.Parallel()

	rules := AllBuiltInDetectionRules()
	if len(rules) != 7 {
		t.Fatalf("expected 7 built-in rules, got %d", len(rules))
	}
	seen := map[DetectionRuleCode]struct{}{}
	for _, rule := range rules {
		if err := rule.Validate(); err != nil {
			t.Fatalf("catalog rule %s is invalid: %v", rule.Code, err)
		}
		if !rule.Enabled {
			t.Fatalf("catalog rule %s must be enabled", rule.Code)
		}
		if _, exists := seen[rule.Code]; exists {
			t.Fatalf("duplicate catalog rule %s", rule.Code)
		}
		seen[rule.Code] = struct{}{}
	}
	required := []DetectionRuleCode{
		DetectionRuleErrorLevel, DetectionRuleFatalLevel, DetectionRulePanic, DetectionRuleStackTrace,
		DetectionRuleHTTP5xx, DetectionRuleProcessCrash, DetectionRuleContainerRestart,
	}
	for _, code := range required {
		if _, ok := seen[code]; !ok {
			t.Fatalf("catalog missing %s", code)
		}
	}
}

func newTestProject(t *testing.T, now time.Time) *Project {
	t.Helper()
	account, err := NewGitHubAccount(uuid.New(), "Akritas", GitHubAccountOrganization, GitHubAuthenticationGitHubApp, "unknowns24", IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	repository, err := NewGitHubRepository(account.ID, "repo-1", "Unknowns24", "akritas", "main", false, "https://github.com/Unknowns24/akritas")
	if err != nil {
		t.Fatal(err)
	}
	server, err := NewDokployServer(uuid.New(), "demo", "https://dokploy.example.com", "server-1", IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	application, err := NewDokployApplication(server.ID, "app-1", "instance-1", "api", "production", DokployApplicationRunning)
	if err != nil {
		t.Fatal(err)
	}
	project, err := NewProject(uuid.New(), "sentinel-api", "demo", repository, application, DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatal(err)
	}
	return project
}
