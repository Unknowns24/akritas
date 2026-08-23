package monitoring

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) GetCurrentCheckpoint(ctx context.Context, projectID uuid.UUID, lock bool) (*domain.MonitoringCheckpoint, error) {
	var checkpoint domain.MonitoringCheckpoint
	query := txcontext.From(ctx, r.db).WithContext(ctx).Table("monitoring_checkpoints")
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	err := query.Where("project_id = ? AND is_current", projectID).Take(&checkpoint).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, mapError(err)
	}
	return &checkpoint, nil
}

func (r *Repository) CreateCheckpoint(ctx context.Context, checkpoint *domain.MonitoringCheckpoint) error {
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("monitoring_checkpoints").Create(checkpoint).Error; err != nil {
		return mapError(err)
	}
	return nil
}

func (r *Repository) RotateCheckpoint(ctx context.Context, checkpoint *domain.MonitoringCheckpoint) error {
	db := txcontext.From(ctx, r.db).WithContext(ctx)
	result := db.Table("monitoring_checkpoints").Where("project_id = ? AND is_current AND next_finalize_at IS NULL", checkpoint.ProjectID).Update("is_current", false)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		var currentCount int64
		if err := db.Table("monitoring_checkpoints").Where("project_id = ? AND is_current", checkpoint.ProjectID).Count(&currentCount).Error; err != nil {
			return mapError(err)
		}
		if currentCount > 0 {
			return domain.ErrMonitoringConcurrentModification
		}
	}
	if err := db.Table("monitoring_checkpoints").Create(checkpoint).Error; err != nil {
		return mapError(err)
	}
	return nil
}

func (r *Repository) UpdateCheckpoint(ctx context.Context, checkpoint *domain.MonitoringCheckpoint, expectedVersion int64) error {
	result := txcontext.From(ctx, r.db).WithContext(ctx).Table("monitoring_checkpoints").Where("id = ? AND version = ?", checkpoint.ID, expectedVersion).
		Select("initial_backfill_pending", "cursor_timestamp", "cursor_ordinal", "cursor_content_hash", "version", "assembly_state", "next_finalize_at", "updated_at").Updates(checkpoint)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected != 1 {
		return domain.ErrMonitoringConcurrentModification
	}
	return nil
}
