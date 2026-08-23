package investigation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type fakeIncidentReader struct {
	exists bool
	err    error
}

func (f *fakeIncidentReader) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return f.exists, f.err
}

type fakeInvestigationStore struct {
	createErr         error
	updateErr         error
	created           *domain.Investigation
	updated           []domain.Investigation
	findByIDResult    *domain.Investigation
	findByIDErr       error
	listResult        paging.Slice[domain.Investigation]
	listErr           error
	activeForIncident bool
	activeErr         error
}

func (f *fakeInvestigationStore) Create(ctx context.Context, value *domain.Investigation) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = value
	return nil
}

func (f *fakeInvestigationStore) Update(ctx context.Context, value *domain.Investigation) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	f.updated = append(f.updated, *value)
	return nil
}

func (f *fakeInvestigationStore) FindByID(ctx context.Context, id uuid.UUID) (*domain.Investigation, error) {
	if f.findByIDErr != nil {
		return nil, f.findByIDErr
	}
	return f.findByIDResult, nil
}

func (f *fakeInvestigationStore) ListByIncident(ctx context.Context, incidentID uuid.UUID, params paging.Params) (paging.Slice[domain.Investigation], error) {
	return f.listResult, f.listErr
}

func (f *fakeInvestigationStore) ExistsActiveForIncident(ctx context.Context, incidentID uuid.UUID) (bool, error) {
	return f.activeForIncident, f.activeErr
}

type fakeOperationStore struct {
	createErr       error
	updateErr       error
	created         *domain.Operation
	updated         []domain.Operation
	findByIDResult  *domain.Operation
	findByIDErr     error
	findByKeyResult *domain.Operation
	findByKeyErr    error
	findByKeyCalled bool
}

func (f *fakeOperationStore) Create(ctx context.Context, value *domain.Operation) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = value
	return nil
}

func (f *fakeOperationStore) Update(ctx context.Context, value *domain.Operation) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	f.updated = append(f.updated, *value)
	return nil
}

func (f *fakeOperationStore) FindByID(ctx context.Context, id uuid.UUID) (*domain.Operation, error) {
	if f.findByIDErr != nil {
		return nil, f.findByIDErr
	}
	return f.findByIDResult, nil
}

func (f *fakeOperationStore) FindByIdempotencyKey(ctx context.Context, key string) (*domain.Operation, error) {
	f.findByKeyCalled = true
	if f.findByKeyErr != nil {
		return nil, f.findByKeyErr
	}
	return f.findByKeyResult, nil
}

type fakeInvestigationDispatcher struct {
	dispatched                   bool
	investigationID, operationID uuid.UUID
}

func (f *fakeInvestigationDispatcher) Dispatch(investigationID, operationID uuid.UUID) {
	f.dispatched = true
	f.investigationID = investigationID
	f.operationID = operationID
}

type fakeInvestigationRunner struct {
	result out.InvestigationRunResult
	err    error
}

func (f *fakeInvestigationRunner) Run(ctx context.Context, investigation domain.Investigation) (out.InvestigationRunResult, error) {
	return f.result, f.err
}
