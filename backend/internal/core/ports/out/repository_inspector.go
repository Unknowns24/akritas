package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type RepositoryScope struct {
	Account domain.GitHubAccount
	Owner   string
	Name    string
	Branch  string
}

type RepositoryCodeMatch struct {
	Path       string `json:"path"`
	Repository string `json:"repository"`
	URL        string `json:"url"`
}

type RepositoryFile struct {
	Path    string `json:"path"`
	Ref     string `json:"ref"`
	Content string `json:"content"`
	SHA     string `json:"sha"`
}

type RepositoryCommitSummary struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	URL     string `json:"url"`
}

type RepositoryCommitFile struct {
	Filename string `json:"filename"`
	Status   string `json:"status"`
	Patch    string `json:"patch"`
}

type RepositoryCommit struct {
	SHA     string                 `json:"sha"`
	Message string                 `json:"message"`
	Author  string                 `json:"author"`
	Date    string                 `json:"date"`
	URL     string                 `json:"url"`
	Files   []RepositoryCommitFile `json:"files"`
}

// RepositoryInspector is the complete H3 allowlist. It intentionally exposes
// no mutation, shell, branch, issue, commit, push or pull-request capability.
type RepositoryInspector interface {
	SearchCode(context.Context, domain.GitHubAccount, string, string, string) ([]RepositoryCodeMatch, error)
	ReadFile(context.Context, domain.GitHubAccount, string, string, string, string) (RepositoryFile, error)
	ListRecentCommits(context.Context, domain.GitHubAccount, string, string, string, int) ([]RepositoryCommitSummary, error)
	ReadCommit(context.Context, domain.GitHubAccount, string, string, string) (RepositoryCommit, error)
	ReadDiff(context.Context, domain.GitHubAccount, string, string, string) (string, error)
}
