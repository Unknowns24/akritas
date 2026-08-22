package schema

import (
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_01_AddIntegrationLookups() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_01_add_integration_lookups",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&domain.GitHubAccount{}, &domain.DokployServer{})
		},
		Rollback: func(tx *gorm.DB) error {
			if tx.Migrator().HasTable(&domain.GitHubAccount{}) {
				if err := tx.Migrator().DropTable(&domain.GitHubAccount{}); err != nil {
					return err
				}
			}
			if tx.Migrator().HasTable(&domain.DokployServer{}) {
				return tx.Migrator().DropTable(&domain.DokployServer{})
			}
			return nil
		},
	}
}
