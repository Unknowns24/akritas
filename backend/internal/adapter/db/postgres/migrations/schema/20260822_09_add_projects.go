package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_09_AddProjects() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_09_add_projects",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`CREATE TABLE projects (
                id uuid PRIMARY KEY,
                name varchar(150) NOT NULL,
                description varchar(2000) NOT NULL DEFAULT '',
                monitoring_status varchar(32) NOT NULL CHECK (monitoring_status IN ('disabled','starting','monitoring','degraded','error')),
                health_status varchar(32) NOT NULL CHECK (health_status IN ('healthy','warning','critical','unknown')),
                github_account_id uuid NOT NULL REFERENCES github_accounts(id) ON DELETE RESTRICT,
                repository_identifier varchar(255) NOT NULL,
                repository_owner varchar(255) NOT NULL,
                repository_name varchar(255) NOT NULL,
                repository_full_name varchar(511) NOT NULL,
                default_branch varchar(255) NOT NULL,
                repository_private boolean NOT NULL,
                repository_html_url text NOT NULL,
                dokploy_server_id uuid NOT NULL REFERENCES dokploy_servers(id) ON DELETE RESTRICT,
                application_identifier varchar(255) NOT NULL,
                instance_identifier varchar(255) NOT NULL,
                application_display_name varchar(255) NOT NULL,
                application_environment varchar(255) NOT NULL DEFAULT '',
                application_status varchar(32) NOT NULL CHECK (application_status IN ('running','stopped','degraded','unknown')),
                monitoring_enabled boolean NOT NULL DEFAULT false,
                error_patterns jsonb NOT NULL DEFAULT '[]'::jsonb,
                ignored_patterns jsonb NOT NULL DEFAULT '[]'::jsonb,
                grouping_window_ns bigint NOT NULL CHECK (grouping_window_ns > 0),
                context_before integer NOT NULL CHECK (context_before BETWEEN 0 AND 1000),
                context_after integer NOT NULL CHECK (context_after BETWEEN 0 AND 1000),
                last_observed_at timestamptz NULL,
                created_at timestamptz NOT NULL,
                updated_at timestamptz NOT NULL,
                CONSTRAINT ck_projects_monitoring_state CHECK (
                    (monitoring_enabled = false AND monitoring_status = 'disabled') OR
                    (monitoring_enabled = true AND monitoring_status <> 'disabled')
                ),
				CONSTRAINT uq_projects_dokploy_application UNIQUE (dokploy_server_id, application_identifier)
			)`).Error; err != nil {
				return err
			}
			for _, statement := range []string{
				"CREATE UNIQUE INDEX uq_projects_name_lower ON projects (LOWER(name))",
				"CREATE INDEX ix_projects_github_account_id ON projects (github_account_id)",
				"CREATE INDEX ix_projects_created_at_id ON projects (created_at, id)",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS projects").Error },
	}
}
