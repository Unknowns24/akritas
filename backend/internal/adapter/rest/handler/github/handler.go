package github

import (
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandler = errors.New("invalid GitHub REST handler configuration")

type Handler struct {
	accounts portsin.GitHubAccountUseCase
	apps     portsin.GitHubAppUseCase
	paging   pagination.Config
}

func New(accounts portsin.GitHubAccountUseCase, apps portsin.GitHubAppUseCase, paging pagination.Config) (*Handler, error) {
	if accounts == nil || apps == nil || len(paging.Secret) < 32 || paging.TTL <= 0 {
		return nil, ErrInvalidHandler
	}
	return &Handler{accounts: accounts, apps: apps, paging: paging}, nil
}
