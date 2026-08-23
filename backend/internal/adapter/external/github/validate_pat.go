package github

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) ValidatePAT(ctx context.Context, request portsout.GitHubPATValidationRequest) (portsout.GitHubPATValidation, error) {
	if strings.TrimSpace(request.Token) == "" || request.AccountType.Validate() != nil || strings.TrimSpace(request.AccountIdentifier) == "" {
		return portsout.GitHubPATValidation{}, domain.ErrGitHubCredentialRejected
	}
	var identity struct {
		Login string `json:"login"`
		Type  string `json:"type"`
	}
	response, err := c.doJSON(ctx, http.MethodGet, "/user", request.Token, nil, &identity)
	if err != nil {
		return portsout.GitHubPATValidation{}, normalizeCredentialError(err)
	}
	identifier := identity.Login
	if request.AccountType == domain.GitHubAccountOrganization {
		var organization struct {
			Login string `json:"login"`
		}
		_, err = c.doJSON(ctx, http.MethodGet, "/orgs/"+url.PathEscape(request.AccountIdentifier), request.Token, nil, &organization)
		if err != nil {
			return portsout.GitHubPATValidation{}, normalizeCredentialError(err)
		}
		identifier = organization.Login
	}
	if !strings.EqualFold(identifier, request.AccountIdentifier) {
		return portsout.GitHubPATValidation{}, domain.ErrGitHubCredentialRejected
	}
	return portsout.GitHubPATValidation{AccountIdentifier: identifier, ClassicScopes: parseScopes(response.Header.Get("X-OAuth-Scopes"))}, nil
}

func parseScopes(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if scope := strings.TrimSpace(part); scope != "" {
			result = append(result, scope)
		}
	}
	return result
}
