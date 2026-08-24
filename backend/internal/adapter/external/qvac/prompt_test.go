package qvac

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func TestUserPromptContainsActualBoundedRedactedEvidence(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	id := uuid.New()
	evidence, err := domain.NewEvidence(id, uuid.New(), domain.EvidenceLogExcerpt, "database refused", strings.Repeat("context ", 1000)+" token=super-secret", now)
	if err != nil {
		t.Fatal(err)
	}
	runContext := portsout.InvestigationRunContext{
		Investigation: domain.Investigation{ID: evidence.InvestigationID},
		Incident:      domain.Incident{ID: uuid.New(), Title: "database connection refused", Severity: domain.SeverityError},
		Project:       domain.Project{Name: "api"}, Repository: portsout.RepositoryScope{Owner: "acme", Name: "api", Branch: "main"},
		Evidence: []domain.Evidence{*evidence},
	}
	prompt := userPrompt(runContext, maximumInitialPromptBytes)
	if !strings.Contains(prompt, "UNTRUSTED_DATA_BEGIN") || !strings.Contains(prompt, id.String()) || !strings.Contains(prompt, "database refused") {
		t.Fatalf("prompt does not carry actual framed Evidence: %s", prompt)
	}
	if strings.Contains(prompt, "super-secret") {
		t.Fatal("prompt leaked a credential-like value")
	}
	if budget := initialEvidenceBudget(10000); budget >= maximumInitialPromptBytes || budget != 5424 {
		t.Fatalf("dynamic context budget=%d", budget)
	}
	if initialEvidenceBudget(defaultContextSize) != maximumInitialPromptBytes {
		t.Fatal("default context must cap initial Evidence at 24 KiB")
	}
}

func TestBuildUserPromptReportsOnlyEvidenceActuallyShown(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	shown, err := domain.NewEvidence(uuid.New(), uuid.New(), domain.EvidenceLogExcerpt, "trigger", "database refused", now)
	if err != nil {
		t.Fatal(err)
	}
	omitted, err := domain.NewEvidence(uuid.New(), shown.InvestigationID, domain.EvidenceLogExcerpt, "later", strings.Repeat("x", 512), now)
	if err != nil {
		t.Fatal(err)
	}

	prompt, shownIDs := buildUserPrompt(portsout.InvestigationRunContext{Evidence: []domain.Evidence{*shown, *omitted}}, 180)
	if _, ok := shownIDs[shown.ID]; !ok || !strings.Contains(prompt, shown.ID.String()) {
		t.Fatalf("first Evidence was not shown: ids=%v prompt=%s", shownIDs, prompt)
	}
	if _, ok := shownIDs[omitted.ID]; ok || strings.Contains(prompt, omitted.ID.String()) {
		t.Fatalf("omitted Evidence was incorrectly marked citable: ids=%v prompt=%s", shownIDs, prompt)
	}
}

func TestToolDataEnvelopeStaysValidJSONWithinPayloadLimit(t *testing.T) {
	t.Parallel()
	evidenceID := uuid.New()
	envelope := toolDataEnvelope(strings.Repeat(`\\\"`, maximumToolPayloadBytes), evidenceID, maximumToolPayloadBytes)
	if len(envelope) > maximumToolPayloadBytes {
		t.Fatalf("tool envelope has %d bytes", len(envelope))
	}
	var decoded struct {
		EvidenceID string `json:"evidence_id"`
		Data       any    `json:"data"`
	}
	if err := json.Unmarshal([]byte(envelope), &decoded); err != nil {
		t.Fatalf("tool envelope is not valid JSON: %v", err)
	}
	if decoded.EvidenceID != evidenceID.String() {
		t.Fatalf("evidence_id=%q", decoded.EvidenceID)
	}
}
