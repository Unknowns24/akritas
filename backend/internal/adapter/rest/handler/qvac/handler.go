package qvac

import (
	"errors"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandler = errors.New("invalid QVAC REST handler configuration")

type Handler struct {
	qvac portsin.QvacUseCase
}

func New(qvac portsin.QvacUseCase) (*Handler, error) {
	if qvac == nil {
		return nil, ErrInvalidHandler
	}
	return &Handler{qvac: qvac}, nil
}
