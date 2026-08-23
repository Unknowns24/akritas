package remediation

import (
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

// remediationRecord is private infrastructure state: domain.Remediation
// carries no GORM tags, so persistence maps through this record rather
// than reusing the domain struct directly (ADR-012 only covers tagged
// entities). Changes and ValidationResults stay in their own persistence
// surfaces; this record owns the Remediation lifecycle row and external PR
// reference.
type remediationRecord struct {
	ID                    uuid.UUID  `gorm:"column:id;type:uuid;primaryKey"`
	IncidentID            uuid.UUID  `gorm:"column:incident_id;type:uuid"`
	InvestigationID       *uuid.UUID `gorm:"column:investigation_id;type:uuid"`
	Status                string     `gorm:"column:status"`
	BranchName            string     `gorm:"column:branch_name"`
	ChangesSummary        string     `gorm:"column:changes_summary"`
	FailureUserMessage    string     `gorm:"column:failure_user_message"`
	PullRequestNumber     int        `gorm:"column:pull_request_number"`
	PullRequestURL        string     `gorm:"column:pull_request_url"`
	PullRequestRepository string     `gorm:"column:pull_request_repository"`
	PullRequestBranch     string     `gorm:"column:pull_request_branch"`
	PullRequestCreatedAt  *time.Time `gorm:"column:pull_request_created_at"`
	CreatedAt             time.Time  `gorm:"column:created_at"`
	UpdatedAt             time.Time  `gorm:"column:updated_at"`
}

func fromDomain(value *domain.Remediation) remediationRecord {
	record := remediationRecord{
		ID: value.ID, IncidentID: value.IncidentID, Status: string(value.Status),
		BranchName: value.BranchName, ChangesSummary: value.ChangesSummary,
		FailureUserMessage: value.FailureUserMessage, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt,
	}
	if value.InvestigationID != uuid.Nil {
		record.InvestigationID = &value.InvestigationID
	}
	if value.PullRequestReference != nil {
		record.PullRequestNumber = value.PullRequestReference.Number
		record.PullRequestURL = value.PullRequestReference.URL
		record.PullRequestRepository = value.PullRequestReference.Repository
		record.PullRequestBranch = value.PullRequestReference.Branch
		record.PullRequestCreatedAt = &value.PullRequestReference.CreatedAt
	}
	return record
}

func (r remediationRecord) toDomain() *domain.Remediation {
	value := &domain.Remediation{
		ID: r.ID, IncidentID: r.IncidentID, Status: domain.RemediationStatus(r.Status),
		BranchName: r.BranchName, ChangesSummary: r.ChangesSummary,
		Changes: []domain.CodeChange{}, ValidationResults: []domain.ValidationResult{},
		FailureUserMessage: r.FailureUserMessage, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
	if r.InvestigationID != nil {
		value.InvestigationID = *r.InvestigationID
	}
	if r.PullRequestNumber > 0 || r.PullRequestURL != "" || r.PullRequestRepository != "" || r.PullRequestBranch != "" || r.PullRequestCreatedAt != nil {
		createdAt := time.Time{}
		if r.PullRequestCreatedAt != nil {
			createdAt = *r.PullRequestCreatedAt
		}
		value.PullRequestReference = &domain.PullRequestReference{
			Number: r.PullRequestNumber, URL: r.PullRequestURL, Repository: r.PullRequestRepository,
			Branch: r.PullRequestBranch, CreatedAt: createdAt,
		}
	}
	return value
}
