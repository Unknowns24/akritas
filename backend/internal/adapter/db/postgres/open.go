package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type tableNamer struct {
	schema.NamingStrategy
}

func (n tableNamer) TableName(name string) string {
	switch name {
	case "GitHubAccount":
		return "github_accounts"
	case "DokployServer":
		return "dokploy_servers"
	case "Project":
		return "projects"
	default:
		return n.NamingStrategy.TableName(name)
	}
}

func Open(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("postgres DSN is required")
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Silent),
		NamingStrategy: tableNamer{NamingStrategy: schema.NamingStrategy{}},
	})
}
