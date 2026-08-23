package monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) NextIncidentKey(ctx context.Context) (string, error) {
	var next int64
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Raw("SELECT nextval('incident_key_sequence')").Scan(&next).Error; err != nil {
		return "", mapError(err)
	}
	return fmt.Sprintf("AKR-%d", next), nil
}

func (r *Repository) FindGroupableIncident(ctx context.Context, projectID uuid.UUID, fingerprint string, occurredAt time.Time, window time.Duration) (*domain.Incident, error) {
	var incident domain.Incident
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("incidents").Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("project_id = ? AND fingerprint = ? AND phase NOT IN ? AND last_seen_at <= ? AND last_seen_at >= ?", projectID, fingerprint, []domain.IncidentPhase{domain.IncidentPhaseCompleted, domain.IncidentPhaseFailed}, occurredAt, occurredAt.Add(-window)).
		Order("last_seen_at DESC, id DESC").Take(&incident).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, mapError(err)
	}
	return &incident, nil
}

func (r *Repository) CreateIncident(ctx context.Context, incident *domain.Incident) error {
	if incident == nil || incident.Validate() != nil {
		return domain.ErrInvalidIncident
	}
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("incidents").Create(incident).Error; err != nil {
		return mapError(err)
	}
	return nil
}

func (r *Repository) UpdateIncident(ctx context.Context, incident *domain.Incident) error {
	result := txcontext.From(ctx, r.db).WithContext(ctx).Table("incidents").Where("id = ?", incident.ID).Updates(map[string]any{
		"severity": incident.Severity, "last_seen_at": incident.LastSeenAt, "occurrence_count": incident.OccurrenceCount,
	})
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected != 1 {
		return domain.ErrIncidentNotFound
	}
	return nil
}
