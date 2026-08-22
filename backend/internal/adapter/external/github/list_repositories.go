package github

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
)

func (c *Client) ListRepositories(ctx context.Context, account domain.GitHubAccount, params paging.Params) (paging.Slice[domain.GitHubRepository], error) {
	if params.Limit < 1 {
		params.Limit = 25
	}
	if params.Limit > 100 {
		return paging.Slice[domain.GitHubRepository]{}, domain.ErrInvalidGitHubRepository
	}
	pageNumber := 1
	if position := providerBoundary(params, "provider_page"); position != "" {
		parsed, err := strconv.Atoi(position)
		if err != nil || parsed < 1 {
			return paging.Slice[domain.GitHubRepository]{}, domain.ErrInvalidGitHubRepository
		}
		pageNumber = parsed
	}
	token, err := c.accountToken(ctx, account)
	if err != nil {
		return paging.Slice[domain.GitHubRepository]{}, err
	}
	defer wipe(token)
	path := "/user/repos"
	appAuthentication := account.AuthenticationMethod == domain.GitHubAuthenticationGitHubApp
	if appAuthentication {
		path = "/installation/repositories"
	} else if account.AccountType == domain.GitHubAccountOrganization {
		path = "/orgs/" + url.PathEscape(account.AccountIdentifier) + "/repos"
	}
	values := url.Values{"per_page": {strconv.Itoa(params.Limit)}, "page": {strconv.Itoa(pageNumber)}}
	if account.AccountType == domain.GitHubAccountPersonal {
		values.Set("affiliation", "owner,collaborator,organization_member")
	}
	var repositories []repositoryDTO
	if appAuthentication {
		var response struct {
			TotalCount   int64           `json:"total_count"`
			Repositories []repositoryDTO `json:"repositories"`
		}
		_, err = c.doJSON(ctx, http.MethodGet, path+"?"+values.Encode(), string(token), nil, &response)
		repositories = response.Repositories
	} else {
		_, err = c.doJSON(ctx, http.MethodGet, path+"?"+values.Encode(), string(token), nil, &repositories)
	}
	if err != nil {
		return paging.Slice[domain.GitHubRepository]{}, normalizeProviderError(err)
	}
	items := make([]domain.GitHubRepository, 0, len(repositories))
	for _, repository := range repositories {
		if !strings.EqualFold(repository.Owner.Login, account.AccountIdentifier) {
			continue
		}
		nameLike := params.Filters["name_like"]
		if nameLike != "" && !strings.Contains(strings.ToLower(repository.Name), strings.ToLower(nameLike)) {
			continue
		}
		mapped, mapErr := domain.NewGitHubRepository(account.ID, strconv.FormatInt(repository.ID, 10), repository.Owner.Login, repository.Name, repository.DefaultBranch, repository.Private, repository.HTMLURL)
		if mapErr == nil {
			items = append(items, mapped)
		}
	}
	result := paging.Slice[domain.GitHubRepository]{Items: items, Total: int64(len(items))}
	if len(repositories) == params.Limit {
		result.NextBoundary = map[string]string{"provider_page": strconv.Itoa(pageNumber + 1)}
	}
	if pageNumber > 1 {
		result.PrevBoundary = map[string]string{"provider_page": strconv.Itoa(pageNumber - 1)}
	}
	return result, nil
}

func providerBoundary(params paging.Params, field string) string {
	if params.Cursor == nil {
		return ""
	}
	if value := params.Cursor.After[field]; value != "" {
		return value
	}
	return params.Cursor.Before[field]
}

type repositoryDTO struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	DefaultBranch string `json:"default_branch"`
	Private       bool   `json:"private"`
	HTMLURL       string `json:"html_url"`
	Owner         struct {
		Login string `json:"login"`
	} `json:"owner"`
}
