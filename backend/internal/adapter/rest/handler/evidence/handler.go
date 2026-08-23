package evidence

import (
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandler = errors.New("invalid Evidence REST handler configuration")

type Handler struct {
	evidence portsin.EvidenceUseCase
	paging   pagination.Config
}

func New(evidence portsin.EvidenceUseCase, paging pagination.Config) (*Handler, error) {
	if evidence == nil || len(paging.Secret) < 32 || paging.TTL <= 0 {
		return nil, ErrInvalidHandler
	}
	return &Handler{evidence: evidence, paging: paging}, nil
}
