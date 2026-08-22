package domain

import (
	"strings"
	"time"
)

type GitHubIssueReference struct {
	Number     int
	URL        string
	Repository string
	CreatedAt  time.Time
}

func NewGitHubIssueReference(number int, issueURL, repository string, createdAt time.Time) (GitHubIssueReference, error) {
	reference := GitHubIssueReference{
		Number: number, URL: strings.TrimSpace(issueURL), Repository: strings.TrimSpace(repository), CreatedAt: createdAt,
	}
	if err := reference.Validate(); err != nil {
		return GitHubIssueReference{}, err
	}
	return reference, nil
}

func (r GitHubIssueReference) Validate() error {
	if r.Number < 1 || !validHTTPURL(r.URL) || !nonBlank(r.Repository) || !validTime(r.CreatedAt) {
		return ErrInvalidGitHubIssueReference.Wrap(validationCause("GitHub issue reference"))
	}
	return nil
}
