package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_08_AddDokployComposeSources() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_08_add_dokploy_compose_sources",
		Migrate: func(tx *gorm.DB) error {
			statements := []string{
				"ALTER TABLE dokploy_servers ADD COLUMN compose_count integer NOT NULL DEFAULT 0 CHECK (compose_count >= 0)",
				"ALTER TABLE projects DROP CONSTRAINT uq_projects_dokploy_application",
				"ALTER TABLE projects RENAME COLUMN application_identifier TO source_resource_identifier",
				"ALTER TABLE projects RENAME COLUMN instance_identifier TO source_instance_identifier",
				"ALTER TABLE projects RENAME COLUMN application_display_name TO source_display_name",
				"ALTER TABLE projects RENAME COLUMN application_environment TO source_environment",
				"ALTER TABLE projects RENAME COLUMN application_status TO source_status",
				"ALTER TABLE projects ADD COLUMN source_type varchar(32) NOT NULL DEFAULT 'application'",
				"ALTER TABLE projects ADD COLUMN source_service_name varchar(255) NOT NULL DEFAULT ''",
				"ALTER TABLE projects ADD COLUMN source_runtime_type varchar(32) NOT NULL DEFAULT ''",
				"ALTER TABLE projects ADD COLUMN source_provider_server_id varchar(255) NOT NULL DEFAULT ''",
				`ALTER TABLE projects ADD CONSTRAINT ck_projects_dokploy_source CHECK (
                    (source_type = 'application' AND source_service_name = '' AND source_runtime_type = '' AND source_provider_server_id = '') OR
                    (source_type = 'compose_service' AND source_service_name <> '' AND source_runtime_type IN ('docker-compose','stack'))
                )`,
				"CREATE UNIQUE INDEX uq_projects_dokploy_source ON projects (dokploy_server_id, source_type, source_resource_identifier, COALESCE(source_service_name, ''))",
				"ALTER TABLE monitoring_checkpoints RENAME COLUMN source_application_id TO source_resource_id",
				"ALTER TABLE monitoring_checkpoints ADD COLUMN source_type varchar(32) NOT NULL DEFAULT 'application' CHECK (source_type IN ('application','compose_service'))",
				"ALTER TABLE monitoring_checkpoints ADD COLUMN source_service_name varchar(255) NOT NULL DEFAULT ''",
				"ALTER TABLE log_events RENAME COLUMN source_application_id TO source_resource_id",
				"ALTER TABLE log_events ADD COLUMN source_type varchar(32) NOT NULL DEFAULT 'application' CHECK (source_type IN ('application','compose_service'))",
				"ALTER TABLE log_events ADD COLUMN source_service_name varchar(255) NOT NULL DEFAULT ''",
			}
			for _, statement := range statements {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			statements := []string{
				`DO $$ BEGIN
                    IF EXISTS (SELECT 1 FROM projects WHERE source_type <> 'application') THEN
                        RAISE EXCEPTION 'cannot rollback Dokploy compose sources while compose_service projects exist';
                    END IF;
                END $$`,
				"ALTER TABLE log_events DROP COLUMN source_service_name",
				"ALTER TABLE log_events DROP COLUMN source_type",
				"ALTER TABLE log_events RENAME COLUMN source_resource_id TO source_application_id",
				"ALTER TABLE monitoring_checkpoints DROP COLUMN source_service_name",
				"ALTER TABLE monitoring_checkpoints DROP COLUMN source_type",
				"ALTER TABLE monitoring_checkpoints RENAME COLUMN source_resource_id TO source_application_id",
				"DROP INDEX uq_projects_dokploy_source",
				"ALTER TABLE projects DROP CONSTRAINT ck_projects_dokploy_source",
				"ALTER TABLE projects DROP COLUMN source_provider_server_id",
				"ALTER TABLE projects DROP COLUMN source_runtime_type",
				"ALTER TABLE projects DROP COLUMN source_service_name",
				"ALTER TABLE projects DROP COLUMN source_type",
				"ALTER TABLE projects RENAME COLUMN source_status TO application_status",
				"ALTER TABLE projects RENAME COLUMN source_environment TO application_environment",
				"ALTER TABLE projects RENAME COLUMN source_display_name TO application_display_name",
				"ALTER TABLE projects RENAME COLUMN source_instance_identifier TO instance_identifier",
				"ALTER TABLE projects RENAME COLUMN source_resource_identifier TO application_identifier",
				"ALTER TABLE projects ADD CONSTRAINT uq_projects_dokploy_application UNIQUE (dokploy_server_id, application_identifier)",
				"ALTER TABLE dokploy_servers DROP COLUMN compose_count",
			}
			for _, statement := range statements {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
}
