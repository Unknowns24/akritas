package qvac

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
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
		"evidence_ids":[],
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

func TestParseInvestigationResultAcceptsConfidenceBoundariesAndKnownEvidence(t *testing.T) {
	t.Parallel()
	id := uuid.New()
	for _, confidence := range []float64{0, 1} {
		raw := fmt.Sprintf(`{"summary":"bounded","root_cause":"","root_cause_status":"unknown","resolution_status":"requires_human","confidence":%v,"hypotheses":[],"evidence_ids":[%q],"relevant_files":[],"relevant_commits":[],"recommended_actions":[]}`, confidence, id.String())
		result, err := parseInvestigationResult(raw, map[uuid.UUID]struct{}{id: {}})
		if err != nil || len(result.EvidenceIDs) != 1 {
			t.Fatalf("confidence %v: result=%+v err=%v", confidence, result, err)
		}
	}
}

func TestParseInvestigationResultAcceptsSuspectedRequiresHuman(t *testing.T) {
	t.Parallel()
	raw := `{"summary":"likely race","root_cause":"shared map race","root_cause_status":"suspected","resolution_status":"requires_human","confidence":0.55,"hypotheses":["race"],"evidence_ids":[],"relevant_files":[],"relevant_commits":[],"recommended_actions":["reproduce with race detector"]}`
	result, err := parseInvestigationResult(raw)
	if err != nil || result.RootCauseStatus != domain.RootCauseSuspected || result.ResolutionStatus != domain.ResolutionRequiresHuman {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}

func TestParseInvestigationResultRejectsMissingExtraAndUnknownEvidence(t *testing.T) {
	t.Parallel()
	unknown := uuid.New()
	cases := []string{
		`{"summary":"x"}`,
		`{"summary":"x","root_cause":"","root_cause_status":"unknown","resolution_status":"requires_human","confidence":0,"hypotheses":[],"evidence_ids":[],"relevant_files":[],"relevant_commits":[],"recommended_actions":[],"extra":true}`,
		fmt.Sprintf(`{"summary":"x","root_cause":"","root_cause_status":"unknown","resolution_status":"requires_human","confidence":0,"hypotheses":[],"evidence_ids":[%q],"relevant_files":[],"relevant_commits":[],"recommended_actions":[]}`, unknown.String()),
	}
	for _, raw := range cases {
		if _, err := parseInvestigationResult(raw); !errors.Is(err, ErrInvalidModelOutput) {
			t.Fatalf("expected strict rejection for %s, got %v", raw, err)
		}
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
		"evidence_ids":[],
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
		"evidence_ids":[],
		"relevant_files":[],
		"relevant_commits":[],
		"recommended_actions":[]
	}`
	_, err := parseInvestigationResult(raw)
	if !errors.Is(err, ErrInvalidModelOutput) {
		t.Fatalf("expected ErrInvalidModelOutput, got %v", err)
	}
}
