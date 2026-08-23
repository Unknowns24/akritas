package out

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type IssueContent struct {
	Title string
	Body  string
}

type PublishedIssue struct {
	Number    int
	URL       string
	CreatedAt time.Time
}

// IssuePublisher is the only H4 GitHub mutation capability. It intentionally
// exposes no branches, contents, commits, pull requests, merge or deploy.
type IssuePublisher interface {
	PublishIssue(context.Context, domain.GitHubAccount, domain.GitHubRepository, IssueContent) (PublishedIssue, error)
}
