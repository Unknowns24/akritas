package project

import "context"

type Create interface {
	Create(ctx context.Context, command CreateCommand) (*Result, error)
}
