package operation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

type UseCase struct {
	operations portsout.OperationStore
}

func New(operations portsout.OperationStore) portsin.OperationUseCase {
	return &UseCase{operations: operations}
}

func (uc *UseCase) GetOperation(ctx context.Context, id uuid.UUID) (*domain.Operation, error) {
	return uc.operations.FindByID(ctx, id)
}
