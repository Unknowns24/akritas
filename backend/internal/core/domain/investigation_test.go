package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestInvestigationLifecycleAndDefensiveCopies(t *testing.T) {
	t.Parallel()

	createdAt := time.Now().UTC()
	investigation, err := NewInvestigation(uuid.New(), uuid.New(), createdAt)
	if err != nil {
		t.Fatal(err)
	}
	if err := investigation.Start(createdAt.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	hypotheses := []string{"nil dependency"}
	files := []string{"internal/service.go"}
	commits := []string{"abc123"}
	actions := []string{"add guard"}
	if err := investigation.Complete(
		createdAt.Add(time.Minute),
		"summary",
		"nil dependency",
		RootCauseIdentified,
		ResolutionFixable,
		0.8,
		hypotheses,
		files,
		commits,
		actions,
	); err != nil {
		t.Fatal(err)
	}
	hypotheses[0], files[0], commits[0], actions[0] = "mutated", "mutated", "mutated", "mutated"
	if investigation.Hypotheses[0] == "mutated" || investigation.RelevantFiles[0] == "mutated" || investigation.RelevantCommits[0] == "mutated" || investigation.RecommendedActions[0] == "mutated" {
		t.Fatal("investigation retained caller-owned slices")
	}
	if err := investigation.Start(createdAt.Add(time.Hour)); !errors.Is(err, ErrInvestigationTransition) {
		t.Fatalf("completed investigation restarted: %v", err)
	}
}

func TestEvidenceIsAlwaysSanitized(t *testing.T) {
	t.Parallel()

	evidence, err := NewEvidence(uuid.New(), uuid.New(), EvidenceLogExcerpt, "safe summary", "redacted content", time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if !evidence.Redacted {
		t.Fatal("evidence must be marked redacted")
	}
}

func TestInvestigationFailureIsTerminal(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	investigation, _ := NewInvestigation(uuid.New(), uuid.New(), now)
	_ = investigation.Start(now.Add(time.Second))
	if err := investigation.Fail(now.Add(time.Minute), "No se pudo completar la investigación."); err != nil {
		t.Fatal(err)
	}
	if investigation.Status != InvestigationStatusFailed || investigation.FinishedAt == nil {
		t.Fatalf("unexpected failed investigation: %+v", investigation)
	}
	if err := investigation.Fail(now.Add(2*time.Minute), "again"); !errors.Is(err, ErrInvestigationTransition) {
		t.Fatalf("terminal investigation changed state: %v", err)
	}
}
