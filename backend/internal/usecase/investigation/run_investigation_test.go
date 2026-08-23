package investigation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func newRunFixture(t *testing.T, now time.Time) (*startDeps, uuid.UUID, uuid.UUID) {
	t.Helper()
	deps := newStartDeps()
	deps.now = now
	investigation, err := domain.NewInvestigation(uuid.New(), uuid.New(), now)
	if err != nil {
		t.Fatal(err)
	}
	operation, err := domain.NewOperation(uuid.New(), domain.OperationTypeInvestigation, nil, nil, nil, "queued", now)
	if err != nil {
		t.Fatal(err)
	}
	deps.investigations.findByIDResult = investigation
	deps.operations.findByIDResult = operation
	return deps, investigation.ID, operation.ID
}

func fixtureEvidence(t *testing.T, investigationID uuid.UUID, now time.Time) domain.Evidence {
	t.Helper()
	value, err := domain.NewEvidence(uuid.New(), investigationID, domain.EvidenceDeploymentMetadata, "summary", "content", now)
	if err != nil {
		t.Fatal(err)
	}
	return *value
}

func TestRunInvestigationCompletesOnRunnerSuccess(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	deps, investigationID, operationID := newRunFixture(t, now)
	deps.assembler.result = []domain.Evidence{fixtureEvidence(t, investigationID, now)}
	rootCause := domain.RootCauseIdentified
	resolution := domain.ResolutionFixable
	deps.runner.result = out.InvestigationRunResult{
		Summary: "root cause found", RootCause: "nil pointer", RootCauseStatus: rootCause, ResolutionStatus: resolution,
		Confidence: 0.9, Hypotheses: []string{"h1"}, RelevantFiles: []string{"main.go"}, RelevantCommits: []string{"abc123"}, RecommendedActions: []string{"patch it"},
	}

	if err := deps.runUseCase().Execute(context.Background(), investigationID, operationID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(deps.investigations.updated) != 3 {
		t.Fatalf("expected running, evidence count, then completed persisted, got %d updates", len(deps.investigations.updated))
	}
	evidenceUpdate := deps.investigations.updated[1]
	if evidenceUpdate.EvidenceCount != 1 {
		t.Fatalf("expected EvidenceCount to be persisted before the runner runs, got %+v", evidenceUpdate)
	}
	final := deps.investigations.updated[len(deps.investigations.updated)-1]
	if final.Status != domain.InvestigationStatusCompleted || final.Summary != "root cause found" || final.EvidenceCount != 1 {
		t.Fatalf("expected the investigation to complete with the runner result and keep EvidenceCount, got %+v", final)
	}

	if len(deps.operations.updated) != 2 {
		t.Fatalf("expected running then succeeded persisted, got %d updates", len(deps.operations.updated))
	}
	finalOp := deps.operations.updated[len(deps.operations.updated)-1]
	if finalOp.Status != domain.OperationStatusSucceeded {
		t.Fatalf("expected the operation to succeed, got %+v", finalOp)
	}

	if len(deps.evidence.created) != 1 {
		t.Fatalf("expected the assembled evidence to be persisted, got %d", len(deps.evidence.created))
	}
}

func TestRunInvestigationFailsOnRunnerError(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	deps, investigationID, operationID := newRunFixture(t, now)
	deps.runner.err = errors.New("qvac: not implemented yet, pending PB-028+")

	if err := deps.runUseCase().Execute(context.Background(), investigationID, operationID); err != nil {
		t.Fatalf("a failed investigation is a valid outcome, not a Go error: %v", err)
	}

	finalInvestigation := deps.investigations.updated[len(deps.investigations.updated)-1]
	if finalInvestigation.Status != domain.InvestigationStatusFailed || finalInvestigation.FailureUserMessage == "" {
		t.Fatalf("expected the investigation to fail with a message, got %+v", finalInvestigation)
	}

	finalOperation := deps.operations.updated[len(deps.operations.updated)-1]
	if finalOperation.Status != domain.OperationStatusFailed || finalOperation.UserMessage == "" {
		t.Fatalf("expected the operation to fail with a message, got %+v", finalOperation)
	}
}

func TestRunInvestigationKeepsAssembledEvidenceWhenRunnerLaterFails(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	deps, investigationID, operationID := newRunFixture(t, now)
	deps.assembler.result = []domain.Evidence{fixtureEvidence(t, investigationID, now), fixtureEvidence(t, investigationID, now)}
	deps.runner.err = errors.New("qvac: not implemented yet, pending PB-028+")

	if err := deps.runUseCase().Execute(context.Background(), investigationID, operationID); err != nil {
		t.Fatal(err)
	}

	if len(deps.evidence.created) != 2 {
		t.Fatalf("expected assembled evidence to survive a later runner failure, got %d", len(deps.evidence.created))
	}
	finalInvestigation := deps.investigations.updated[len(deps.investigations.updated)-1]
	if finalInvestigation.Status != domain.InvestigationStatusFailed || finalInvestigation.EvidenceCount != 2 {
		t.Fatalf("expected the failed investigation to keep its real EvidenceCount, got %+v", finalInvestigation)
	}
}

func TestRunInvestigationContinuesWhenAssemblerProducesNoEvidence(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	deps, investigationID, operationID := newRunFixture(t, now)
	rootCause := domain.RootCauseIdentified
	resolution := domain.ResolutionFixable
	deps.runner.result = out.InvestigationRunResult{
		Summary: "s", RootCause: "c", RootCauseStatus: rootCause, ResolutionStatus: resolution, Confidence: 0.5,
	}

	if err := deps.runUseCase().Execute(context.Background(), investigationID, operationID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(deps.evidence.created) != 0 {
		t.Fatalf("expected no evidence to be created, got %d", len(deps.evidence.created))
	}
	final := deps.investigations.updated[len(deps.investigations.updated)-1]
	if final.EvidenceCount != 0 {
		t.Fatalf("expected EvidenceCount to be 0 when the assembler produces nothing, got %+v", final)
	}
}

func TestRunInvestigationPropagatesAssemblerInfrastructureErrorWithoutInvokingRunner(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	deps, investigationID, operationID := newRunFixture(t, now)
	wantErr := errors.New("assembler unavailable")
	deps.assembler.err = wantErr

	err := deps.runUseCase().Execute(context.Background(), investigationID, operationID)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected the assembler error to propagate, got %v", err)
	}
	if deps.runner.calls != 0 {
		t.Fatal("runner must not be invoked when evidence assembly fails")
	}
	for _, update := range deps.investigations.updated {
		if update.Status == domain.InvestigationStatusFailed || update.Status == domain.InvestigationStatusCompleted {
			t.Fatal("an assembler infrastructure error must not route through failInvestigation")
		}
	}
}

func TestRunInvestigationPersistsRunningBeforeInvokingRunner(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	deps, investigationID, operationID := newRunFixture(t, now)
	deps.runner.err = errors.New("boom")

	if err := deps.runUseCase().Execute(context.Background(), investigationID, operationID); err != nil {
		t.Fatal(err)
	}

	if deps.investigations.updated[0].Status != domain.InvestigationStatusRunning {
		t.Fatal("expected the running state to be persisted before the runner is invoked")
	}
	if deps.operations.updated[0].Status != domain.OperationStatusRunning {
		t.Fatal("expected the running operation state to be persisted before the runner is invoked")
	}
}
