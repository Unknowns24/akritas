package project

import (
	"time"

	dokploydto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dokploy"
	githubdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/github"
)

type ProjectSummaryDTO struct {
	ID               string                        `json:"id"`
	Name             string                        `json:"name"`
	Description      string                        `json:"description"`
	MonitoringStatus string                        `json:"monitoring_status"`
	HealthStatus     string                        `json:"health_status"`
	GitHubRepository githubdto.GitHubRepositoryDTO `json:"github_repository"`
	DokploySource    dokploydto.DokploySourceDTO   `json:"dokploy_source"`
	LastObservedAt   *time.Time                    `json:"last_observed_at,omitempty"`
	CreatedAt        time.Time                     `json:"created_at"`
	UpdatedAt        time.Time                     `json:"updated_at"`
}
