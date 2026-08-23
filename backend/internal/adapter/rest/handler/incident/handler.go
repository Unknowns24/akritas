package incident

import (
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandler = errors.New("invalid Incident REST handler configuration")

type Handler struct {
	incidents portsin.IncidentUseCase
	paging    pagination.Config
}

func New(incidents portsin.IncidentUseCase, paging pagination.Config) (*Handler, error) {
	if incidents == nil || len(paging.Secret) < 32 || paging.TTL <= 0 {
		return nil, ErrInvalidHandler
	}
	return &Handler{incidents: incidents, paging: paging}, nil
}
