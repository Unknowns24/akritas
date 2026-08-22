package project

import (
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type CreateCommand struct {
	Name                    string
	Description             string
	GitHubAccountID         uuid.UUID
	RepositoryIdentifier    string
	DefaultBranch           string
	DokployServerID         uuid.UUID
	ApplicationIdentifier   string
	MonitoringConfiguration domain.MonitoringConfiguration
}

type UpdateCommand struct {
	ID                    uuid.UUID
	Name                  *string
	Description           *string
	GitHubAccountID       *uuid.UUID
	RepositoryIdentifier  *string
	DefaultBranch         *string
	DokployServerID       *uuid.UUID
	ApplicationIdentifier *string
}

type MonitoringCommand struct {
	ProjectID               uuid.UUID
	MonitoringConfiguration domain.MonitoringConfiguration
}

type Result struct {
	Project               *domain.Project
	BuiltInDetectionRules []domain.BuiltInDetectionRule
}

func NewResult(project *domain.Project) *Result {
	if project == nil {
		return nil
	}
	return &Result{Project: project, BuiltInDetectionRules: domain.AllBuiltInDetectionRules()}
}
