package project

import (
	"context"

	"github.com/google/uuid"
)

type Get interface {
	Get(ctx context.Context, id uuid.UUID) (*Result, error)
}
