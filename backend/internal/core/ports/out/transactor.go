package out

import "context"

// Transactor runs fn atomically. Repositories that participate detect the
// active transaction via the context WithinTransaction passes to fn and use
// it instead of their own connection.
type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
