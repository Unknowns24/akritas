package remediation

import (
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

// remediationRecord is private infrastructure state: domain.Remediation
// carries no GORM tags, so persistence maps through this record rather
// than reusing the domain struct directly (ADR-012 only covers tagged
// entities). It deliberately mirrors only the minimal columns this task's
// table has — Changes, ValidationResults and PullRequestReference are not
// persisted here; that is deferred to AKR-49/55+.
type remediationRecord struct {
	ID                 uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	IncidentID         uuid.UUID `gorm:"column:incident_id;type:uuid"`
	Status             string    `gorm:"column:status"`
	BranchName         string    `gorm:"column:branch_name"`
	ChangesSummary     string    `gorm:"column:changes_summary"`
	FailureUserMessage string    `gorm:"column:failure_user_message"`
	CreatedAt          time.Time `gorm:"column:created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at"`
}

func fromDomain(value *domain.Remediation) remediationRecord {
	return remediationRecord{
		ID: value.ID, IncidentID: value.IncidentID, Status: string(value.Status),
		BranchName: value.BranchName, ChangesSummary: value.ChangesSummary,
		FailureUserMessage: value.FailureUserMessage, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt,
	}
}

func (r remediationRecord) toDomain() *domain.Remediation {
	return &domain.Remediation{
		ID: r.ID, IncidentID: r.IncidentID, Status: domain.RemediationStatus(r.Status),
		BranchName: r.BranchName, ChangesSummary: r.ChangesSummary,
		Changes: []domain.CodeChange{}, ValidationResults: []domain.ValidationResult{},
		FailureUserMessage: r.FailureUserMessage, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}
