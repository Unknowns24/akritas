package incident

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (uc *UseCase) Get(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	if id == uuid.Nil {
		return nil, domain.ErrIncidentNotFound
	}
	incident, err := uc.incidents.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if uc.investigations != nil {
		latest, latestErr := uc.investigations.FindLatestByIncident(ctx, id)
		if latestErr != nil && !errors.Is(latestErr, domain.ErrInvestigationNotFound) {
			return nil, latestErr
		}
		incident.LatestInvestigation = latest
	}
	if uc.issueRefs != nil {
		reference, referenceErr := uc.issueRefs.FindLatestByIncident(ctx, id)
		if referenceErr != nil {
			return nil, referenceErr
		}
		incident.GitHubIssueReference = reference
	}
	return incident, nil
}
