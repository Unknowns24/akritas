package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidationResultLifecycle(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	result, err := NewValidationResult(uuid.New(), uuid.New(), ValidationTypeTest, "go test ./...", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := result.Start(now.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := result.Pass(now.Add(time.Minute), "all tests passed", "ok"); err != nil {
		t.Fatal(err)
	}
	if result.OutputRedacted || result.Status != ValidationStatusPassed {
		t.Fatalf("unexpected validation result: %+v", result)
	}
	if err := result.Fail(now.Add(2*time.Minute), "late", "late"); !errors.Is(err, ErrValidationTransition) {
		t.Fatalf("terminal validation changed state: %v", err)
	}
}

func TestRemediationRequiresPassingValidationBeforePullRequest(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	remediation, err := NewRemediation(uuid.New(), uuid.New(), now)
	if err != nil {
		t.Fatal(err)
	}
	if err := remediation.Start("codex/akr-1-fix", now.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	change, err := NewCodeChange("internal/service.go", CodeChangeModified, "@@ sanitized @@")
	if err != nil {
		t.Fatal(err)
	}
	if err := remediation.AddChange(change, now.Add(2*time.Second)); err != nil {
		t.Fatal(err)
	}
	validation, _ := NewValidationResult(uuid.New(), remediation.ID, ValidationTypeTest, "go test ./...", now)
	_ = validation.Start(now.Add(time.Second))
	_ = validation.Pass(now.Add(time.Minute), "passed", "ok")
	if err := remediation.AddValidationResult(*validation, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := remediation.MarkValidated(now.Add(2 * time.Minute)); err != nil {
		t.Fatal(err)
	}
	pr, _ := NewPullRequestReference(9, "https://example.com/pull/9", "owner/repo", "codex/akr-1-fix", now)
	if err := remediation.CreatePullRequest(pr, now.Add(3*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if remediation.Status != RemediationStatusPullRequestCreated {
		t.Fatalf("unexpected remediation state: %+v", remediation)
	}
}

func TestRemediationRejectsFailedValidation(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	remediation, _ := NewRemediation(uuid.New(), uuid.New(), now)
	_ = remediation.Start("codex/fix", now)
	validation, _ := NewValidationResult(uuid.New(), remediation.ID, ValidationTypeBuild, "go build ./...", now)
	_ = validation.Start(now)
	_ = validation.Fail(now.Add(time.Second), "build failed", "safe output")
	_ = remediation.AddValidationResult(*validation, now.Add(time.Second))
	if err := remediation.MarkValidated(now.Add(time.Minute)); !errors.Is(err, ErrRemediationTransition) {
		t.Fatalf("failed validation allowed remediation: %v", err)
	}
}

func TestRemediationFailureIsTerminal(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	remediation, _ := NewRemediation(uuid.New(), uuid.New(), now)
	_ = remediation.Start("codex/fix", now.Add(time.Second))
	if err := remediation.Fail("La validación no pudo completarse.", now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if remediation.Status != RemediationStatusFailed || remediation.FailureUserMessage == "" {
		t.Fatalf("unexpected failed remediation: %+v", remediation)
	}
	pr, _ := NewPullRequestReference(1, "https://example.com/pull/1", "owner/repo", "codex/fix", now)
	if err := remediation.CreatePullRequest(pr, now.Add(2*time.Minute)); !errors.Is(err, ErrRemediationTransition) {
		t.Fatalf("failed remediation created PR: %v", err)
	}
}
