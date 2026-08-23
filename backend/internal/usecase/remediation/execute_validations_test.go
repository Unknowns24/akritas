package remediation

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func startedRemediation(t *testing.T, remediations *fakeRemediationStore, id, incidentID uuid.UUID) {
	t.Helper()
	value, err := domain.NewRemediation(id, incidentID, fixedNow()())
	if err != nil {
		t.Fatalf("NewRemediation: %v", err)
	}
	if err := value.Start("akritas/remediation/"+id.String(), fixedNow()()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	remediations.byID[id] = *value
}

func TestExecuteRemediationValidationsNotFound(t *testing.T) {
	remediations := newFakeRemediationStore()
	uc := newTestUseCase(remediations, &fakeRepositoryWorkspace{}, newFakeValidationRunner(), &fakeValidationResultStore{}, &fakeWorkspaceInspector{})

	_, _, err := uc.ExecuteRemediationValidations(context.Background(), portsin.ExecuteRemediationValidationsCommand{
		RemediationID: uuid.New(), WorkspacePath: "/workspace",
	})
	if !errors.Is(err, domain.ErrRemediationNotFound) {
		t.Fatalf("expected ErrRemediationNotFound, got %v", err)
	}
}

func TestExecuteRemediationValidationsWrongStatus(t *testing.T) {
	remediations := newFakeRemediationStore()
	id, incidentID := uuid.New(), uuid.New()
	planned, err := domain.NewRemediation(id, incidentID, fixedNow()())
	if err != nil {
		t.Fatalf("NewRemediation: %v", err)
	}
	remediations.byID[id] = *planned // still "planned", never Start()-ed

	runner := newFakeValidationRunner()
	uc := newTestUseCase(remediations, &fakeRepositoryWorkspace{}, runner, &fakeValidationResultStore{}, &fakeWorkspaceInspector{has: map[string]bool{"go.mod": true}})

	_, _, err = uc.ExecuteRemediationValidations(context.Background(), portsin.ExecuteRemediationValidationsCommand{
		RemediationID: id, WorkspacePath: "/workspace",
	})
	if !errors.Is(err, domain.ErrRemediationTransition) {
		t.Fatalf("expected ErrRemediationTransition, got %v", err)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("expected no validations to run for a remediation not in_progress, got %d calls", len(runner.calls))
	}
}

func TestExecuteRemediationValidationsUnsupportedStack(t *testing.T) {
	remediations := newFakeRemediationStore()
	id, incidentID := uuid.New(), uuid.New()
	startedRemediation(t, remediations, id, incidentID)

	runner := newFakeValidationRunner()
	results := &fakeValidationResultStore{}
	uc := newTestUseCase(remediations, &fakeRepositoryWorkspace{}, runner, results, &fakeWorkspaceInspector{has: map[string]bool{}})

	_, _, err := uc.ExecuteRemediationValidations(context.Background(), portsin.ExecuteRemediationValidationsCommand{
		RemediationID: id, WorkspacePath: "/workspace",
	})
	if !errors.Is(err, domain.ErrValidationStackUnsupported) {
		t.Fatalf("expected ErrValidationStackUnsupported, got %v", err)
	}
	if len(runner.calls) != 0 || len(results.created) != 0 {
		t.Fatalf("expected zero execution/persistence for an unsupported stack, got %d runs, %d persisted", len(runner.calls), len(results.created))
	}
}

func TestExecuteRemediationValidationsAllStepsRunRegardlessOfEarlierFailures(t *testing.T) {
	remediations := newFakeRemediationStore()
	id, incidentID := uuid.New(), uuid.New()
	startedRemediation(t, remediations, id, incidentID)

	runner := newFakeValidationRunner()
	runner.results[portsout.ValidationCommandGoBuild] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 1, Stdout: "build failed TOKEN=secret-value"}
	runner.results[portsout.ValidationCommandGoVet] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 0}
	runner.results[portsout.ValidationCommandGoTest] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 0}

	results := &fakeValidationResultStore{}
	uc := newTestUseCase(remediations, &fakeRepositoryWorkspace{}, runner, results, &fakeWorkspaceInspector{has: map[string]bool{"go.mod": true}})

	remediationOut, resultsOut, err := uc.ExecuteRemediationValidations(context.Background(), portsin.ExecuteRemediationValidationsCommand{
		RemediationID: id, WorkspacePath: "/workspace",
	})
	if err != nil {
		t.Fatalf("ExecuteRemediationValidations: %v", err)
	}
	if len(runner.calls) != 3 {
		t.Fatalf("expected all 3 steps to run despite the build failure, got %d calls", len(runner.calls))
	}
	if len(resultsOut) != 3 || len(results.created) != 3 {
		t.Fatalf("expected 3 persisted results, got %d returned, %d persisted", len(resultsOut), len(results.created))
	}
	passed, failed := 0, 0
	for _, r := range resultsOut {
		switch r.Status {
		case domain.ValidationStatusPassed:
			passed++
		case domain.ValidationStatusFailed:
			failed++
		default:
			t.Fatalf("unexpected terminal status %v", r.Status)
		}
		if r.OutputRedacted && r.Status != domain.ValidationStatusFailed {
			t.Fatal("only the secret-bearing failed validation should be marked redacted")
		}
	}
	if !resultsOut[0].OutputRedacted || strings.Contains(resultsOut[0].OutputExcerpt, "secret-value") {
		t.Fatalf("failed validation output was not truthfully redacted: %+v", resultsOut[0])
	}
	if passed != 2 || failed != 1 {
		t.Fatalf("expected 2 passed and 1 failed, got %d passed, %d failed", passed, failed)
	}
	if remediationOut.Status != domain.RemediationStatusFailed || len(remediations.updated) != 1 {
		t.Fatalf("expected Remediation to fail after failed validation, got remediation=%+v updates=%d", remediationOut, len(remediations.updated))
	}
}

