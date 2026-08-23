package out

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type PullRequestContent struct {
	Title string
	Body  string
}

type PullRequestInput struct {
	BaseBranch string
	HeadBranch string
	Content    PullRequestContent
}

type PublishedPullRequest struct {
	Number    int
	URL       string
	CreatedAt time.Time
}

type PullRequestPublisher interface {
	CreateOrFindPullRequest(context.Context, domain.GitHubAccount, domain.GitHubRepository, PullRequestInput) (PublishedPullRequest, error)
}
