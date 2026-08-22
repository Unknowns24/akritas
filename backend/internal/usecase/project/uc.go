package project

import (
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type UseCase struct {
	projects  out.ProjectRepository
	accounts  out.GitHubAccountRepository
	servers   out.DokployServerRepository
	snapshots out.IntegrationSnapshotResolver
	now       func() time.Time
}

func NewUseCase(
	projects out.ProjectRepository,
	accounts out.GitHubAccountRepository,
	servers out.DokployServerRepository,
	snapshots out.IntegrationSnapshotResolver,
) *UseCase {
	return &UseCase{
		projects: projects, accounts: accounts, servers: servers, snapshots: snapshots,
		now: func() time.Time { return time.Now().UTC() },
	}
}
