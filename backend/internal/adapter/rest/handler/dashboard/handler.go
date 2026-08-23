package dashboard

import (
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandler = errors.New("invalid Dashboard REST handler configuration")

type Handler struct {
	dashboard portsin.DashboardUseCase
	paging    pagination.Config
}

func New(dashboard portsin.DashboardUseCase, paging pagination.Config) (*Handler, error) {
	if dashboard == nil || len(paging.Secret) < 32 || paging.TTL <= 0 {
		return nil, ErrInvalidHandler
	}
	return &Handler{dashboard: dashboard, paging: paging}, nil
}
