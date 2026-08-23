package postgres

import (
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrInvalidDatabaseConfiguration = errors.New("invalid PostgreSQL configuration")

type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func Open(config Config) (*gorm.DB, error) {
	if config.DSN == "" {
		return nil, ErrInvalidDatabaseConfiguration
	}
	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, ErrInvalidDatabaseConfiguration
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, ErrInvalidDatabaseConfiguration
	}
	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, ErrInvalidDatabaseConfiguration
	}
	return db, nil
}
