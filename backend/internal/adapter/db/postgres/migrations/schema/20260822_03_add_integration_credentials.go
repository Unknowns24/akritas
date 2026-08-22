package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_03_AddIntegrationCredentials() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_03_add_integration_credentials",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`CREATE TABLE integration_credentials (
                id uuid PRIMARY KEY,
                owner_type varchar(64) NOT NULL,
                owner_id uuid NOT NULL,
                secret_kind varchar(64) NOT NULL,
                ciphertext bytea NOT NULL,
                nonce bytea NOT NULL,
                encryption_version smallint NOT NULL CHECK (encryption_version > 0),
                created_at timestamptz NOT NULL,
                updated_at timestamptz NOT NULL,
                CONSTRAINT uq_integration_credentials_owner UNIQUE (owner_type, owner_id, secret_kind),
                CONSTRAINT ck_integration_credentials_ciphertext CHECK (octet_length(ciphertext) > 16),
                CONSTRAINT ck_integration_credentials_nonce CHECK (octet_length(nonce) = 12)
            )`).Error
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS integration_credentials").Error },
	}
}
