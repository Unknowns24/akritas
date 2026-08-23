package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_11_AddRuntimeSettings() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_11_add_runtime_settings",
		Migrate: func(tx *gorm.DB) error {
			for _, statement := range []string{
				`CREATE TABLE automation_policy (
					id smallint PRIMARY KEY CHECK (id = 1),
					automatic_investigation boolean NOT NULL,
					automatic_remediation boolean NOT NULL,
					automatic_pull_request boolean NOT NULL,
					updated_at timestamptz NOT NULL,
					CONSTRAINT ck_automation_policy_dependencies CHECK (
						(automatic_investigation = true OR (automatic_remediation = false AND automatic_pull_request = false)) AND
						(automatic_remediation = true OR automatic_pull_request = false)
					)
				)`,
				`INSERT INTO automation_policy (id, automatic_investigation, automatic_remediation, automatic_pull_request, updated_at)
				 VALUES (1, true, true, true, now())`,
				`CREATE TABLE qvac_configurations (
					id smallint PRIMARY KEY CHECK (id = 1),
					endpoint_url text NOT NULL,
					connection_timeout_seconds integer NOT NULL CHECK (connection_timeout_seconds BETWEEN 1 AND 300),
					authentication_type varchar(16) NOT NULL CHECK (authentication_type IN ('none','bearer','basic')),
					basic_username varchar(255) NOT NULL DEFAULT '',
					credential_configured boolean NOT NULL DEFAULT false,
					updated_at timestamptz NOT NULL,
					CONSTRAINT ck_qvac_credential_state CHECK (
						(authentication_type = 'none' AND credential_configured = false AND basic_username = '') OR
						(authentication_type = 'bearer' AND credential_configured = true) OR
						(authentication_type = 'basic' AND credential_configured = true AND basic_username <> '')
					)
				)`,
				`INSERT INTO qvac_configurations (id, endpoint_url, connection_timeout_seconds, authentication_type, basic_username, credential_configured, updated_at)
				 VALUES (1, 'http://127.0.0.1:11434/v1', 30, 'none', '', false, now())`,
				"CREATE INDEX ix_operations_type_updated_at ON operations (type, updated_at DESC)",
				"CREATE INDEX ix_remediations_pull_request_created_at ON remediations (pull_request_created_at DESC) WHERE pull_request_number > 0",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for _, statement := range []string{
				"DROP INDEX IF EXISTS ix_remediations_pull_request_created_at",
				"DROP INDEX IF EXISTS ix_operations_type_updated_at",
				"DROP TABLE IF EXISTS qvac_configurations",
				"DROP TABLE IF EXISTS automation_policy",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
}
