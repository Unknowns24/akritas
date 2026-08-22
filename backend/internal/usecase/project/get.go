package project

import (
	"context"

	inproject "github.com/Unknowns24/akritas/backend/internal/core/ports/in/project"
	"github.com/google/uuid"
)

func (uc *UseCase) Get(ctx context.Context, id uuid.UUID) (*inproject.Result, error) {
	project, err := uc.projects.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return inproject.NewResult(project), nil
}
