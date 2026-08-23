package validationresult

import (
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

// validationResultRecord is private infrastructure state, following the
// same rationale as remediationRecord: domain.ValidationResult carries no
// GORM tags.
type validationResultRecord struct {
	ID             uuid.UUID  `gorm:"column:id;type:uuid;primaryKey"`
	RemediationID  uuid.UUID  `gorm:"column:remediation_id;type:uuid"`
	Type           string     `gorm:"column:type"`
	Name           string     `gorm:"column:name"`
	Status         string     `gorm:"column:status"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	StartedAt      *time.Time `gorm:"column:started_at"`
	FinishedAt     *time.Time `gorm:"column:finished_at"`
	Summary        string     `gorm:"column:summary"`
	OutputExcerpt  string     `gorm:"column:output_excerpt"`
	OutputRedacted bool       `gorm:"column:output_redacted"`
}

func fromDomain(value *domain.ValidationResult) validationResultRecord {
	return validationResultRecord{
		ID: value.ID, RemediationID: value.RemediationID, Type: string(value.Type), Name: value.Name,
		Status: string(value.Status), CreatedAt: value.CreatedAt, StartedAt: value.StartedAt, FinishedAt: value.FinishedAt,
		Summary: value.Summary, OutputExcerpt: value.OutputExcerpt, OutputRedacted: value.OutputRedacted,
	}
}

func (r validationResultRecord) toDomain() domain.ValidationResult {
	return domain.ValidationResult{
		ID: r.ID, RemediationID: r.RemediationID, Type: domain.ValidationType(r.Type), Name: r.Name,
		Status: domain.ValidationStatus(r.Status), CreatedAt: r.CreatedAt, StartedAt: r.StartedAt, FinishedAt: r.FinishedAt,
		Summary: r.Summary, OutputExcerpt: r.OutputExcerpt, OutputRedacted: r.OutputRedacted,
	}
}
