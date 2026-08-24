package qvac

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

type modelResultDTO struct {
	Summary            *string   `json:"summary"`
	RootCause          *string   `json:"root_cause"`
	RootCauseStatus    *string   `json:"root_cause_status"`
	ResolutionStatus   *string   `json:"resolution_status"`
	Confidence         *float64  `json:"confidence"`
	Hypotheses         *[]string `json:"hypotheses"`
	EvidenceIDs        *[]string `json:"evidence_ids"`
	RelevantFiles      *[]string `json:"relevant_files"`
	RelevantCommits    *[]string `json:"relevant_commits"`
	RecommendedActions *[]string `json:"recommended_actions"`
}

var investigationResultSchema = json.RawMessage(`{
  "type": "object",
  "additionalProperties": false,
  "required": ["summary", "root_cause", "root_cause_status", "resolution_status", "confidence", "hypotheses", "evidence_ids", "relevant_files", "relevant_commits", "recommended_actions"],
  "properties": {
    "summary": {"type": "string", "minLength": 1, "maxLength": 10000},
    "root_cause": {"type": "string", "maxLength": 20000},
    "root_cause_status": {"type": "string", "enum": ["identified", "suspected", "unknown"]},
    "resolution_status": {"type": "string", "enum": ["fixable", "requires_human"]},
    "confidence": {"type": "number", "minimum": 0, "maximum": 1},
    "hypotheses": {"type": "array", "maxItems": 50, "items": {"type": "string", "minLength": 1, "maxLength": 5000}},
    "evidence_ids": {"type": "array", "maxItems": 25, "uniqueItems": true, "items": {"type": "string", "format": "uuid"}},
    "relevant_files": {"type": "array", "maxItems": 100, "items": {"type": "string", "minLength": 1, "maxLength": 4096}},
    "relevant_commits": {"type": "array", "maxItems": 100, "items": {"type": "string", "minLength": 1, "maxLength": 64}},
    "recommended_actions": {"type": "array", "maxItems": 50, "items": {"type": "string", "minLength": 1, "maxLength": 5000}}
  }
}`)

