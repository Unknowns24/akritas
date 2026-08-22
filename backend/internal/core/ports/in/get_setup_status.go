package in

import "context"

type SetupStatus struct {
	SetupRequired    bool
	RegistrationOpen bool
}

type GetSetupStatusUseCase interface {
	Execute(ctx context.Context) (SetupStatus, error)
}
