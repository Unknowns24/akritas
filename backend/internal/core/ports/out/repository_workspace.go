package out

import (
	"context"
	"time"
)

// CreateBranchInput carries everything CreateBranch needs. BaseBranch is
// received as input rather than resolved internally: resolving the correct
// base/default branch is a GitHub-API concern owned by the caller.
// BranchName is computed by the caller from a Remediation identity it
// controls, so branch creation never depends on IssueReference or a
// GitHub Issue existing.
type CreateBranchInput struct {
	WorkspacePath string
	BaseBranch    string
	BranchName    string
}

// CreateBranchOutput returns enough information for the caller to later
// associate the branch with a Remediation without re-inspecting the
// workspace.
type CreateBranchOutput struct {
	BranchName string
	BaseBranch string
	BaseCommit string
	CreatedAt  time.Time
}

type CommitAllInput struct {
	WorkspacePath string
	BranchName    string
	Message       string
}

type CommitAllOutput struct {
	SHA       string
	Summary   string
	CreatedAt time.Time
}

type PushBranchInput struct {
	WorkspacePath string
	BranchName    string
}

type PushBranchOutput struct {
	BranchName string
	PushedAt   time.Time
}

// RepositoryWorkspace is the minimal, write-capable local-git output port.
// It is the mutation counterpart to the read-only RepositoryInspector (H3),
// and follows the same allowlist discipline: it exposes exactly one
// capability per Git mutation and MUST NOT grow a Run(command string)-shaped
// or otherwise unconstrained passthrough.
type RepositoryWorkspace interface {
	CreateBranch(ctx context.Context, input CreateBranchInput) (CreateBranchOutput, error)
	CommitAll(ctx context.Context, input CommitAllInput) (CommitAllOutput, error)
	PushBranch(ctx context.Context, input PushBranchInput) (PushBranchOutput, error)
}
