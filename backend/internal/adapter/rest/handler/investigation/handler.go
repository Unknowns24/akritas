package investigation

import (
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandler = errors.New("invalid Investigation REST handler configuration")

type Handler struct {
	investigations portsin.InvestigationUseCase
	paging         pagination.Config
}

func New(investigations portsin.InvestigationUseCase, paging pagination.Config) (*Handler, error) {
	if investigations == nil || len(paging.Secret) < 32 || paging.TTL <= 0 {
		return nil, ErrInvalidHandler
	}
	return &Handler{investigations: investigations, paging: paging}, nil
}
