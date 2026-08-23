package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type FileContent = portsout.RepositoryFile

func (c *Client) ReadFile(ctx context.Context, account domain.GitHubAccount, owner, repo, filePath, ref string) (FileContent, error) {
	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	ref = strings.TrimSpace(ref)
	cleanPath, err := sanitizeRepositoryPath(filePath)
	if err != nil {
		return FileContent{}, err
	}
	if owner == "" || repo == "" {
		return FileContent{}, domain.ErrIntegrationUnavailable
	}
	token, err := c.accountToken(ctx, account)
	if err != nil {
		return FileContent{}, err
	}
	defer wipe(token)

	endpoint := fmt.Sprintf("/repos/%s/%s/contents/%s", url.PathEscape(owner), url.PathEscape(repo), escapePathSegments(cleanPath))
	if ref != "" {
		endpoint += "?" + url.Values{"ref": {ref}}.Encode()
	}
	var payload struct {
		Type     string `json:"type"`
		Path     string `json:"path"`
		SHA      string `json:"sha"`
		Encoding string `json:"encoding"`
		Content  string `json:"content"`
	}
	_, err = c.doJSON(ctx, http.MethodGet, endpoint, string(token), nil, &payload)
	if err != nil {
		return FileContent{}, normalizeProviderError(err)
	}
	if payload.Type != "file" {
		return FileContent{}, domain.ErrIntegrationUnavailable.Wrap(errString("path is not a file"))
	}
	content := payload.Content
	if strings.EqualFold(payload.Encoding, "base64") {
		decoded, decodeErr := base64.StdEncoding.DecodeString(strings.ReplaceAll(payload.Content, "\n", ""))
		if decodeErr != nil {
			return FileContent{}, domain.ErrIntegrationUnavailable.Wrap(decodeErr)
		}
		content = string(decoded)
	}
	return FileContent{Path: payload.Path, Ref: ref, Content: content, SHA: payload.SHA}, nil
}

func escapePathSegments(value string) string {
	parts := strings.Split(value, "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}
