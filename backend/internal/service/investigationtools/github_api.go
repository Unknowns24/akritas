package investigationtools

import (
	"context"

	githubexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/github"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// GitHubAPI is the read-only GitHub surface used by investigation tools.
type GitHubAPI interface {
	SearchCode(ctx context.Context, account domain.GitHubAccount, owner, repo, query string) ([]githubexternal.CodeSearchMatch, error)
	ReadFile(ctx context.Context, account domain.GitHubAccount, owner, repo, path, ref string) (githubexternal.FileContent, error)
	ListRecentCommits(ctx context.Context, account domain.GitHubAccount, owner, repo, branch string, limit int) ([]githubexternal.CommitSummary, error)
	ReadCommit(ctx context.Context, account domain.GitHubAccount, owner, repo, sha string) (githubexternal.CommitDetail, error)
	ReadDiff(ctx context.Context, account domain.GitHubAccount, owner, repo, sha string) (string, error)
}

// Ensure the production GitHub adapter satisfies GitHubAPI.
var _ GitHubAPI = (*githubexternal.Client)(nil)
