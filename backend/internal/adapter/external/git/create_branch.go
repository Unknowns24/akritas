package git

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) CreateBranch(ctx context.Context, input out.CreateBranchInput) (out.CreateBranchOutput, error) {
	if err := validateRefName(input.BaseBranch); err != nil {
		return out.CreateBranchOutput{}, err
	}
	if err := validateRefName(input.BranchName); err != nil {
		return out.CreateBranchOutput{}, err
	}
	if input.BranchName == input.BaseBranch {
		return out.CreateBranchOutput{}, ErrProtectedBranchTarget
	}

	if _, _, err := c.runGit(ctx, input.WorkspacePath, "rev-parse", "--is-inside-work-tree"); err != nil {
		return out.CreateBranchOutput{}, ErrInvalidWorkspace.Wrap(errNotAGitRepository)
	}

	baseCommit, _, err := c.runGit(ctx, input.WorkspacePath, "rev-parse", "--verify", "refs/heads/"+input.BaseBranch)
	if err != nil {
		if ctx.Err() != nil {
			return out.CreateBranchOutput{}, ErrGitCommandFailed.Wrap(ctx.Err())
		}
		return out.CreateBranchOutput{}, ErrBaseBranchNotFound
	}

	if _, _, err := c.runGit(ctx, input.WorkspacePath, "rev-parse", "--verify", "refs/heads/"+input.BranchName); err == nil {
		return out.CreateBranchOutput{}, ErrBranchAlreadyExists
	}

	if _, _, err := c.runGit(ctx, input.WorkspacePath, "checkout", "-b", input.BranchName, input.BaseBranch); err != nil {
		if ctx.Err() != nil {
			return out.CreateBranchOutput{}, ErrGitCommandFailed.Wrap(ctx.Err())
		}
		return out.CreateBranchOutput{}, ErrGitCommandFailed.Wrap(err)
	}

	return out.CreateBranchOutput{
		BranchName: input.BranchName,
		BaseBranch: input.BaseBranch,
		BaseCommit: baseCommit,
		CreatedAt:  time.Now().UTC(),
	}, nil
}
