package remediation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

type fakeRemediationStore struct {
	byID       map[uuid.UUID]domain.Remediation
	createErr  error
	getErr     error
	created    []domain.Remediation
	createCalls int
}

func newFakeRemediationStore() *fakeRemediationStore {
	return &fakeRemediationStore{byID: map[uuid.UUID]domain.Remediation{}}
}

func (f *fakeRemediationStore) Create(ctx context.Context, value *domain.Remediation) error {
	f.createCalls++
	if f.createErr != nil {
		return f.createErr
	}
	f.byID[value.ID] = *value
	f.created = append(f.created, *value)
	return nil
}

func (f *fakeRemediationStore) Get(ctx context.Context, id uuid.UUID) (*domain.Remediation, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	value, ok := f.byID[id]
	if !ok {
		return nil, domain.ErrRemediationNotFound
	}
	return &value, nil
}

type createBranchCall struct {
	input portsout.CreateBranchInput
}

type fakeRepositoryWorkspace struct {
	output portsout.CreateBranchOutput
	err    error
	calls  []createBranchCall
}

func (f *fakeRepositoryWorkspace) CreateBranch(ctx context.Context, input portsout.CreateBranchInput) (portsout.CreateBranchOutput, error) {
	f.calls = append(f.calls, createBranchCall{input: input})
	if f.err != nil {
		return portsout.CreateBranchOutput{}, f.err
	}
	// Mimic the real adapter's contract: it echoes back the requested
	// branch/base names. Tests only need to override this when they care
	// about a specific BaseCommit/CreatedAt.
	output := f.output
	if output.BranchName == "" {
		output.BranchName = input.BranchName
	}
	if output.BaseBranch == "" {
		output.BaseBranch = input.BaseBranch
	}
	return output, nil
}

type fakeWorkspaceInspector struct {
	has map[string]bool
	err error
}

func (f *fakeWorkspaceInspector) HasFile(ctx context.Context, workspacePath, relativePath string) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return f.has[relativePath], nil
}

type runCall struct {
	command       portsout.ValidationCommand
	workspacePath string
}

type fakeValidationRunner struct {
	results map[portsout.ValidationCommand]portsout.ExecutionResult
	errs    map[portsout.ValidationCommand]error
	calls   []runCall
}

func newFakeValidationRunner() *fakeValidationRunner {
	return &fakeValidationRunner{
		results: map[portsout.ValidationCommand]portsout.ExecutionResult{},
		errs:    map[portsout.ValidationCommand]error{},
	}
}

func (f *fakeValidationRunner) Run(ctx context.Context, command portsout.ValidationCommand, workspacePath string) (portsout.ExecutionResult, error) {
	f.calls = append(f.calls, runCall{command: command, workspacePath: workspacePath})
	if err, ok := f.errs[command]; ok {
		return portsout.ExecutionResult{}, err
	}
	return f.results[command], nil
}

type fakeValidationResultStore struct {
	createErr error
	created   []domain.ValidationResult
}

func (f *fakeValidationResultStore) Create(ctx context.Context, value *domain.ValidationResult) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = append(f.created, *value)
	return nil
}

func (f *fakeValidationResultStore) ListByRemediation(ctx context.Context, remediationID uuid.UUID) ([]domain.ValidationResult, error) {
	var results []domain.ValidationResult
	for _, value := range f.created {
		if value.RemediationID == remediationID {
			results = append(results, value)
		}
	}
	return results, nil
}
