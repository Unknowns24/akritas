package qvac

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type modelResultDTO struct {
	Summary            string   `json:"summary"`
	RootCause          string   `json:"root_cause"`
	RootCauseStatus    string   `json:"root_cause_status"`
	ResolutionStatus   string   `json:"resolution_status"`
	Confidence         float64  `json:"confidence"`
	Hypotheses         []string `json:"hypotheses"`
	RelevantFiles      []string `json:"relevant_files"`
	RelevantCommits    []string `json:"relevant_commits"`
	RecommendedActions []string `json:"recommended_actions"`
}

var investigationResultSchema = json.RawMessage(`{
  "type": "object",
  "additionalProperties": false,
  "required": ["summary", "root_cause", "root_cause_status", "resolution_status", "confidence", "hypotheses", "relevant_files", "relevant_commits", "recommended_actions"],
  "properties": {
    "summary": {"type": "string"},
    "root_cause": {"type": "string"},
    "root_cause_status": {"type": "string", "enum": ["identified", "suspected", "unknown"]},
    "resolution_status": {"type": "string", "enum": ["fixable", "requires_human"]},
    "confidence": {"type": "number"},
    "hypotheses": {"type": "array", "items": {"type": "string"}},
    "relevant_files": {"type": "array", "items": {"type": "string"}},
    "relevant_commits": {"type": "array", "items": {"type": "string"}},
    "recommended_actions": {"type": "array", "items": {"type": "string"}}
  }
}`)

func parseInvestigationResult(raw string) (out.InvestigationRunResult, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return out.InvestigationRunResult{}, fmt.Errorf("%w: empty content", ErrInvalidModelOutput)
	}
	var dto modelResultDTO
	if err := json.Unmarshal([]byte(raw), &dto); err != nil {
		return out.InvestigationRunResult{}, fmt.Errorf("%w: %v", ErrInvalidModelOutput, err)
	}
	if strings.TrimSpace(dto.Summary) == "" {
		return out.InvestigationRunResult{}, fmt.Errorf("%w: summary is required", ErrInvalidModelOutput)
	}
	rootCauseStatus := domain.RootCauseStatus(strings.TrimSpace(dto.RootCauseStatus))
	if err := rootCauseStatus.Validate(); err != nil {
		return out.InvestigationRunResult{}, fmt.Errorf("%w: root_cause_status", ErrInvalidModelOutput)
	}
	resolutionStatus := domain.ResolutionStatus(strings.TrimSpace(dto.ResolutionStatus))
	if err := resolutionStatus.Validate(); err != nil {
		return out.InvestigationRunResult{}, fmt.Errorf("%w: resolution_status", ErrInvalidModelOutput)
	}
	if dto.Confidence < 0 || dto.Confidence > 1 {
		return out.InvestigationRunResult{}, fmt.Errorf("%w: confidence", ErrInvalidModelOutput)
	}
	return out.InvestigationRunResult{
		Summary:            strings.TrimSpace(dto.Summary),
		RootCause:          strings.TrimSpace(dto.RootCause),
		RootCauseStatus:    rootCauseStatus,
		ResolutionStatus:   resolutionStatus,
		Confidence:         dto.Confidence,
		Hypotheses:         cloneStrings(dto.Hypotheses),
		RelevantFiles:      cloneStrings(dto.RelevantFiles),
		RelevantCommits:    cloneStrings(dto.RelevantCommits),
		RecommendedActions: cloneStrings(dto.RecommendedActions),
	}, nil
}

func cloneStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	cloned := make([]string, len(values))
	copy(cloned, values)
	return cloned
}
