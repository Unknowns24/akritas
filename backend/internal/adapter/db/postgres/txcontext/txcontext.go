// Package txcontext propagates an active GORM transaction through
// context.Context so unrelated repositories can participate in the same
// transaction without depending on each other directly.
package txcontext

import (
	"context"

	"gorm.io/gorm"
)

type key struct{}

// With returns a context carrying the active transaction.
func With(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, key{}, tx)
}

// From returns the active transaction from ctx, or fallback if none is set.
func From(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(key{}).(*gorm.DB); ok {
		return tx
	}
	return fallback
}
