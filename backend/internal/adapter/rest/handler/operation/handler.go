package operation

import (
	"errors"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandler = errors.New("invalid Operation REST handler configuration")

type Handler struct {
	operations portsin.OperationUseCase
}

func New(operations portsin.OperationUseCase) (*Handler, error) {
	if operations == nil {
		return nil, ErrInvalidHandler
	}
	return &Handler{operations: operations}, nil
}
