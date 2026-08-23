package qvacconfig

import (
	"context"
	"errors"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"gorm.io/gorm"
)

var ErrInvalidRepository = errors.New("invalid QVAC configuration repository")

type Repository struct{ db *gorm.DB }

type record struct {
	ID                       int       `gorm:"column:id"`
	EndpointURL              string    `gorm:"column:endpoint_url"`
	ConnectionTimeoutSeconds int       `gorm:"column:connection_timeout_seconds"`
	AuthenticationType       string    `gorm:"column:authentication_type"`
	BasicUsername            string    `gorm:"column:basic_username"`
	CredentialConfigured     bool      `gorm:"column:credential_configured"`
	UpdatedAt                time.Time `gorm:"column:updated_at"`
}

func New(db *gorm.DB) (*Repository, error) {
	if db == nil {
		return nil, ErrInvalidRepository
	}
	return &Repository{db: db}, nil
}

func (r *Repository) Get(ctx context.Context) (domain.QvacConfiguration, error) {
	var row record
	if err := r.db.WithContext(ctx).Table("qvac_configurations").Where("id = 1").Take(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.DefaultQvacConfiguration(time.Now().UTC()), nil
		}
		return domain.QvacConfiguration{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	return domain.NewQvacConfiguration(row.EndpointURL, row.ConnectionTimeoutSeconds, domain.QvacAuthenticationType(row.AuthenticationType), row.CredentialConfigured, row.BasicUsername, row.UpdatedAt)
}

func (r *Repository) Put(ctx context.Context, value domain.QvacConfiguration) error {
	if err := value.Validate(); err != nil {
		return err
	}
	row := record{
		ID: 1, EndpointURL: value.EndpointURL, ConnectionTimeoutSeconds: value.ConnectionTimeoutSeconds,
		AuthenticationType: string(value.AuthenticationType), BasicUsername: value.BasicUsername,
		CredentialConfigured: value.CredentialConfigured, UpdatedAt: value.UpdatedAt,
	}
	return r.db.WithContext(ctx).Table("qvac_configurations").Save(&row).Error
}
