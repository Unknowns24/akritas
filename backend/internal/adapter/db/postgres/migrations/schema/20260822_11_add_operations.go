package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// operations is generic async-command infrastructure, not exclusive to
// investigation: resource_id is polymorphic (no foreign key) so remediation
// and pull_request can reuse the same table once they queue work through it.
func SCHEMA_20260822_11_AddOperations() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_11_add_operations",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`CREATE TABLE operations (
                id uuid PRIMARY KEY,
                type varchar(32) NOT NULL CHECK (type IN ('system_diagnostics','investigation','remediation','pull_request')),
                status varchar(32) NOT NULL CHECK (status IN ('queued','running','succeeded','failed')),
                resource_type varchar(32) NULL CHECK (resource_type IN ('system','investigation','remediation','pull_request')),
                resource_id uuid NULL,
                user_message varchar(1000) NOT NULL DEFAULT '',
                failure_code varchar(16) NULL,
                idempotency_key varchar(64) NULL,
                created_at timestamptz NOT NULL,
                updated_at timestamptz NOT NULL,
                finished_at timestamptz NULL,
                CONSTRAINT ck_operations_resource_pair CHECK ((resource_type IS NULL) = (resource_id IS NULL))
            )`).Error; err != nil {
				return err
			}
			for _, statement := range []string{
				"CREATE UNIQUE INDEX uq_operations_idempotency_key ON operations (idempotency_key) WHERE idempotency_key IS NOT NULL",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS operations").Error },
	}
}
