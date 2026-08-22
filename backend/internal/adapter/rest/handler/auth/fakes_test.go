package auth_test

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var errUnexpected = errors.New("unexpected infrastructure failure")

type fakeGetSetupStatusUseCase struct {
	status in.SetupStatus
	err    error
}

func (f *fakeGetSetupStatusUseCase) Execute(ctx context.Context) (in.SetupStatus, error) {
	return f.status, f.err
}

type fakeStartAdministratorSetupUseCase struct {
	output       in.StartAdministratorSetupOutput
	err          error
	receivedArgs in.StartAdministratorSetupInput
}

func (f *fakeStartAdministratorSetupUseCase) Execute(ctx context.Context, input in.StartAdministratorSetupInput) (in.StartAdministratorSetupOutput, error) {
	f.receivedArgs = input
	return f.output, f.err
}

type fakeVerifyAdministratorSetupUseCase struct {
	output       in.VerifyAdministratorSetupOutput
	err          error
	receivedArgs in.VerifyAdministratorSetupInput
}

func (f *fakeVerifyAdministratorSetupUseCase) Execute(ctx context.Context, input in.VerifyAdministratorSetupInput) (in.VerifyAdministratorSetupOutput, error) {
	f.receivedArgs = input
	return f.output, f.err
}
