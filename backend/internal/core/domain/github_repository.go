package domain

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type GitHubRepository struct {
	GitHubAccountID      uuid.UUID `gorm:"type:uuid;not null;index;column:github_account_id"`
	RepositoryIdentifier string    `gorm:"not null;column:repository_identifier"`
	Owner                string    `gorm:"not null;column:repository_owner"`
	Name                 string    `gorm:"not null;column:repository_name"`
	FullName             string    `gorm:"not null;column:repository_full_name"`
	DefaultBranch        string    `gorm:"not null;column:default_branch"`
	Private              bool      `gorm:"not null;column:repository_private"`
	HTMLURL              string    `gorm:"not null;column:repository_html_url"`
}

func NewGitHubRepository(
	githubAccountID uuid.UUID,
	repositoryIdentifier, owner, name, defaultBranch string,
	private bool,
	htmlURL string,
) (GitHubRepository, error) {
	repository := GitHubRepository{
		GitHubAccountID: githubAccountID, RepositoryIdentifier: strings.TrimSpace(repositoryIdentifier),
		Owner: strings.TrimSpace(owner), Name: strings.TrimSpace(name), DefaultBranch: strings.TrimSpace(defaultBranch),
		Private: private, HTMLURL: strings.TrimSpace(htmlURL),
	}
	repository.FullName = fmt.Sprintf("%s/%s", repository.Owner, repository.Name)
	if err := repository.Validate(); err != nil {
		return GitHubRepository{}, err
	}
	return repository, nil
}

func (r GitHubRepository) Validate() error {
	if r.GitHubAccountID == uuid.Nil || !nonBlank(r.RepositoryIdentifier) || !nonBlank(r.Owner) ||
		!nonBlank(r.Name) || !nonBlank(r.DefaultBranch) || r.FullName != r.Owner+"/"+r.Name || !validHTTPURL(r.HTMLURL) {
		return ErrInvalidGitHubRepository.Wrap(validationCause("GitHub repository"))
	}
	return nil
}
