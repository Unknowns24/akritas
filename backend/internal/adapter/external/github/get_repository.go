package github

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (c *Client) GetRepository(ctx context.Context, account domain.GitHubAccount, identifier string) (domain.GitHubRepository, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return domain.GitHubRepository{}, domain.ErrIntegrationNotFound
	}
	path := ""
	if _, err := strconv.ParseInt(identifier, 10, 64); err == nil {
		path = "/repositories/" + url.PathEscape(identifier)
	} else {
		parts := strings.Split(identifier, "/")
		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
			return domain.GitHubRepository{}, domain.ErrIntegrationNotFound
		}
		path = "/repos/" + url.PathEscape(strings.TrimSpace(parts[0])) + "/" + url.PathEscape(strings.TrimSpace(parts[1]))
	}
	token, err := c.accountToken(ctx, account)
	if err != nil {
		return domain.GitHubRepository{}, err
	}
	defer wipe(token)
	var payload repositoryDTO
	_, err = c.doJSON(ctx, http.MethodGet, path, string(token), nil, &payload)
	if err != nil {
		var providerErr *providerError
		if errors.As(err, &providerErr) && providerErr.Status == http.StatusNotFound {
			return domain.GitHubRepository{}, domain.ErrIntegrationNotFound
		}
		return domain.GitHubRepository{}, normalizeProviderError(err)
	}
	if !strings.EqualFold(payload.Owner.Login, account.AccountIdentifier) {
		return domain.GitHubRepository{}, domain.ErrIntegrationNotFound
	}
	return domain.NewGitHubRepository(account.ID, strconv.FormatInt(payload.ID, 10), payload.Owner.Login, payload.Name, payload.DefaultBranch, payload.Private, payload.HTMLURL)
}
