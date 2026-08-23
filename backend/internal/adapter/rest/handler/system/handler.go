package system

import (
	"errors"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandler = errors.New("invalid System REST handler configuration")

type Handler struct {
	system portsin.SystemUseCase
}

func New(system portsin.SystemUseCase) (*Handler, error) {
	if system == nil {
		return nil, ErrInvalidHandler
	}
	return &Handler{system: system}, nil
}
