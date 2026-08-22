package out

import "context"

type AdministratorRepository interface {
	ExistsActive(ctx context.Context) (bool, error)
}
