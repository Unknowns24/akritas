package remediation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/service/validationpolicy"
	"github.com/google/uuid"
)

func fixedNow() func() time.Time {
	t := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	return func() time.Time { return t }
}

func newTestUseCase(remediations *fakeRemediationStore, workspace *fakeRepositoryWorkspace, runner *fakeValidationRunner, results portsout.ValidationResultStore, inspector *fakeWorkspaceInspector) portsin.RemediationUseCase {
	policy := validationpolicy.New(inspector)
	return New(workspace, runner, remediations, results, policy, uuid.New, fixedNow())
}

func TestCreateRemediationBranchNewID(t *testing.T) {
	remediations := newFakeRemediationStore()
	ws := &fakeRepositoryWorkspace{}
	ws.output.BaseCommit = "deadbeef"
	ws.output.CreatedAt = time.Now()

	uc := newTestUseCase(remediations, ws, newFakeValidationRunner(), &fakeValidationResultStore{}, &fakeWorkspaceInspector{})

	remediationID := uuid.New()
	incidentID := uuid.New()
	got, err := uc.CreateRemediationBranch(context.Background(), portsin.CreateRemediationBranchCommand{
		RemediationID: remediationID, IncidentID: incidentID, WorkspacePath: "/workspace", BaseBranch: "main",
	})
	if err != nil {
		t.Fatalf("CreateRemediationBranch: %v", err)
	}
	if got.Status != domain.RemediationStatusInProgress {
		t.Fatalf("expected in_progress, got %v", got.Status)
	}
	wantBranch := remediationBranchName(remediationID)
	if got.BranchName != wantBranch {
		t.Fatalf("expected branch %q, got %q", wantBranch, got.BranchName)
	}
	if len(ws.calls) != 1 {
		t.Fatalf("expected exactly 1 CreateBranch call, got %d", len(ws.calls))
	}
	if ws.calls[0].input.BranchName != wantBranch || ws.calls[0].input.BaseBranch != "main" || ws.calls[0].input.WorkspacePath != "/workspace" {
		t.Fatalf("unexpected CreateBranch input: %+v", ws.calls[0].input)
	}
	if remediations.createCalls != 1 {
		t.Fatalf("expected exactly 1 Create call, got %d", remediations.createCalls)
	}
}

func TestCreateRemediationBranchIdempotentReplay(t *testing.T) {
	remediations := newFakeRemediationStore()
	ws := &fakeRepositoryWorkspace{}
	uc := newTestUseCase(remediations, ws, newFakeValidationRunner(), &fakeValidationResultStore{}, &fakeWorkspaceInspector{})

	remediationID := uuid.New()
	cmd := portsin.CreateRemediationBranchCommand{RemediationID: remediationID, IncidentID: uuid.New(), WorkspacePath: "/workspace", BaseBranch: "main"}

	first, err := uc.CreateRemediationBranch(context.Background(), cmd)
	if err != nil {
		t.Fatalf("first call: %v", err)
	}
	second, err := uc.CreateRemediationBranch(context.Background(), cmd)
	if err != nil {
		t.Fatalf("second call: %v", err)
	}
	if second.ID != first.ID || second.BranchName != first.BranchName {
		t.Fatalf("expected idempotent replay to return the same remediation, got %+v vs %+v", first, second)
	}
	if len(ws.calls) != 1 {
		t.Fatalf("expected CreateBranch to be called exactly once across both calls, got %d", len(ws.calls))
	}
	if remediations.createCalls != 1 {
		t.Fatalf("expected Create to be called exactly once across both calls, got %d", remediations.createCalls)
	}
}

func TestCreateRemediationBranchWorkspaceFailureLeavesNoPersistedRow(t *testing.T) {
	remediations := newFakeRemediationStore()
	ws := &fakeRepositoryWorkspace{err: errors.New("git failure")}
	uc := newTestUseCase(remediations, ws, newFakeValidationRunner(), &fakeValidationResultStore{}, &fakeWorkspaceInspector{})

	_, err := uc.CreateRemediationBranch(context.Background(), portsin.CreateRemediationBranchCommand{
		RemediationID: uuid.New(), IncidentID: uuid.New(), WorkspacePath: "/workspace", BaseBranch: "main",
	})
	if err == nil {
		t.Fatal("expected an error")
	}
	if remediations.createCalls != 0 {
		t.Fatalf("expected no Create call after a workspace failure, got %d", remediations.createCalls)
	}
}

func TestCreateRemediationBranchStoreFailureAfterSuccessfulBranch(t *testing.T) {
	remediations := newFakeRemediationStore()
	remediations.createErr = errors.New("db down")
	ws := &fakeRepositoryWorkspace{}
	uc := newTestUseCase(remediations, ws, newFakeValidationRunner(), &fakeValidationResultStore{}, &fakeWorkspaceInspector{})

	_, err := uc.CreateRemediationBranch(context.Background(), portsin.CreateRemediationBranchCommand{
		RemediationID: uuid.New(), IncidentID: uuid.New(), WorkspacePath: "/workspace", BaseBranch: "main",
	})
	if err == nil {
		t.Fatal("expected an error")
	}
	if len(ws.calls) != 1 {
		t.Fatalf("expected the branch to have been created before the store failure, got %d calls", len(ws.calls))
	}
}

func TestRemediationBranchNameIsDeterministic(t *testing.T) {
	id := uuid.New()
	if remediationBranchName(id) != remediationBranchName(id) {
		t.Fatal("expected remediationBranchName to be deterministic for the same ID")
	}
	other := uuid.New()
	if remediationBranchName(id) == remediationBranchName(other) {
		t.Fatal("expected different IDs to produce different branch names")
	}
}
