package remediation

import (
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandler = errors.New("invalid Remediation REST handler configuration")

type Handler struct {
	remediations  portsin.RemediationUseCase
	paging        pagination.Config
	workspaceRoot string
}

func New(remediations portsin.RemediationUseCase, paging pagination.Config, workspaceRoot string) (*Handler, error) {
	if remediations == nil || len(paging.Secret) < 32 || paging.TTL <= 0 {
		return nil, ErrInvalidHandler
	}
	return &Handler{remediations: remediations, paging: paging, workspaceRoot: workspaceRoot}, nil
}
