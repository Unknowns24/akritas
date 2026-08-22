package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_02_AddDokployServers() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_02_add_dokploy_servers",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`CREATE TABLE dokploy_servers (
                id uuid PRIMARY KEY,
                name varchar(200) NOT NULL,
                base_url text NOT NULL,
                server_identifier char(64) NOT NULL,
                connection_status varchar(32) NOT NULL CHECK (connection_status IN ('pending','connected','authentication_failed','unavailable')),
                credential_configured boolean NOT NULL DEFAULT false,
                application_count integer NOT NULL DEFAULT 0 CHECK (application_count >= 0),
                last_synced_at timestamptz NULL,
                created_at timestamptz NOT NULL,
                updated_at timestamptz NOT NULL,
                CONSTRAINT uq_dokploy_servers_base_url UNIQUE (base_url),
                CONSTRAINT uq_dokploy_servers_identifier UNIQUE (server_identifier)
            )`).Error
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS dokploy_servers").Error },
	}
}
