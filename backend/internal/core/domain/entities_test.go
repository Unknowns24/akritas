package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPublishedEnums(t *testing.T) {
	t.Parallel()
	validChecks := []struct {
		name string
		err  error
	}{
		{"integration", IntegrationStatusConnected.Validate()},
		{"connection test", ConnectionTestConnected.Validate()},
		{"github account type", GitHubAccountPersonal.Validate()},
		{"github auth", GitHubAuthenticationPersonalAccessToken.Validate()},
		{"monitoring", MonitoringStatusMonitoring.Validate()},
		{"health", ProjectHealthHealthy.Validate()},
		{"severity", SeverityCritical.Validate()},
		{"incident phase", IncidentPhasePublishingIssue.Validate()},
		{"terminal outcome", TerminalOutcomeRemediationFailed.Validate()},
		{"investigation", InvestigationStatusRunning.Validate()},
		{"root cause", RootCauseSuspected.Validate()},
		{"resolution", ResolutionRequiresHuman.Validate()},
		{"evidence", EvidenceCodeLocation.Validate()},
		{"remediation", RemediationStatusValidated.Validate()},
		{"validation status", ValidationStatusPassed.Validate()},
		{"validation type", ValidationTypeStaticAnalysis.Validate()},
		{"change type", CodeChangeDeleted.Validate()},
		{"detection rule", DetectionRuleContainerRestart.Validate()},
	}
	for _, check := range validChecks {
		if check.err != nil {
			t.Fatalf("%s rejected a published value: %v", check.name, check.err)
		}
	}

	checks := []struct {
		name string
		err  error
	}{
		{"integration", IntegrationStatus("bad").Validate()},
		{"connection test", ConnectionTestStatus("bad").Validate()},
		{"github account type", GitHubAccountType("bad").Validate()},
		{"github auth", GitHubAuthenticationMethod("bad").Validate()},
		{"monitoring", MonitoringStatus("bad").Validate()},
		{"health", ProjectHealthStatus("bad").Validate()},
		{"severity", Severity("bad").Validate()},
		{"incident phase", IncidentPhase("bad").Validate()},
		{"terminal outcome", TerminalOutcome("bad").Validate()},
		{"investigation", InvestigationStatus("bad").Validate()},
		{"root cause", RootCauseStatus("bad").Validate()},
		{"resolution", ResolutionStatus("bad").Validate()},
		{"evidence", EvidenceType("bad").Validate()},
		{"remediation", RemediationStatus("bad").Validate()},
		{"validation status", ValidationStatus("bad").Validate()},
		{"validation type", ValidationType("bad").Validate()},
		{"change type", CodeChangeType("bad").Validate()},
		{"detection rule", DetectionRuleCode("bad").Validate()},
	}
	for _, check := range checks {
		if check.err == nil {
			t.Fatalf("%s accepted an unknown value", check.name)
		}
	}
}

func TestBuiltInDetectionRuleIsAlwaysEnabled(t *testing.T) {
	t.Parallel()

	rule, err := NewBuiltInDetectionRule(DetectionRulePanic, "Panic")
	if err != nil {
		t.Fatal(err)
	}
	if !rule.Enabled {
		t.Fatal("built-in rule must be enabled")
	}
	rule.Enabled = false
	if !errors.Is(rule.Validate(), ErrInvalidDetectionRule) {
		t.Fatal("disabled built-in rule must be rejected")
	}
}

func TestAdministratorConstructor(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	administrator, err := NewAdministrator(uuid.New(), "admin@example.com", "Akritas Admin", now)
	if err != nil {
		t.Fatalf("valid administrator rejected: %v", err)
	}
	if administrator.UpdatedAt != now {
		t.Fatal("administrator timestamps were not initialized together")
	}
	if _, err := NewAdministrator(uuid.Nil, "invalid", "", now); !errors.Is(err, ErrInvalidAdministrator) {
		t.Fatalf("expected invalid administrator error, got %v", err)
	}
}

func TestLogEventIsSanitizedAndCopiesContext(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	record, err := NewSanitizedLogRecord(now, LogStreamStderr, "panic: redacted")
	if err != nil {
		t.Fatal(err)
	}
	rules := []string{"panic"}
	context := []SanitizedLogRecord{record}
	event, err := NewLogEvent(uuid.New(), uuid.New(), now, SeverityCritical, "panic", "sha256:abcdef", rules, context, nil)
	if err != nil {
		t.Fatal(err)
	}
	rules[0] = "mutated"
	context[0].Message = "mutated"
	if event.DetectionRules[0] == "mutated" || event.ContextBefore[0].Message == "mutated" {
		t.Fatal("log event retained caller-owned slices")
	}
	if !event.RawContextRedacted || !event.ContextBefore[0].Redacted {
		t.Fatal("log event context must be redacted")
	}
	if _, err := NewSanitizedLogRecord(time.Time{}, LogStream("bad"), ""); !errors.Is(err, ErrInvalidSanitizedLogRecord) {
		t.Fatalf("expected invalid sanitized record error, got %v", err)
	}
}

func TestIntegrationAndProjectConstructors(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
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
	config := DefaultMonitoringConfiguration()
	config.ErrorPatterns = []string{"original"}
	project, err := NewProject(uuid.New(), "sentinel-api", "demo", repository, application, config, now)
	if err != nil {
		t.Fatal(err)
	}
	if project.MonitoringStatus != MonitoringStatusDisabled || project.HealthStatus != ProjectHealthUnknown {
		t.Fatalf("unexpected initial project state: %+v", project)
	}
	config.ErrorPatterns[0] = "mutated"
	if project.MonitoringConfiguration.ErrorPatterns[0] == "mutated" {
		t.Fatal("project retained caller-owned monitoring slices")
	}

	_, err = NewProject(uuid.Nil, "", "", repository, application, config, now)
	if !errors.Is(err, ErrInvalidProject) {
		t.Fatalf("expected invalid project error, got %v", err)
	}
}
