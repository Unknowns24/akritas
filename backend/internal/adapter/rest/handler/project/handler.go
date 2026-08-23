package project

import (
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandler = errors.New("invalid Project REST handler configuration")

type Handler struct {
	projects portsin.ProjectUseCase
	paging   pagination.Config
}

func New(projects portsin.ProjectUseCase, paging pagination.Config) (*Handler, error) {
	if projects == nil || len(paging.Secret) < 32 || paging.TTL <= 0 {
		return nil, ErrInvalidHandler
	}
	return &Handler{projects: projects, paging: paging}, nil
}
