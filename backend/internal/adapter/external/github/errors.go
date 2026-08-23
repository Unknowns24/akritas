package github

import "github.com/Unknowns24/akritas/backend/internal/core/domain"

var ErrInvalidGitHubAppPrivateKey = &domain.Error{
	Code: "0x302001I", Message: "invalid GitHub App private key", UserMessage: "No se pudo autenticar la integración con GitHub.",
}

func Catalog() map[string]*domain.Error {
	return map[string]*domain.Error{"ErrInvalidGitHubAppPrivateKey": ErrInvalidGitHubAppPrivateKey}
}
