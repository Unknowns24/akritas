package githubaccount

import (
	"strings"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

type UseCase struct {
	store   portsout.GitHubAccountStore
	gateway portsout.GitHubGateway
	usage   portsout.IntegrationUsageReader
	newID   func() uuid.UUID
	now     func() time.Time
}

func New(store portsout.GitHubAccountStore, gateway portsout.GitHubGateway, usage portsout.IntegrationUsageReader, newID func() uuid.UUID, now func() time.Time) portsin.GitHubAccountUseCase {
	return &UseCase{store: store, gateway: gateway, usage: usage, newID: newID, now: now}
}

func wipe(value []byte) {
	for index := range value {
		value[index] = 0
	}
}

func classicScopesSatisfy(accountType domain.GitHubAccountType, scopes []string) bool {
	if len(scopes) == 0 {
		return true
	}
	available := make(map[string]bool, len(scopes))
	for _, scope := range scopes {
		available[strings.ToLower(strings.TrimSpace(scope))] = true
	}
	if !available["repo"] {
		return false
	}
	return accountType != domain.GitHubAccountOrganization || available["read:org"]
}
