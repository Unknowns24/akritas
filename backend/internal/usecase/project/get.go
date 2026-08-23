package project

import (
	"context"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/google/uuid"
)

func (uc *UseCase) Get(ctx context.Context, id uuid.UUID) (*portsin.ProjectResult, error) {
	project, err := uc.projects.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return result(project), nil
}
