package git

import (
	"context"
	"strings"
	"time"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) CommitAll(ctx context.Context, input portsout.CommitAllInput) (portsout.CommitAllOutput, error) {
	if err := validateRefName(input.BranchName); err != nil {
		return portsout.CommitAllOutput{}, err
	}
	message := strings.TrimSpace(input.Message)
	if message == "" || len(message) > 1000 {
		return portsout.CommitAllOutput{}, ErrGitCommandFailed
	}
	currentBranch, _, err := c.runGit(ctx, input.WorkspacePath, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return portsout.CommitAllOutput{}, ErrInvalidWorkspace.Wrap(errNotAGitRepository)
	}
	if currentBranch != input.BranchName {
		return portsout.CommitAllOutput{}, ErrProtectedBranchTarget
	}
	if _, _, err = c.runGit(ctx, input.WorkspacePath, "add", "-A"); err != nil {
		return portsout.CommitAllOutput{}, ErrGitCommandFailed.Wrap(err)
	}
	summary, _, err := c.runGit(ctx, input.WorkspacePath, "diff", "--cached", "--name-status")
	if err != nil {
		return portsout.CommitAllOutput{}, ErrGitCommandFailed.Wrap(err)
	}
	if strings.TrimSpace(summary) == "" {
		return portsout.CommitAllOutput{}, ErrGitCommandFailed
	}
	if _, _, err = c.runGit(ctx, input.WorkspacePath, "commit", "-m", message); err != nil {
		if ctx.Err() != nil {
			return portsout.CommitAllOutput{}, ErrGitCommandFailed.Wrap(ctx.Err())
		}
		return portsout.CommitAllOutput{}, ErrGitCommandFailed.Wrap(err)
	}
	sha, _, err := c.runGit(ctx, input.WorkspacePath, "rev-parse", "HEAD")
	if err != nil || strings.TrimSpace(sha) == "" {
		return portsout.CommitAllOutput{}, ErrGitCommandFailed.Wrap(err)
	}
	return portsout.CommitAllOutput{
		SHA:       strings.TrimSpace(sha),
		Summary:   summary,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (c *Client) PushBranch(ctx context.Context, input portsout.PushBranchInput) (portsout.PushBranchOutput, error) {
	if err := validateRefName(input.BranchName); err != nil {
		return portsout.PushBranchOutput{}, err
	}
	if _, _, err := c.runGit(ctx, input.WorkspacePath, "rev-parse", "--verify", "refs/heads/"+input.BranchName); err != nil {
		return portsout.PushBranchOutput{}, ErrBaseBranchNotFound.Wrap(err)
	}
	refspec := "refs/heads/" + input.BranchName + ":refs/heads/" + input.BranchName
	if _, _, err := c.runGit(ctx, input.WorkspacePath, "push", "origin", refspec); err != nil {
		if ctx.Err() != nil {
			return portsout.PushBranchOutput{}, ErrGitCommandFailed.Wrap(ctx.Err())
		}
		return portsout.PushBranchOutput{}, ErrGitCommandFailed.Wrap(err)
	}
	return portsout.PushBranchOutput{BranchName: input.BranchName, PushedAt: time.Now().UTC()}, nil
}