func TestExecuteRemediationValidationsRunnerErrorDistinguishableFromTestFailure(t *testing.T) {
	remediations := newFakeRemediationStore()
	id, incidentID := uuid.New(), uuid.New()
	startedRemediation(t, remediations, id, incidentID)

	runner := newFakeValidationRunner()
	runner.errs[portsout.ValidationCommandGoBuild] = errors.New("could not start process")
	runner.results[portsout.ValidationCommandGoVet] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 1, Stdout: "vet found a real issue"}
	runner.results[portsout.ValidationCommandGoTest] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 0}

	results := &fakeValidationResultStore{}
	uc := newTestUseCase(remediations, &fakeRepositoryWorkspace{}, runner, results, &fakeWorkspaceInspector{has: map[string]bool{"go.mod": true}})

	_, resultsOut, err := uc.ExecuteRemediationValidations(context.Background(), portsin.ExecuteRemediationValidationsCommand{
		RemediationID: id, WorkspacePath: "/workspace",
	})
	if err != nil {
		t.Fatalf("ExecuteRemediationValidations: %v", err)
	}

	var buildResult, vetResult domain.ValidationResult
	for _, r := range resultsOut {
		switch r.Type {
		case domain.ValidationTypeBuild:
			buildResult = r
		case domain.ValidationTypeStaticAnalysis:
			vetResult = r
		}
	}
	if buildResult.Status != domain.ValidationStatusFailed || vetResult.Status != domain.ValidationStatusFailed {
		t.Fatalf("expected both build and vet to be recorded as failed, got build=%v vet=%v", buildResult.Status, vetResult.Status)
	}
	if buildResult.Summary == vetResult.Summary {
		t.Fatalf("expected a runner-execution-error summary distinguishable from a genuine validation-failure summary, both were %q", buildResult.Summary)
	}
}

