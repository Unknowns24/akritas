package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type CodeSearchMatch = portsout.RepositoryCodeMatch

func (c *Client) SearchCode(ctx context.Context, account domain.GitHubAccount, owner, repo, query string) ([]CodeSearchMatch, error) {
	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	query = strings.TrimSpace(query)
	if owner == "" || repo == "" || query == "" {
		return nil, domain.ErrIntegrationUnavailable
	}
	token, err := c.accountToken(ctx, account)
	if err != nil {
		return nil, err
	}
	defer wipe(token)

	q := fmt.Sprintf("%s repo:%s/%s", query, owner, repo)
	path := "/search/code?" + url.Values{"q": {q}, "per_page": {"10"}}.Encode()
	var payload struct {
		Items []struct {
			Path       string `json:"path"`
			HTMLURL    string `json:"html_url"`
			Repository struct {
				FullName string `json:"full_name"`
			} `json:"repository"`
		} `json:"items"`
	}
	_, err = c.doJSON(ctx, http.MethodGet, path, string(token), nil, &payload)
	if err != nil {
		return nil, normalizeProviderError(err)
	}
	matches := make([]CodeSearchMatch, 0, len(payload.Items))
	for _, item := range payload.Items {
		if !strings.EqualFold(item.Repository.FullName, owner+"/"+repo) {
			continue
		}
		matches = append(matches, CodeSearchMatch{Path: item.Path, Repository: item.Repository.FullName, URL: item.HTMLURL})
	}
	return matches, nil
}
