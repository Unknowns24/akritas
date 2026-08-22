package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type GitHubAccountStore interface {
	CreateWithCredential(context.Context, *domain.GitHubAccount, SecretValue) error
	Get(context.Context, uuid.UUID) (*domain.GitHubAccount, error)
	List(context.Context, paging.Params) (paging.Slice[domain.GitHubAccount], error)
	Update(context.Context, *domain.GitHubAccount, *SecretValue) error
	UpdateConnection(context.Context, *domain.GitHubAccount) error
	Delete(context.Context, uuid.UUID) error
}

type GitHubPATValidationRequest struct {
	AccountType       domain.GitHubAccountType
	AccountIdentifier string
	Token             string
}

type GitHubPATValidation struct {
	AccountIdentifier string
	ClassicScopes     []string
}

type GitHubGateway interface {
	ValidatePAT(context.Context, GitHubPATValidationRequest) (GitHubPATValidation, error)
	TestConnection(context.Context, domain.GitHubAccount) (ProviderConnectionResult, error)
	ListRepositories(context.Context, domain.GitHubAccount, paging.Params) (paging.Slice[domain.GitHubRepository], error)
}
