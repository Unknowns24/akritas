package github

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type createIssueRequestDTO struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type createIssueResponseDTO struct {
	Number    int       `json:"number"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Client) PublishIssue(ctx context.Context, account domain.GitHubAccount, repository domain.GitHubRepository, content portsout.IssueContent) (portsout.PublishedIssue, error) {
	if repository.Validate() != nil || repository.GitHubAccountID != account.ID || content.Title == "" || content.Body == "" {
		return portsout.PublishedIssue{}, domain.ErrInvalidGitHubIssueReference
	}
	token, err := c.accountToken(ctx, account)
	if err != nil {
		return portsout.PublishedIssue{}, err
	}
	defer wipe(token)
	payload, err := json.Marshal(createIssueRequestDTO{Title: content.Title, Body: content.Body})
	if err != nil {
		return portsout.PublishedIssue{}, domain.ErrIntegrationUnavailable
	}
	path := "/repos/" + url.PathEscape(repository.Owner) + "/" + url.PathEscape(repository.Name) + "/issues"
	var response createIssueResponseDTO
	if _, err = c.doJSON(ctx, http.MethodPost, path, string(token), bytes.NewReader(payload), &response); err != nil {
		return portsout.PublishedIssue{}, normalizeProviderError(err)
	}
	if response.Number < 1 || response.HTMLURL == "" || response.CreatedAt.IsZero() {
		return portsout.PublishedIssue{}, domain.ErrIntegrationUnavailable
	}
	return portsout.PublishedIssue{Number: response.Number, URL: response.HTMLURL, CreatedAt: response.CreatedAt.UTC()}, nil
}
