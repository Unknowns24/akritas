package monitoring

import (
	"context"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) OccurrenceExists(ctx context.Context, projectID uuid.UUID, occurrenceKey string) (bool, error) {
	var count int64
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("log_events").Where("project_id = ? AND occurrence_key = ?", projectID, occurrenceKey).Count(&count).Error
	return count > 0, mapErrorOrNil(err)
}

func (r *Repository) CreateLogEvent(ctx context.Context, event *domain.LogEvent) error {
	if event == nil || event.IncidentID == uuid.Nil || event.SourceType.Validate() != nil || strings.TrimSpace(event.SourceResourceID) == "" || strings.TrimSpace(event.SourceInstanceID) == "" || strings.TrimSpace(event.OccurrenceKey) == "" || event.Validate() != nil {
		return domain.ErrInvalidLogEvent
	}
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("log_events").Create(event).Error; err != nil {
		return mapError(err)
	}
	return nil
}

func mapErrorOrNil(err error) error {
	if err == nil {
		return nil
	}
	return mapError(err)
}
