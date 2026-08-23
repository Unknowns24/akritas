package incident

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// Update persists workflow-owned Incident state without touching H2 grouping
// fields such as occurrence_count or last_seen_at.
func (r *Repository) Update(ctx context.Context, incident *domain.Incident) error {
	if incident == nil || incident.Validate() != nil {
		return domain.ErrInvalidIncident
	}
	result := txcontext.From(ctx, r.db).WithContext(ctx).Table("incidents").Where("id = ?", incident.ID).Updates(map[string]any{
		"phase":             incident.Phase,
		"terminal_outcome":  incident.TerminalOutcome,
		"summary":           incident.Summary,
		"root_cause_status": incident.RootCauseStatus,
		"resolution_status": incident.ResolutionStatus,
		"confidence":        incident.Confidence,
	})
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected != 1 {
		return domain.ErrIncidentNotFound
	}
	return nil
}
