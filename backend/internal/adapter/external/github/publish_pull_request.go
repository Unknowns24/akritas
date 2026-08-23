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

type createPullRequestRequestDTO struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

type pullRequestResponseDTO struct {
	Number    int       `json:"number"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
	Head      struct {
		Ref  string `json:"ref"`
		Repo struct {
			FullName string `json:"full_name"`
		} `json:"repo"`
	} `json:"head"`
	Base struct {
		Ref  string `json:"ref"`
		Repo struct {
			FullName string `json:"full_name"`
		} `json:"repo"`
	} `json:"base"`
}

func (c *Client) CreateOrFindPullRequest(ctx context.Context, account domain.GitHubAccount, repository domain.GitHubRepository, input portsout.PullRequestInput) (portsout.PublishedPullRequest, error) {
	if repository.Validate() != nil || repository.GitHubAccountID != account.ID ||
		input.BaseBranch == "" || input.HeadBranch == "" || input.Content.Title == "" || input.Content.Body == "" ||
		input.BaseBranch == input.HeadBranch {
		return portsout.PublishedPullRequest{}, domain.ErrInvalidPullRequestReference
	}
	token, err := c.accountToken(ctx, account)
	if err != nil {
		return portsout.PublishedPullRequest{}, err
	}
	defer wipe(token)

	if existing, found, err := c.findPullRequest(ctx, repository, string(token), input.BaseBranch, input.HeadBranch); err != nil {
		return portsout.PublishedPullRequest{}, err
	} else if found {
		return existing, nil
	}

	payload, err := json.Marshal(createPullRequestRequestDTO{
		Title: input.Content.Title, Body: input.Content.Body,
		Head: input.HeadBranch, Base: input.BaseBranch,
	})
	if err != nil {
		return portsout.PublishedPullRequest{}, domain.ErrIntegrationUnavailable
	}
	path := "/repos/" + url.PathEscape(repository.Owner) + "/" + url.PathEscape(repository.Name) + "/pulls"
	var response pullRequestResponseDTO
	if _, err = c.doJSON(ctx, http.MethodPost, path, string(token), bytes.NewReader(payload), &response); err != nil {
		if reconciled, found, reconcileErr := c.findPullRequest(ctx, repository, string(token), input.BaseBranch, input.HeadBranch); reconcileErr == nil && found {
			return reconciled, nil
		}
		return portsout.PublishedPullRequest{}, normalizeProviderError(err)
	}
	return mapPullRequestResponse(repository, input.BaseBranch, input.HeadBranch, response)
}

func (c *Client) findPullRequest(ctx context.Context, repository domain.GitHubRepository, token, baseBranch, headBranch string) (portsout.PublishedPullRequest, bool, error) {
	query := url.Values{}
	query.Set("state", "open")
	query.Set("base", baseBranch)
	query.Set("head", repository.Owner+":"+headBranch)
	path := "/repos/" + url.PathEscape(repository.Owner) + "/" + url.PathEscape(repository.Name) + "/pulls?" + query.Encode()
	var response []pullRequestResponseDTO
	if _, err := c.doJSON(ctx, http.MethodGet, path, token, nil, &response); err != nil {
		return portsout.PublishedPullRequest{}, false, normalizeProviderError(err)
	}
	for _, candidate := range response {
		mapped, err := mapPullRequestResponse(repository, baseBranch, headBranch, candidate)
		if err == nil {
			return mapped, true, nil
		}
	}
	return portsout.PublishedPullRequest{}, false, nil
}

func mapPullRequestResponse(repository domain.GitHubRepository, baseBranch, headBranch string, response pullRequestResponseDTO) (portsout.PublishedPullRequest, error) {
	if response.Number < 1 || response.HTMLURL == "" || response.CreatedAt.IsZero() ||
		response.Base.Ref != baseBranch || response.Head.Ref != headBranch ||
		response.Base.Repo.FullName != repository.FullName || response.Head.Repo.FullName != repository.FullName {
		return portsout.PublishedPullRequest{}, domain.ErrIntegrationUnavailable
	}
	return portsout.PublishedPullRequest{Number: response.Number, URL: response.HTMLURL, CreatedAt: response.CreatedAt.UTC()}, nil
}
