package project

import "context"

type Update interface {
	Update(ctx context.Context, command UpdateCommand) (*Result, error)
}