func TestExecuteRemediationValidationsTimeoutDistinguishableSummary(t *testing.T) {
	remediations := newFakeRemediationStore()
	id, incidentID := uuid.New(), uuid.New()
	startedRemediation(t, remediations, id, incidentID)

	runner := newFakeValidationRunner()
	runner.results[portsout.ValidationCommandGoBuild] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 0}
	runner.results[portsout.ValidationCommandGoVet] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 0}
	runner.results[portsout.ValidationCommandGoTest] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeTimedOut}

	results := &fakeValidationResultStore{}
	uc := newTestUseCase(remediations, &fakeRepositoryWorkspace{}, runner, results, &fakeWorkspaceInspector{has: map[string]bool{"go.mod": true}})

	_, resultsOut, err := uc.ExecuteRemediationValidations(context.Background(), portsin.ExecuteRemediationValidationsCommand{
		RemediationID: id, WorkspacePath: "/workspace",
	})
	if err != nil {
		t.Fatalf("ExecuteRemediationValidations: %v", err)
	}

	var testResult domain.ValidationResult
	for _, r := range resultsOut {
		if r.Type == domain.ValidationTypeTest {
			testResult = r
		}
	}
	if testResult.Status != domain.ValidationStatusFailed {
		t.Fatalf("expected a timed-out step to be recorded as failed, got %v", testResult.Status)
	}
	if !strings.Contains(strings.ToLower(testResult.Summary), "tiempo") {
		t.Fatalf("expected a timeout-specific summary, got %q", testResult.Summary)
	}
}

func TestExecuteRemediationValidationsTruncatesOversizedOutput(t *testing.T) {
	remediations := newFakeRemediationStore()
	id, incidentID := uuid.New(), uuid.New()
	startedRemediation(t, remediations, id, incidentID)

	huge := strings.Repeat("x", 200000)
	runner := newFakeValidationRunner()
	runner.results[portsout.ValidationCommandGoBuild] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 0, Stdout: huge}
	runner.results[portsout.ValidationCommandGoVet] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 0}
	runner.results[portsout.ValidationCommandGoTest] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 0}

	results := &fakeValidationResultStore{}
	uc := newTestUseCase(remediations, &fakeRepositoryWorkspace{}, runner, results, &fakeWorkspaceInspector{has: map[string]bool{"go.mod": true}})

	_, resultsOut, err := uc.ExecuteRemediationValidations(context.Background(), portsin.ExecuteRemediationValidationsCommand{
		RemediationID: id, WorkspacePath: "/workspace",
	})
	if err != nil {
		t.Fatalf("ExecuteRemediationValidations: %v", err)
	}
	for _, r := range resultsOut {
		if err := r.Validate(); err != nil {
			t.Fatalf("expected every persisted result to satisfy domain.Validate(), got %v for %+v", err, r)
		}
		if len(r.OutputExcerpt) > 50000 {
			t.Fatalf("expected OutputExcerpt to be capped at 50000 bytes, got %d", len(r.OutputExcerpt))
		}
	}
}

func TestExecuteRemediationValidationsStoreFailureMidPlanKeepsPriorResults(t *testing.T) {
	remediations := newFakeRemediationStore()
	id, incidentID := uuid.New(), uuid.New()
	startedRemediation(t, remediations, id, incidentID)

	runner := newFakeValidationRunner()
	runner.results[portsout.ValidationCommandGoBuild] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 0}
	runner.results[portsout.ValidationCommandGoVet] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 0}
	runner.results[portsout.ValidationCommandGoTest] = portsout.ExecutionResult{Outcome: portsout.ExecutionOutcomeCompleted, ExitCode: 0}

	results := &failNthValidationResultStore{failAt: 2}
	uc := newTestUseCase(remediations, &fakeRepositoryWorkspace{}, runner, results, &fakeWorkspaceInspector{has: map[string]bool{"go.mod": true}})

	_, _, err := uc.ExecuteRemediationValidations(context.Background(), portsin.ExecuteRemediationValidationsCommand{
		RemediationID: id, WorkspacePath: "/workspace",
	})
	if err == nil {
		t.Fatal("expected an error from the second Create call")
	}
	if len(results.created) != 1 {
		t.Fatalf("expected the first result to remain persisted despite the later failure, got %d", len(results.created))
	}
}

type failNthValidationResultStore struct {
	calls   int
	failAt  int
	created []domain.ValidationResult
}

func (f *failNthValidationResultStore) Create(ctx context.Context, value *domain.ValidationResult) error {
	f.calls++
	if f.calls == f.failAt {
		return errors.New("db down")
	}
	f.created = append(f.created, *value)
	return nil
}

func (f *failNthValidationResultStore) ListByRemediation(ctx context.Context, remediationID uuid.UUID) ([]domain.ValidationResult, error) {
	return f.created, nil
}
