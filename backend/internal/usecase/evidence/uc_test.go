package evidence

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type fakeInvestigationStore struct {
	findByIDResult *domain.Investigation
	findByIDErr    error
}

func (f *fakeInvestigationStore) Create(ctx context.Context, value *domain.Investigation) error {
	return nil
}
func (f *fakeInvestigationStore) Update(ctx context.Context, value *domain.Investigation) error {
	return nil
}
func (f *fakeInvestigationStore) FindByID(ctx context.Context, id uuid.UUID) (*domain.Investigation, error) {
	if f.findByIDErr != nil {
		return nil, f.findByIDErr
	}
	return f.findByIDResult, nil
}
func (f *fakeInvestigationStore) ListByIncident(ctx context.Context, incidentID uuid.UUID, params paging.Params) (paging.Slice[domain.Investigation], error) {
	return paging.Slice[domain.Investigation]{}, nil
}
func (f *fakeInvestigationStore) ExistsActiveForIncident(ctx context.Context, incidentID uuid.UUID) (bool, error) {
	return false, nil
}

type fakeEvidenceStore struct {
	listResult   paging.Slice[domain.Evidence]
	listErr      error
	listCalled   bool
	listInvestID uuid.UUID
}

func (f *fakeEvidenceStore) Create(ctx context.Context, value *domain.Evidence) error { return nil }
func (f *fakeEvidenceStore) ListByInvestigation(ctx context.Context, investigationID uuid.UUID, params paging.Params) (paging.Slice[domain.Evidence], error) {
	f.listCalled = true
	f.listInvestID = investigationID
	return f.listResult, f.listErr
}

func TestListInvestigationEvidenceRejectsMissingInvestigation(t *testing.T) {
	t.Parallel()
	investigations := &fakeInvestigationStore{findByIDErr: domain.ErrInvestigationNotFound}
	store := &fakeEvidenceStore{}
	uc := New(investigations, store)

	_, err := uc.ListInvestigationEvidence(context.Background(), uuid.New(), paging.Params{})
	if !errors.Is(err, domain.ErrInvestigationNotFound) {
		t.Fatalf("expected ErrInvestigationNotFound, got %v", err)
	}
	if store.listCalled {
		t.Fatal("evidence must not be listed when the investigation does not exist")
	}
}

func TestListInvestigationEvidenceReturnsStorePage(t *testing.T) {
	t.Parallel()
	investigationID := uuid.New()
	investigations := &fakeInvestigationStore{findByIDResult: &domain.Investigation{ID: investigationID, CreatedAt: time.Now()}}
	store := &fakeEvidenceStore{listResult: paging.Slice[domain.Evidence]{Items: []domain.Evidence{{ID: uuid.New()}}, Total: 1}}
	uc := New(investigations, store)

	page, err := uc.ListInvestigationEvidence(context.Background(), investigationID, paging.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Items) != 1 {
		t.Fatal("expected the store's page to be returned unchanged")
	}
	if !store.listCalled || store.listInvestID != investigationID {
		t.Fatal("expected the list to be scoped to the given investigation")
	}
}
