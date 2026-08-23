package dokploy

import (
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandler = errors.New("invalid Dokploy REST handler configuration")

type Handler struct {
	servers portsin.DokployServerUseCase
	paging  pagination.Config
}

func New(servers portsin.DokployServerUseCase, paging pagination.Config) (*Handler, error) {
	if servers == nil || len(paging.Secret) < 32 || paging.TTL <= 0 {
		return nil, ErrInvalidHandler
	}
	return &Handler{servers: servers, paging: paging}, nil
}
