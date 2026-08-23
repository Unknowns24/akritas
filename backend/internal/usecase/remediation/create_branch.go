package remediation

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// CreateRemediationBranch is idempotent by RemediationID: a replay with the
// same ID returns the already-persisted Remediation without touching the
// git workspace again, per the background-processes policy's "do not
// create multiple remediation branches for the same successful workflow"
// rule.
func (uc *UseCase) CreateRemediationBranch(ctx context.Context, cmd portsin.CreateRemediationBranchCommand) (*domain.Remediation, error) {
	existing, err := uc.remediations.Get(ctx, cmd.RemediationID)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, domain.ErrRemediationNotFound) {
		return nil, err
	}

	remediation, err := domain.NewRemediation(cmd.RemediationID, cmd.IncidentID, uc.now())
	if err != nil {
		return nil, err
	}

	branchName := remediationBranchName(cmd.RemediationID)
	output, err := uc.workspace.CreateBranch(ctx, portsout.CreateBranchInput{
		WorkspacePath: cmd.WorkspacePath, BaseBranch: cmd.BaseBranch, BranchName: branchName,
	})
	if err != nil {
		return nil, err
	}

	if err := remediation.Start(output.BranchName, uc.now()); err != nil {
		return nil, err
	}
	if err := uc.remediations.Create(ctx, remediation); err != nil {
		return nil, err
	}
	return remediation, nil
}
