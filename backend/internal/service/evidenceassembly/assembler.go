// Package evidenceassembly builds the initial H3 context exclusively from
// persisted H2 Incident/LogEvent data and non-secret Project metadata.
package evidenceassembly

import (
	"context"
	"sort"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

const (
	maximumInitialEvidence = 17
	maximumCorpusBytes     = 128 << 10
)

type incidentEvidenceReader interface {
	portsout.IncidentGetter
	portsout.IncidentLogEventLister
}

type Assembler struct {
	incidents incidentEvidenceReader
	projects  portsout.ProjectStore
	accounts  portsout.GitHubAccountReader
	newID     func() uuid.UUID
	now       func() time.Time
}

func New(incidents incidentEvidenceReader, projects portsout.ProjectStore, accounts portsout.GitHubAccountReader, newID func() uuid.UUID, now func() time.Time) portsout.InvestigationContextAssembler {
	return &Assembler{incidents: incidents, projects: projects, accounts: accounts, newID: newID, now: now}
}

func (a *Assembler) Assemble(ctx context.Context, investigation domain.Investigation) (portsout.InvestigationRunContext, error) {
	incident, err := a.incidents.Get(ctx, investigation.IncidentID)
	if err != nil {
		return portsout.InvestigationRunContext{}, err
	}
	project, err := a.projects.Get(ctx, incident.ProjectID)
	if err != nil {
		return portsout.InvestigationRunContext{}, err
	}
	account, err := a.accounts.Get(ctx, project.GitHubRepository.GitHubAccountID)
	if err != nil {
		return portsout.InvestigationRunContext{}, err
	}

	page, err := a.incidents.ListLogEvents(ctx, incident.ID, paging.Params{
		Limit: maximumInitialEvidence,
		Sort:  []ukerpagination.SortExpression{{Field: "timestamp", Direction: ukerpagination.DirectionDesc}},
	})
	if err != nil {
		return portsout.InvestigationRunContext{}, err
	}
	sort.SliceStable(page.Items, func(left, right int) bool {
		if page.Items[left].Timestamp.Equal(page.Items[right].Timestamp) {
			return page.Items[left].ID.String() < page.Items[right].ID.String()
		}
		return page.Items[left].Timestamp.After(page.Items[right].Timestamp)
	})

	evidence, err := a.buildEvidence(investigation.ID, *project, page.Items)
	if err != nil {
		return portsout.InvestigationRunContext{}, err
	}
	return portsout.InvestigationRunContext{
		Investigation: investigation,
		Incident:      *incident,
		Project:       *project,
		Repository: portsout.RepositoryScope{
			Account: *account,
			Owner:   project.GitHubRepository.Owner,
			Name:    project.GitHubRepository.Name,
			Branch:  project.GitHubRepository.DefaultBranch,
		},
		Evidence: evidence,
	}, nil
}

func (a *Assembler) buildEvidence(investigationID uuid.UUID, project domain.Project, events []domain.LogEvent) ([]domain.Evidence, error) {
	now := a.now().UTC()
	metadata, err := deploymentMetadataEvidence(a.newID(), investigationID, project, now)
	if err != nil {
		return nil, err
	}
	candidates := []domain.Evidence{*metadata}
	stackAdded := false
	for _, event := range events {
		logEvidence, buildErr := logExcerptEvidence(a.newID(), investigationID, event, now)
		if buildErr != nil {
			return nil, buildErr
		}
		candidates = append(candidates, *logEvidence)
		if !stackAdded && hasStackTrace(event) {
			stack, stackErr := stackTraceEvidence(a.newID(), investigationID, event, now)
			if stackErr != nil {
				return nil, stackErr
			}
			candidates = append(candidates, *stack)
			stackAdded = true
		}
	}

	selected := make([]domain.Evidence, 0, maximumInitialEvidence)
	used := 0
	for _, candidate := range candidates {
		if len(selected) >= maximumInitialEvidence {
			break
		}
		size := evidenceSize(candidate)
		if used+size > maximumCorpusBytes {
			continue
		}
		selected = append(selected, candidate)
		used += size
	}
	return selected, nil
}

func evidenceSize(evidence domain.Evidence) int {
	return len(evidence.Summary) + len(evidence.Content) + len(evidence.FilePath) + len(evidence.CommitSHA) + len(evidence.Patch)
}