func parseInvestigationResult(raw string, allowedGroups ...map[uuid.UUID]struct{}) (portsout.InvestigationRunResult, error) {
	allowedEvidence := map[uuid.UUID]struct{}{}
	if len(allowedGroups) > 0 && allowedGroups[0] != nil {
		allowedEvidence = allowedGroups[0]
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return portsout.InvestigationRunResult{}, invalidOutput("empty content")
	}
	raw = extractJSONObject(raw)
	decoder := json.NewDecoder(bytes.NewBufferString(raw))
	decoder.DisallowUnknownFields()
	var dto modelResultDTO
	if err := decoder.Decode(&dto); err != nil {
		return portsout.InvestigationRunResult{}, invalidOutput(err.Error())
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return portsout.InvestigationRunResult{}, invalidOutput("multiple JSON values")
	}
	if dto.Summary == nil || dto.RootCause == nil || dto.RootCauseStatus == nil || dto.ResolutionStatus == nil || dto.Confidence == nil ||
		dto.Hypotheses == nil || dto.EvidenceIDs == nil || dto.RelevantFiles == nil || dto.RelevantCommits == nil || dto.RecommendedActions == nil {
		return portsout.InvestigationRunResult{}, invalidOutput("required field missing")
	}
	if strings.TrimSpace(*dto.Summary) == "" || len(*dto.Summary) > 10000 || len(*dto.RootCause) > 20000 {
		return portsout.InvestigationRunResult{}, invalidOutput("summary or root_cause")
	}
	rootCauseStatus := domain.RootCauseStatus(strings.TrimSpace(*dto.RootCauseStatus))
	if err := rootCauseStatus.Validate(); err != nil {
		return portsout.InvestigationRunResult{}, invalidOutput("root_cause_status")
	}
	resolutionStatus := domain.ResolutionStatus(strings.TrimSpace(*dto.ResolutionStatus))
	if err := resolutionStatus.Validate(); err != nil {
		return portsout.InvestigationRunResult{}, invalidOutput("resolution_status")
	}
	if *dto.Confidence < 0 || *dto.Confidence > 1 {
		return portsout.InvestigationRunResult{}, invalidOutput("confidence")
	}
	if err := validateStringArray(*dto.Hypotheses, 50, 5000); err != nil {
		return portsout.InvestigationRunResult{}, invalidOutput("hypotheses")
	}
	if err := validateStringArray(*dto.RelevantFiles, 100, 4096); err != nil {
		return portsout.InvestigationRunResult{}, invalidOutput("relevant_files")
	}
	if err := validateStringArray(*dto.RelevantCommits, 100, 64); err != nil {
		return portsout.InvestigationRunResult{}, invalidOutput("relevant_commits")
	}
	if err := validateStringArray(*dto.RecommendedActions, 50, 5000); err != nil {
		return portsout.InvestigationRunResult{}, invalidOutput("recommended_actions")
	}
	if len(*dto.EvidenceIDs) > 25 {
		return portsout.InvestigationRunResult{}, invalidOutput("evidence_ids")
	}
	evidenceIDs := make([]uuid.UUID, 0, len(*dto.EvidenceIDs))
	seen := make(map[uuid.UUID]struct{}, len(*dto.EvidenceIDs))
	for _, rawID := range *dto.EvidenceIDs {
		id, err := uuid.Parse(rawID)
		if err != nil {
			return portsout.InvestigationRunResult{}, invalidOutput("evidence_ids")
		}
		if _, duplicate := seen[id]; duplicate {
			return portsout.InvestigationRunResult{}, invalidOutput("duplicate evidence_id")
		}
		if _, shown := allowedEvidence[id]; !shown {
			return portsout.InvestigationRunResult{}, invalidOutput("unknown evidence_id")
		}
		seen[id] = struct{}{}
		evidenceIDs = append(evidenceIDs, id)
	}
	return portsout.InvestigationRunResult{
		Summary: strings.TrimSpace(*dto.Summary), RootCause: strings.TrimSpace(*dto.RootCause),
		RootCauseStatus: rootCauseStatus, ResolutionStatus: resolutionStatus, Confidence: *dto.Confidence,
		Hypotheses: cloneStrings(*dto.Hypotheses), EvidenceIDs: evidenceIDs,
		RelevantFiles: cloneStrings(*dto.RelevantFiles), RelevantCommits: cloneStrings(*dto.RelevantCommits),
		RecommendedActions: cloneStrings(*dto.RecommendedActions),
	}, nil
}

func validateStringArray(values []string, maximumItems, maximumLength int) error {
	if len(values) > maximumItems {
		return fmt.Errorf("too many values")
	}
	for _, value := range values {
		if strings.TrimSpace(value) == "" || len(value) > maximumLength {
			return fmt.Errorf("invalid value")
		}
	}
	return nil
}

func invalidOutput(detail string) error {
	return fmt.Errorf("%w: %s", ErrInvalidModelOutput, detail)
}

func cloneStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return append([]string(nil), values...)
}

func extractJSONObject(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "{") {
		if start, end, ok := firstBalancedJSONObject(raw); ok && start == 0 && strings.TrimSpace(raw[end:]) != "" {
			return raw[:end]
		}
		return raw
	}
	if start := strings.Index(raw, "```"); start >= 0 {
		fenced := raw[start+len("```"):]
		if newline := strings.IndexByte(fenced, '\n'); newline >= 0 {
			fenced = fenced[newline+1:]
		}
		if end := strings.Index(fenced, "```"); end >= 0 {
			raw = strings.TrimSpace(fenced[:end])
			if strings.HasPrefix(raw, "{") {
				return raw
			}
		}
	}
	if start, end, ok := firstBalancedJSONObject(raw); ok {
		return raw[start:end]
	}
	return raw
}

func firstBalancedJSONObject(raw string) (int, int, bool) {
	start := strings.IndexByte(raw, '{')
	if start < 0 {
		return 0, 0, false
	}
	depth := 0
	inString := false
	escaped := false
	for index := start; index < len(raw); index++ {
		ch := raw[index]
		if inString {
			if escaped {
				escaped = false
				continue
			}
			switch ch {
			case '\\':
				escaped = true
			case '"':
				inString = false
			}
			continue
		}
		switch ch {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return start, index + 1, true
			}
		}
	}
	return 0, 0, false
}
