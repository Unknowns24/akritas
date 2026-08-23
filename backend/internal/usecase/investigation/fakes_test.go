package investigation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type fakeIncidentStore struct {
	exists    bool
	getResult *domain.Incident
	getErr    error
	updated   []domain.Incident
}

func (f *fakeIncidentStore) Get(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if !f.exists {
		return nil, domain.ErrIncidentNotFound
	}
	if f.getResult == nil {
		return &domain.Incident{ID: id, ProjectID: uuid.New(), Phase: domain.IncidentPhaseDetected}, nil
	}
	return f.getResult, f.getErr
}

func (f *fakeIncidentStore) Lock(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	return f.Get(ctx, id)
}

func (f *fakeIncidentStore) Update(ctx context.Context, value *domain.Incident) error {
	f.updated = append(f.updated, *value)
	return nil
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
	open              []domain.Investigation
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

func (f *fakeInvestigationStore) ListOpen(context.Context) ([]domain.Investigation, error) {
	return append([]domain.Investigation(nil), f.open...), nil
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

func (f *fakeOperationStore) FindByResource(context.Context, domain.OperationResourceType, uuid.UUID) (*domain.Operation, error) {
	if f.findByIDErr != nil {
		return nil, f.findByIDErr
	}
	if f.findByIDResult == nil {
		return nil, domain.ErrOperationNotFound
	}
	return f.findByIDResult, nil
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
	result  out.InvestigationRunResult
	err     error
	calls   int
	context out.InvestigationRunContext
}

func (f *fakeInvestigationRunner) Run(ctx context.Context, runContext out.InvestigationRunContext) (out.InvestigationRunResult, error) {
	f.calls++
	f.context = runContext
	return f.result, f.err
}

type fakeEvidenceStore struct {
	createErr error
	created   []domain.Evidence
}

func (f *fakeEvidenceStore) Create(ctx context.Context, value *domain.Evidence) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = append(f.created, *value)
	return nil
}

func (f *fakeEvidenceStore) ListByInvestigation(ctx context.Context, investigationID uuid.UUID, params paging.Params) (paging.Slice[domain.Evidence], error) {
	return paging.Slice[domain.Evidence]{Items: append([]domain.Evidence(nil), f.created...), Total: int64(len(f.created))}, nil
}

type fakeEvidenceAssembler struct {
	result out.InvestigationRunContext
	err    error
}

func (f *fakeEvidenceAssembler) Assemble(ctx context.Context, investigation domain.Investigation) (out.InvestigationRunContext, error) {
	f.result.Investigation = investigation
	return f.result, f.err
}

type fakeTransactor struct{}

func (fakeTransactor) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type failNthTransactor struct {
	calls  int
	failAt int
	err    error
}

func (f *failNthTransactor) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	f.calls++
	if f.calls == f.failAt {
		return f.err
	}
	return fn(ctx)
}
