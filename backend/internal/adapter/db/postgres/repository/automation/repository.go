package automation

import (
	"context"
	"errors"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"gorm.io/gorm"
)

var ErrInvalidRepository = errors.New("invalid automation repository configuration")

type Repository struct{ db *gorm.DB }

type record struct {
	ID                     int       `gorm:"column:id"`
	AutomaticInvestigation bool      `gorm:"column:automatic_investigation"`
	AutomaticRemediation   bool      `gorm:"column:automatic_remediation"`
	AutomaticPullRequest   bool      `gorm:"column:automatic_pull_request"`
	UpdatedAt              time.Time `gorm:"column:updated_at"`
}

func New(db *gorm.DB) (*Repository, error) {
	if db == nil {
		return nil, ErrInvalidRepository
	}
	return &Repository{db: db}, nil
}

func (r *Repository) Get(ctx context.Context) (domain.AutomationPolicy, error) {
	var row record
	if err := r.db.WithContext(ctx).Table("automation_policy").Where("id = 1").Take(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.DefaultAutomationPolicy(), nil
		}
		return domain.AutomationPolicy{}, domain.ErrInvalidAutomationPolicy.Wrap(err)
	}
	return domain.NewAutomationPolicy(row.AutomaticInvestigation, row.AutomaticRemediation, row.AutomaticPullRequest)
}

func (r *Repository) Put(ctx context.Context, value domain.AutomationPolicy) error {
	if err := value.Validate(); err != nil {
		return err
	}
	row := record{
		ID: 1, AutomaticInvestigation: value.AutomaticInvestigation,
		AutomaticRemediation: value.AutomaticRemediation, AutomaticPullRequest: value.AutomaticPullRequest,
		UpdatedAt: time.Now().UTC(),
	}
	return r.db.WithContext(ctx).Table("automation_policy").Save(&row).Error
}
