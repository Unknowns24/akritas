// Package evidenceassembly builds the initial H3 context exclusively from
// persisted H2 Incident/LogEvent data and non-secret Project metadata.
package evidenceassembly

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/Unknowns24/akritas/backend/internal/service/commitcorrelation"
	"github.com/Unknowns24/akritas/backend/internal/service/evidencesafety"
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
	incidents    incidentEvidenceReader
	projects     portsout.ProjectStore
	accounts     portsout.GitHubAccountReader
	commitReader commitReader
	newID        func() uuid.UUID
	now          func() time.Time
}

type commitReader interface {
	ListRecentCommits(context.Context, domain.GitHubAccount, string, string, string, int) ([]portsout.RepositoryCommitSummary, error)
}

func New(incidents incidentEvidenceReader, projects portsout.ProjectStore, accounts portsout.GitHubAccountReader, newID func() uuid.UUID, now func() time.Time) portsout.InvestigationContextAssembler {
	return &Assembler{incidents: incidents, projects: projects, accounts: accounts, newID: newID, now: now}
}

func NewWithCommitCorrelation(incidents incidentEvidenceReader, projects portsout.ProjectStore, accounts portsout.GitHubAccountReader, commitReader commitReader, newID func() uuid.UUID, now func() time.Time) portsout.InvestigationContextAssembler {
	return &Assembler{incidents: incidents, projects: projects, accounts: accounts, commitReader: commitReader, newID: newID, now: now}
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
	evidence = a.appendCommitEvidence(ctx, evidence, investigation.ID, *incident, *project, *account)
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

func (a *Assembler) appendCommitEvidence(ctx context.Context, current []domain.Evidence, investigationID uuid.UUID, incident domain.Incident, project domain.Project, account domain.GitHubAccount) []domain.Evidence {
	if a.commitReader == nil || len(current) >= maximumInitialEvidence {
		return current
	}
	commits, err := a.commitReader.ListRecentCommits(ctx, account, project.GitHubRepository.Owner, project.GitHubRepository.Name, project.GitHubRepository.DefaultBranch, 20)
	if err != nil {
		return current
	}
	correlated := commitcorrelation.Select(commitcorrelation.IncidentWindow{
		FirstSeenAt: incident.FirstSeenAt,
		LastSeenAt:  incident.LastSeenAt,
	}, commits, 3)
	for _, commit := range correlated {
		if len(current) >= maximumInitialEvidence {
			break
		}
		payload := struct {
			SHA       string `json:"sha"`
			Timestamp string `json:"timestamp,omitempty"`
			Author    string `json:"author,omitempty"`
			Message   string `json:"message,omitempty"`
			URL       string `json:"url,omitempty"`
			Reason    string `json:"reason"`
		}{
			SHA: commit.SHA, Author: commit.Author, Message: commit.Message, URL: commit.URL, Reason: commit.Reason,
		}
		if !commit.Timestamp.IsZero() {
			payload.Timestamp = commit.Timestamp.Format(time.RFC3339)
		}
		raw, _ := json.Marshal(payload)
		content := evidencesafety.RedactAndLimit(string(raw), 8000)
		evidence, buildErr := domain.NewEvidence(a.newID(), investigationID, domain.EvidenceCommit, "Commit potencialmente relacionado por ventana temporal; no es causa confirmada.", content, a.now().UTC())
		if buildErr != nil {
			continue
		}
		evidence.CommitSHA = commit.SHA
		current = append(current, *evidence)
	}
	return current
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
