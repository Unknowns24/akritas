package domain

import (
	"errors"
	"testing"
	"time"
)

func TestDefaultMonitoringConfiguration(t *testing.T) {
	t.Parallel()

	config := DefaultMonitoringConfiguration()
	if config.Enabled || len(config.ErrorPatterns) != 0 || len(config.IgnoredPatterns) != 0 {
		t.Fatal("monitoring must be disabled with empty patterns by default")
	}
	if config.GroupingWindow != 30*time.Minute || config.ContextBefore != 20 || config.ContextAfter != 20 {
		t.Fatalf("unexpected monitoring defaults: %+v", config)
	}
	if err := config.Validate(); err != nil {
		t.Fatalf("default configuration must be valid: %v", err)
	}
}

func TestMonitoringConfigurationValidationAndDefensiveCopies(t *testing.T) {
	t.Parallel()

	errorsIn := []string{`database .* unavailable`}
	ignoredIn := []string{`expected healthcheck failure`}
	config, err := NewMonitoringConfiguration(true, errorsIn, ignoredIn, time.Minute, 0, 1000)
	if err != nil {
		t.Fatalf("valid configuration rejected: %v", err)
	}
	errorsIn[0] = "mutated"
	ignoredIn[0] = "mutated"
	if config.ErrorPatterns[0] == "mutated" || config.IgnoredPatterns[0] == "mutated" {
		t.Fatal("constructor retained caller-owned slices")
	}

	invalid := []MonitoringConfiguration{
		{GroupingWindow: 0},
		{GroupingWindow: time.Minute, ContextBefore: -1},
		{GroupingWindow: time.Minute, ContextAfter: 1001},
		{GroupingWindow: time.Minute, ErrorPatterns: []string{"["}},
		{GroupingWindow: time.Minute, IgnoredPatterns: []string{""}},
	}
	for _, candidate := range invalid {
		if err := candidate.Validate(); !errors.Is(err, ErrInvalidMonitoringConfiguration) {
			t.Fatalf("expected monitoring validation error for %+v, got %v", candidate, err)
		}
	}
}

func TestAutomationPolicyDependencies(t *testing.T) {
	t.Parallel()

	valid := []AutomationPolicy{
		{},
		{AutomaticInvestigation: true},
		{AutomaticInvestigation: true, AutomaticRemediation: true},
		DefaultAutomationPolicy(),
	}
	for _, policy := range valid {
		if err := policy.Validate(); err != nil {
			t.Fatalf("valid policy rejected: %+v: %v", policy, err)
		}
	}

	invalid := []AutomationPolicy{
		{AutomaticRemediation: true},
		{AutomaticPullRequest: true},
		{AutomaticInvestigation: true, AutomaticPullRequest: true},
	}
	for _, policy := range invalid {
		if err := policy.Validate(); !errors.Is(err, ErrInvalidAutomationPolicy) {
			t.Fatalf("expected automation validation error for %+v, got %v", policy, err)
		}
	}
}
