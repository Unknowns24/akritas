package project

import (
	"context"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type memoryProjects struct {
	byID map[uuid.UUID]*domain.Project
}

func newMemoryProjects() *memoryProjects {
	return &memoryProjects{byID: map[uuid.UUID]*domain.Project{}}
}

func (m *memoryProjects) Create(_ context.Context, project *domain.Project) error {
	clone := *project
	m.byID[project.ID] = &clone
	return nil
}

func (m *memoryProjects) GetByID(_ context.Context, id uuid.UUID) (*domain.Project, error) {
	project, ok := m.byID[id]
	if !ok {
		return nil, apperr.ErrProjectNotFound
	}
	clone := *project
	return &clone, nil
}

func (m *memoryProjects) GetByNormalizedName(_ context.Context, name string) (*domain.Project, error) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	for _, project := range m.byID {
		if strings.ToLower(project.Name) == normalized {
			clone := *project
			return &clone, nil
		}
	}
	return nil, apperr.ErrProjectNotFound
}

func (m *memoryProjects) GetByDokployApplication(_ context.Context, serverID uuid.UUID, applicationIdentifier string) (*domain.Project, error) {
	identifier := strings.TrimSpace(applicationIdentifier)
	for _, project := range m.byID {
		if project.DokployApplication.DokployServerID == serverID && project.DokployApplication.ApplicationIdentifier == identifier {
			clone := *project
			return &clone, nil
		}
	}
	return nil, apperr.ErrProjectNotFound
}

func (m *memoryProjects) List(_ context.Context, query paging.ListQuery) ([]domain.Project, int64, error) {
	matches := make([]domain.Project, 0, len(m.byID))
	for _, project := range m.byID {
		if query.NameLike != "" && !strings.Contains(strings.ToLower(project.Name), strings.ToLower(query.NameLike)) {
			continue
		}
		if len(query.MonitoringStatusIn) > 0 && !containsStatus(query.MonitoringStatusIn, project.MonitoringStatus) {
			continue
		}
		matches = append(matches, *project)
	}
	total := int64(len(matches))
	start := query.Offset
	if start > len(matches) {
		start = len(matches)
	}
	end := start + query.Limit
	if end > len(matches) {
		end = len(matches)
	}
	return matches[start:end], total, nil
}

func (m *memoryProjects) Update(_ context.Context, project *domain.Project) error {
	if _, ok := m.byID[project.ID]; !ok {
		return apperr.ErrProjectNotFound
	}
	clone := *project
	m.byID[project.ID] = &clone
	return nil
}

func (m *memoryProjects) CountByGitHubAccountID(_ context.Context, accountID uuid.UUID) (int64, error) {
	var total int64
	for _, project := range m.byID {
		if project.GitHubRepository.GitHubAccountID == accountID {
			total++
		}
	}
	return total, nil
}

func containsStatus(allowed []domain.MonitoringStatus, status domain.MonitoringStatus) bool {
	for _, candidate := range allowed {
		if candidate == status {
			return true
		}
	}
	return false
}

type memoryAccounts struct {
	byID map[uuid.UUID]*domain.GitHubAccount
}

func (m *memoryAccounts) GetByID(_ context.Context, id uuid.UUID) (*domain.GitHubAccount, error) {
	account, ok := m.byID[id]
	if !ok {
		return nil, apperr.ErrGitHubAccountNotFound
	}
	clone := *account
	return &clone, nil
}

type memoryServers struct {
	byID map[uuid.UUID]*domain.DokployServer
}

func (m *memoryServers) GetByID(_ context.Context, id uuid.UUID) (*domain.DokployServer, error) {
	server, ok := m.byID[id]
	if !ok {
		return nil, apperr.ErrDokployServerNotFound
	}
	clone := *server
	return &clone, nil
}

type memorySnapshots struct{}

func (memorySnapshots) ResolveGitHubRepository(account *domain.GitHubAccount, repositoryIdentifier, defaultBranch string) (domain.GitHubRepository, error) {
	if account == nil {
		return domain.GitHubRepository{}, apperr.ErrGitHubAccountNotFound
	}
	identifier := strings.TrimSpace(repositoryIdentifier)
	if identifier == "" || strings.TrimSpace(defaultBranch) == "" {
		return domain.GitHubRepository{}, apperr.ErrRepositoryNotResolvable
	}
	owner, name := strings.TrimSpace(account.AccountIdentifier), identifier
	if parts := strings.Split(identifier, "/"); len(parts) == 2 {
		owner, name = strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	} else if strings.Contains(identifier, "/") {
		return domain.GitHubRepository{}, apperr.ErrRepositoryNotResolvable
	}
	if owner == "" || name == "" {
		return domain.GitHubRepository{}, apperr.ErrRepositoryNotResolvable
	}
	return domain.NewGitHubRepository(
		account.ID, identifier, owner, name, defaultBranch, false,
		"https://github.com/"+owner+"/"+name,
	)
}

func (memorySnapshots) ResolveDokployApplication(server *domain.DokployServer, applicationIdentifier string) (domain.DokployApplication, error) {
	identifier := strings.TrimSpace(applicationIdentifier)
	return domain.NewDokployApplication(server.ID, identifier, identifier, identifier, "", domain.DokployApplicationUnknown)
}
