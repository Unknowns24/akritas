package postgres

import (
	"context"

	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type transactor struct {
	db *gorm.DB
}

func NewTransactor(db *gorm.DB) out.Transactor {
	return &transactor{db: db}
}

func (t *transactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(txcontext.With(ctx, tx))
	})
}
