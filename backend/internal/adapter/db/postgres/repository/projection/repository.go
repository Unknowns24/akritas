package projection

import (
	"context"
	"strings"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) count(ctx context.Context, table, where string, args ...any) (int, error) {
	var total int64
	query := txcontext.From(ctx, r.db).WithContext(ctx).Table(table)
	if strings.TrimSpace(where) != "" {
		query = query.Where(where, args...)
	}
	if err := query.Count(&total).Error; err != nil {
		return 0, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	return int(total), nil
}

func (r *Repository) CountGitHubAccounts(ctx context.Context) (int, error) {
	return r.count(ctx, "github_accounts", "")
}

func (r *Repository) CountDokployServers(ctx context.Context) (int, error) {
	return r.count(ctx, "dokploy_servers", "")
}

func (r *Repository) CountMonitoredProjects(ctx context.Context) (int, error) {
	return r.count(ctx, "projects", "monitoring_enabled = true")
}

func (r *Repository) CountActiveIncidents(ctx context.Context) (int, error) {
	return r.count(ctx, "incidents", "phase NOT IN ?", []string{string(domain.IncidentPhaseCompleted), string(domain.IncidentPhaseFailed)})
}

func (r *Repository) CountCompletedIncidents(ctx context.Context) (int, error) {
	return r.count(ctx, "incidents", "phase = ?", string(domain.IncidentPhaseCompleted))
}

func (r *Repository) CountPullRequestsCreated(ctx context.Context) (int, error) {
	return r.count(ctx, "remediations", "pull_request_number > 0")
}

type incidentRow struct {
	domain.Incident `gorm:"embedded"`
	ProjectName     string `gorm:"column:project_name"`
}

func (r incidentRow) domain() domain.Incident {
	value := r.Incident
	value.Project = &domain.ProjectReference{ID: value.ProjectID, Name: r.ProjectName}
	return value
}

func (r *Repository) ListActiveIncidents(ctx context.Context, limit int) ([]domain.Incident, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	var rows []incidentRow
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("incidents").
		Select("incidents.*, projects.name AS project_name").
		Joins("JOIN projects ON projects.id = incidents.project_id").
		Where("incidents.phase NOT IN ?", []string{string(domain.IncidentPhaseCompleted), string(domain.IncidentPhaseFailed)}).
		Order("incidents.last_seen_at DESC, incidents.id DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, domain.ErrIncidentNotFound.Wrap(err)
	}
	items := make([]domain.Incident, 0, len(rows))
	for _, row := range rows {
		items = append(items, row.domain())
	}
	return items, nil
}

type activityRow struct {
	ID          uuid.UUID  `gorm:"column:id"`
	Type        string     `gorm:"column:type"`
	ProjectID   *uuid.UUID `gorm:"column:project_id"`
	ProjectName string     `gorm:"column:project_name"`
	IncidentID  *uuid.UUID `gorm:"column:incident_id"`
	IncidentKey string     `gorm:"column:incident_key"`
	OccurredAt  time.Time  `gorm:"column:occurred_at"`
	Summary     string     `gorm:"column:summary"`
	UserMessage string     `gorm:"column:user_message"`
}

func (r *Repository) ListActivity(ctx context.Context, params paging.Params) (paging.Slice[domain.ActivityEvent], error) {
	limit := params.Limit
	if limit < 1 || limit > 100 {
		limit = 25
	}
	typeFilter := parseCSV(params.Filters["type_in"])
	projectID := strings.TrimSpace(params.Filters["project_id_eq"])
	incidentID := strings.TrimSpace(params.Filters["incident_id_eq"])
	query := `
		SELECT * FROM (
			SELECT i.id, 'incident' AS type, i.project_id, p.name AS project_name, i.id AS incident_id, i.key AS incident_key, i.last_seen_at AS occurred_at, i.title AS summary, i.summary AS user_message
			FROM incidents i JOIN projects p ON p.id = i.project_id
			UNION ALL
			SELECT o.id, o.type AS type, NULL::uuid AS project_id, '' AS project_name, NULL::uuid AS incident_id, '' AS incident_key, o.updated_at AS occurred_at, o.user_message AS summary, o.user_message
			FROM operations o
			UNION ALL
			SELECT r.id, 'remediation' AS type, i.project_id, p.name AS project_name, r.incident_id, i.key AS incident_key, r.updated_at AS occurred_at, 'Remediation ' || r.status AS summary, r.failure_user_message AS user_message
			FROM remediations r JOIN incidents i ON i.id = r.incident_id JOIN projects p ON p.id = i.project_id
			UNION ALL
			SELECT r.id, 'pull_request' AS type, i.project_id, p.name AS project_name, r.incident_id, i.key AS incident_key, r.pull_request_created_at AS occurred_at, r.pull_request_url AS summary, r.pull_request_repository AS user_message
			FROM remediations r JOIN incidents i ON i.id = r.incident_id JOIN projects p ON p.id = i.project_id WHERE r.pull_request_number > 0
		) activity WHERE occurred_at IS NOT NULL`
	args := []any{}
	if len(typeFilter) > 0 {
		query += " AND type IN ?"
		args = append(args, typeFilter)
	}
	if projectID != "" {
		query += " AND project_id = ?"
		args = append(args, projectID)
	}
	if incidentID != "" {
		query += " AND incident_id = ?"
		args = append(args, incidentID)
	}
	countQuery := "SELECT count(*) FROM (" + query + ") count_activity"
	var total int64
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return paging.Slice[domain.ActivityEvent]{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	query += " ORDER BY occurred_at DESC, id DESC LIMIT ?"
	args = append(args, limit)
	var rows []activityRow
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return paging.Slice[domain.ActivityEvent]{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	items := make([]domain.ActivityEvent, 0, len(rows))
	for _, row := range rows {
		item := domain.ActivityEvent{
			ID: row.ID, Type: domain.ActivityType(row.Type), IncidentID: row.IncidentID, IncidentKey: row.IncidentKey,
			OccurredAt: row.OccurredAt, Summary: row.Summary, UserMessage: row.UserMessage,
		}
		if row.ProjectID != nil {
			item.Project = &domain.ProjectReference{ID: *row.ProjectID, Name: row.ProjectName}
		}
		items = append(items, item)
	}
	return paging.Slice[domain.ActivityEvent]{Items: items, Total: total}, nil
}

func (r *Repository) FindLastSystemDiagnostics(ctx context.Context) (*domain.Operation, error) {
	var op domain.Operation
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("operations").
		Where("type = ?", string(domain.OperationTypeSystemDiagnostics)).
		Order("updated_at DESC, id DESC").
		Take(&op).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, domain.ErrOperationNotFound.Wrap(err)
	}
	return &op, nil
}

type pullRequestRow struct {
	ID                    uuid.UUID  `gorm:"column:id"`
	ProjectID             uuid.UUID  `gorm:"column:project_id"`
	ProjectName           string     `gorm:"column:project_name"`
	IncidentID            uuid.UUID  `gorm:"column:incident_id"`
	IncidentKey           string     `gorm:"column:incident_key"`
	RemediationID         uuid.UUID  `gorm:"column:remediation_id"`
	IssueNumber           *int       `gorm:"column:issue_number"`
	IssueURL              string     `gorm:"column:issue_url"`
	IssueRepository       string     `gorm:"column:issue_repository"`
	IssueCreatedAt        *time.Time `gorm:"column:issue_created_at"`
	PullRequestNumber     int        `gorm:"column:pull_request_number"`
	PullRequestURL        string     `gorm:"column:pull_request_url"`
	PullRequestRepository string     `gorm:"column:pull_request_repository"`
	PullRequestBranch     string     `gorm:"column:pull_request_branch"`
	PullRequestCreatedAt  time.Time  `gorm:"column:pull_request_created_at"`
	Title                 string     `gorm:"column:title"`
	ChangesSummary        string     `gorm:"column:changes_summary"`
}

func pullRequestSelect() string {
	return `r.id, p.id AS project_id, p.name AS project_name, i.id AS incident_id, i.key AS incident_key,
		r.id AS remediation_id, gir.issue_number, gir.issue_url, gir.repository AS issue_repository, gir.created_at AS issue_created_at,
		r.pull_request_number, r.pull_request_url, r.pull_request_repository, r.pull_request_branch, r.pull_request_created_at,
		i.title, r.changes_summary`
}

func (r *Repository) pullRequestBase(ctx context.Context) *gorm.DB {
	return txcontext.From(ctx, r.db).WithContext(ctx).Table("remediations r").
		Select(pullRequestSelect()).
		Joins("JOIN incidents i ON i.id = r.incident_id").
		Joins("JOIN projects p ON p.id = i.project_id").
		Joins("LEFT JOIN github_issue_references gir ON gir.incident_id = i.id").
		Where("r.pull_request_number > 0")
}

func (r *Repository) ListPullRequests(ctx context.Context, params paging.Params) (paging.Slice[domain.PullRequestProjection], error) {
	limit := params.Limit
	if limit < 1 || limit > 100 {
		limit = 25
	}
	query := r.pullRequestBase(ctx)
	if value := strings.TrimSpace(params.Filters["project_id_eq"]); value != "" {
		query = query.Where("i.project_id = ?", value)
	}
	if value := strings.TrimSpace(params.Filters["incident_id_eq"]); value != "" {
		query = query.Where("i.id = ?", value)
	}
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return paging.Slice[domain.PullRequestProjection]{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	var rows []pullRequestRow
	if err := query.Order("r.pull_request_created_at DESC, r.id DESC").Limit(limit).Scan(&rows).Error; err != nil {
		return paging.Slice[domain.PullRequestProjection]{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	items := make([]domain.PullRequestProjection, 0, len(rows))
	for _, row := range rows {
		items = append(items, row.toDomain())
	}
	return paging.Slice[domain.PullRequestProjection]{Items: items, Total: total}, nil
}

func (r *Repository) GetPullRequest(ctx context.Context, id uuid.UUID) (*domain.PullRequestProjection, error) {
	var row pullRequestRow
	if err := r.pullRequestBase(ctx).Where("r.id = ?", id).Take(&row).Error; err != nil {
		return nil, domain.ErrRemediationNotFound.Wrap(err)
	}
	value := row.toDomain()
	return &value, nil
}

func (r pullRequestRow) toDomain() domain.PullRequestProjection {
	value := domain.PullRequestProjection{
		ID: r.ID, Project: domain.ProjectReference{ID: r.ProjectID, Name: r.ProjectName},
		IncidentID: r.IncidentID, IncidentKey: r.IncidentKey, RemediationID: r.RemediationID,
		Reference: domain.PullRequestReference{
			Number: r.PullRequestNumber, URL: r.PullRequestURL, Repository: r.PullRequestRepository,
			Branch: r.PullRequestBranch, CreatedAt: r.PullRequestCreatedAt,
		},
		Title: r.Title, ChangesSummary: r.ChangesSummary, CreatedAt: r.PullRequestCreatedAt,
	}
	if r.IssueNumber != nil && r.IssueCreatedAt != nil {
		value.IssueReference = &domain.GitHubIssueReference{Number: *r.IssueNumber, URL: r.IssueURL, Repository: r.IssueRepository, CreatedAt: *r.IssueCreatedAt}
	}
	return value
}

func parseCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	raw := strings.Split(value, ",")
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
