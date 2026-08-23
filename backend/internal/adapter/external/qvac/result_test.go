package qvac

import (
	"errors"
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TestParseInvestigationResultAcceptsValidPayload(t *testing.T) {
	t.Parallel()
	raw := `{
		"summary":"nil deref in handler",
		"root_cause":"missing nil check",
		"root_cause_status":"identified",
		"resolution_status":"fixable",
		"confidence":0.81,
		"hypotheses":["h1"],
		"relevant_files":["main.go"],
		"relevant_commits":["abc"],
		"recommended_actions":["add guard"]
	}`
	result, err := parseInvestigationResult(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Summary != "nil deref in handler" || result.RootCauseStatus != domain.RootCauseIdentified ||
		result.ResolutionStatus != domain.ResolutionFixable || result.Confidence != 0.81 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestParseInvestigationResultRejectsUnknownStatus(t *testing.T) {
	t.Parallel()
	raw := `{
		"summary":"x",
		"root_cause":"y",
		"root_cause_status":"maybe",
		"resolution_status":"fixable",
		"confidence":0.5,
		"hypotheses":[],
		"relevant_files":[],
		"relevant_commits":[],
		"recommended_actions":[]
	}`
	_, err := parseInvestigationResult(raw)
	if !errors.Is(err, ErrInvalidModelOutput) {
		t.Fatalf("expected ErrInvalidModelOutput, got %v", err)
	}
}

func TestParseInvestigationResultRejectsBadConfidence(t *testing.T) {
	t.Parallel()
	raw := `{
		"summary":"x",
		"root_cause":"y",
		"root_cause_status":"unknown",
		"resolution_status":"requires_human",
		"confidence":1.5,
		"hypotheses":[],
		"relevant_files":[],
		"relevant_commits":[],
		"recommended_actions":[]
	}`
	_, err := parseInvestigationResult(raw)
	if !errors.Is(err, ErrInvalidModelOutput) {
		t.Fatalf("expected ErrInvalidModelOutput, got %v", err)
	}
}
