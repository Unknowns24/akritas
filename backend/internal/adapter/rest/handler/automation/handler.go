package automation

import (
	"errors"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandler = errors.New("invalid Automation REST handler configuration")

type Handler struct {
	automation portsin.AutomationUseCase
}

func New(automation portsin.AutomationUseCase) (*Handler, error) {
	if automation == nil {
		return nil, ErrInvalidHandler
	}
	return &Handler{automation: automation}, nil
}
