package automation

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

var ErrInvalidUseCase = errors.New("invalid automation use case configuration")

type UseCase struct {
	store portsout.AutomationPolicyStore
}

func New(store portsout.AutomationPolicyStore) (portsin.AutomationUseCase, error) {
	if store == nil {
		return nil, ErrInvalidUseCase
	}
	return &UseCase{store: store}, nil
}

func (uc *UseCase) GetPolicy(ctx context.Context) (domain.AutomationPolicy, error) {
	return uc.store.Get(ctx)
}

func (uc *UseCase) PutPolicy(ctx context.Context, policy domain.AutomationPolicy) (domain.AutomationPolicy, error) {
	if err := policy.Validate(); err != nil {
		return domain.AutomationPolicy{}, err
	}
	if err := uc.store.Put(ctx, policy); err != nil {
		return domain.AutomationPolicy{}, err
	}
	return policy, nil
}
