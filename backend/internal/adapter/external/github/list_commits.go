package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type CommitSummary struct {
	SHA     string
	Message string
	Author  string
	Date    string
	URL     string
}

func (c *Client) ListRecentCommits(ctx context.Context, account domain.GitHubAccount, owner, repo, branch string, limit int) ([]CommitSummary, error) {
	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	branch = strings.TrimSpace(branch)
	if owner == "" || repo == "" {
		return nil, domain.ErrIntegrationUnavailable
	}
	if limit <= 0 || limit > 30 {
		limit = 10
	}
	token, err := c.accountToken(ctx, account)
	if err != nil {
		return nil, err
	}
	defer wipe(token)

	values := url.Values{"per_page": {strconv.Itoa(limit)}}
	if branch != "" {
		values.Set("sha", branch)
	}
	path := fmt.Sprintf("/repos/%s/%s/commits?%s", url.PathEscape(owner), url.PathEscape(repo), values.Encode())
	var payload []struct {
		SHA     string `json:"sha"`
		HTMLURL string `json:"html_url"`
		Commit  struct {
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
				Date string `json:"date"`
			} `json:"author"`
		} `json:"commit"`
	}
	_, err = c.doJSON(ctx, http.MethodGet, path, string(token), nil, &payload)
	if err != nil {
		return nil, normalizeProviderError(err)
	}
	out := make([]CommitSummary, 0, len(payload))
	for _, item := range payload {
		out = append(out, CommitSummary{
			SHA: item.SHA, Message: item.Commit.Message, Author: item.Commit.Author.Name,
			Date: item.Commit.Author.Date, URL: item.HTMLURL,
		})
	}
	return out, nil
}
