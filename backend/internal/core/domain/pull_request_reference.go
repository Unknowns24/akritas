package domain

import (
	"strings"
	"time"
)

type PullRequestReference struct {
	Number     int
	URL        string
	Repository string
	Branch     string
	CreatedAt  time.Time
}

func NewPullRequestReference(number int, pullRequestURL, repository, branch string, createdAt time.Time) (PullRequestReference, error) {
	reference := PullRequestReference{
		Number: number, URL: strings.TrimSpace(pullRequestURL), Repository: strings.TrimSpace(repository),
		Branch: strings.TrimSpace(branch), CreatedAt: createdAt,
	}
	if err := reference.Validate(); err != nil {
		return PullRequestReference{}, err
	}
	return reference, nil
}

func (r PullRequestReference) Validate() error {
	if r.Number < 1 || !validHTTPURL(r.URL) || !nonBlank(r.Repository) || !nonBlank(r.Branch) || !validTime(r.CreatedAt) {
		return ErrInvalidPullRequestReference.Wrap(validationCause("pull request reference"))
	}
	return nil
}
