package investigationtools

import (
	"context"
	"fmt"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// Scope is the repository investigation tools are allowed to touch.
type Scope struct {
	Account domain.GitHubAccount
	Owner   string
	Name    string
	Branch  string
}

type Resolver struct {
	incidents portsout.IncidentReader
	projects  portsout.ProjectStore
	accounts  portsout.GitHubAccountReader
}

func NewResolver(incidents portsout.IncidentReader, projects portsout.ProjectStore, accounts portsout.GitHubAccountReader) *Resolver {
	return &Resolver{incidents: incidents, projects: projects, accounts: accounts}
}

func (r *Resolver) Resolve(ctx context.Context, investigation domain.Investigation) (Scope, error) {
	if r == nil || r.incidents == nil || r.projects == nil || r.accounts == nil {
		return Scope{}, fmt.Errorf("investigation tool resolver is not configured")
	}
	incident, err := r.incidents.Get(ctx, investigation.IncidentID)
	if err != nil {
		return Scope{}, err
	}
	project, err := r.projects.Get(ctx, incident.ProjectID)
	if err != nil {
		return Scope{}, err
	}
	account, err := r.accounts.Get(ctx, project.GitHubRepository.GitHubAccountID)
	if err != nil {
		return Scope{}, err
	}
	return Scope{
		Account: *account,
		Owner:   project.GitHubRepository.Owner,
		Name:    project.GitHubRepository.Name,
		Branch:  project.GitHubRepository.DefaultBranch,
	}, nil
}
