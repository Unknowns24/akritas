package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type CommitDetail struct {
	SHA     string
	Message string
	Author  string
	Date    string
	URL     string
	Files   []CommitFile
}

type CommitFile struct {
	Filename string
	Status   string
	Patch    string
}

func (c *Client) ReadCommit(ctx context.Context, account domain.GitHubAccount, owner, repo, sha string) (CommitDetail, error) {
	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	sha = strings.TrimSpace(sha)
	if owner == "" || repo == "" || sha == "" || strings.Contains(sha, "/") || strings.Contains(sha, "..") {
		return CommitDetail{}, domain.ErrIntegrationUnavailable
	}
	token, err := c.accountToken(ctx, account)
	if err != nil {
		return CommitDetail{}, err
	}
	defer wipe(token)

	path := fmt.Sprintf("/repos/%s/%s/commits/%s", url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(sha))
	var payload struct {
		SHA     string `json:"sha"`
		HTMLURL string `json:"html_url"`
		Commit  struct {
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
				Date string `json:"date"`
			} `json:"author"`
		} `json:"commit"`
		Files []struct {
			Filename string `json:"filename"`
			Status   string `json:"status"`
			Patch    string `json:"patch"`
		} `json:"files"`
	}
	_, err = c.doJSON(ctx, http.MethodGet, path, string(token), nil, &payload)
	if err != nil {
		return CommitDetail{}, normalizeProviderError(err)
	}
	files := make([]CommitFile, 0, len(payload.Files))
	for _, file := range payload.Files {
		files = append(files, CommitFile{Filename: file.Filename, Status: file.Status, Patch: file.Patch})
	}
	return CommitDetail{
		SHA: payload.SHA, Message: payload.Commit.Message, Author: payload.Commit.Author.Name,
		Date: payload.Commit.Author.Date, URL: payload.HTMLURL, Files: files,
	}, nil
}

func (c *Client) ReadDiff(ctx context.Context, account domain.GitHubAccount, owner, repo, sha string) (string, error) {
	detail, err := c.ReadCommit(ctx, account, owner, repo, sha)
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	for _, file := range detail.Files {
		builder.WriteString("file: ")
		builder.WriteString(file.Filename)
		builder.WriteString(" (")
		builder.WriteString(file.Status)
		builder.WriteString(")\n")
		builder.WriteString(file.Patch)
		builder.WriteString("\n\n")
	}
	return builder.String(), nil
}
