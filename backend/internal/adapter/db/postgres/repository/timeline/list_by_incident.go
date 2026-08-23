package timeline

import (
	"context"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repository) ListTimeline(ctx context.Context, incidentID uuid.UUID, params paging.Params) (paging.Slice[domain.TimelineEvent], error) {
	base := timelineBase(txcontext.From(ctx, r.db).WithContext(ctx), incidentID)
	query := applyTimelinePaging(base, params)
	var items []domain.TimelineEvent
	if err := query.Find(&items).Error; err != nil {
		return paging.Slice[domain.TimelineEvent]{}, mapError(err)
	}
	for i := range items {
		if err := items[i].Validate(); err != nil {
			return paging.Slice[domain.TimelineEvent]{}, mapInvalidEvent()
		}
	}
	countBase := timelineBase(txcontext.From(ctx, r.db).WithContext(ctx), incidentID)
	var total int64
	if err := countBase.Count(&total).Error; err != nil {
		return paging.Slice[domain.TimelineEvent]{}, mapError(err)
	}
	return paging.Slice[domain.TimelineEvent]{Items: items, Total: total}, nil
}

func applyTimelinePaging(query *gorm.DB, params paging.Params) *gorm.DB {
	order := timelineOrder(params)
	if order == "" {
		order = "timeline.occurred_at ASC, timeline.id ASC"
	}
	query = query.Order(order)
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	return query
}

func timelineOrder(params paging.Params) string {
	parts := make([]string, 0, len(params.Sort))
	for _, sort := range params.Sort {
		field := ""
		switch sort.Field {
		case "occurred_at":
			field = "timeline.occurred_at"
		case "id":
			field = "timeline.id"
		default:
			continue
		}
		direction := "ASC"
		if strings.EqualFold(string(sort.Direction), "desc") {
			direction = "DESC"
		}
		parts = append(parts, field+" "+direction)
	}
	return strings.Join(parts, ", ")
}

func timelineBase(db *gorm.DB, incidentID uuid.UUID) *gorm.DB {
	return db.Table(`(
        SELECT md5('incident_detected:' || incidents.id::text)::uuid AS id,
               incidents.id AS incident_id,
               'incident_detected' AS type,
               incidents.first_seen_at AS occurred_at,
               'Incident detected' AS summary,
               incidents.key AS detail
          FROM incidents
         WHERE incidents.id = ?
        UNION ALL
        SELECT md5('investigation_started:' || investigations.id::text)::uuid AS id,
               investigations.incident_id AS incident_id,
               'investigation_started' AS type,
               investigations.started_at AS occurred_at,
               'Investigation started' AS summary,
               investigations.id::text AS detail
          FROM investigations
         WHERE investigations.incident_id = ? AND investigations.started_at IS NOT NULL
        UNION ALL
        SELECT md5('tool_used:' || evidence.id::text)::uuid AS id,
               investigations.incident_id AS incident_id,
               'tool_used' AS type,
               evidence.created_at AS occurred_at,
               'Investigation evidence captured' AS summary,
               evidence.type::text AS detail
          FROM evidence
          JOIN investigations ON investigations.id = evidence.investigation_id
         WHERE investigations.incident_id = ?
        UNION ALL
        SELECT md5('root_cause_classified:' || investigations.id::text)::uuid AS id,
               investigations.incident_id AS incident_id,
               'root_cause_classified' AS type,
               investigations.finished_at AS occurred_at,
               'Root cause classified' AS summary,
               concat_ws(' / ', investigations.root_cause_status, investigations.resolution_status) AS detail
          FROM investigations
         WHERE investigations.incident_id = ? AND investigations.status = 'completed' AND investigations.finished_at IS NOT NULL
        UNION ALL
        SELECT md5('issue_created:' || github_issue_references.investigation_id::text)::uuid AS id,
               github_issue_references.incident_id AS incident_id,
               'issue_created' AS type,
               github_issue_references.created_at AS occurred_at,
               'GitHub Issue created' AS summary,
               github_issue_references.repository || '#' || github_issue_references.issue_number::text AS detail
          FROM github_issue_references
         WHERE github_issue_references.incident_id = ?
        UNION ALL
        SELECT md5('workflow_failed:' || investigations.id::text)::uuid AS id,
               investigations.incident_id AS incident_id,
               'workflow_failed' AS type,
               investigations.finished_at AS occurred_at,
               'Investigation workflow failed' AS summary,
               investigations.failure_user_message AS detail
          FROM investigations
         WHERE investigations.incident_id = ? AND investigations.status = 'failed' AND investigations.finished_at IS NOT NULL
    ) AS timeline`, incidentID, incidentID, incidentID, incidentID, incidentID, incidentID)
}
