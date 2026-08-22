package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type CreateGitHubPATAccountCommand struct {
	DisplayName         string
	AccountType         domain.GitHubAccountType
	AccountIdentifier   string
	PersonalAccessToken string
}

type UpdateGitHubAccountCommand struct {
	DisplayName         *string
	PersonalAccessToken *string
}

type GitHubAccountUseCase interface {
	CreatePAT(context.Context, CreateGitHubPATAccountCommand) (*domain.GitHubAccount, error)
	Get(context.Context, uuid.UUID) (*domain.GitHubAccount, error)
	List(context.Context, paging.Params) (paging.Slice[domain.GitHubAccount], error)
	Update(context.Context, uuid.UUID, UpdateGitHubAccountCommand) (*domain.GitHubAccount, error)
	Delete(context.Context, uuid.UUID) error
	TestConnection(context.Context, uuid.UUID) (ConnectionTestResult, error)
	ListRepositories(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.GitHubRepository], error)
}
