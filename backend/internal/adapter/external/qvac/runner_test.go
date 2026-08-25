package qvac

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func validResultJSON() string {
	return `{
		"summary":"crash in worker",
		"root_cause":"nil map",
		"root_cause_status":"identified",
		"resolution_status":"fixable",
		"confidence":0.9,
		"hypotheses":["h"],
		"evidence_ids":[],
		"relevant_files":["worker.go"],
		"relevant_commits":["deadbeef"],
		"recommended_actions":["initialize map"]
	}`
}

func TestRunnerCompletesWithoutTools(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": validResultJSON()}}},
		})
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewRunner(client, nil, RunnerConfig{})
	if err != nil {
		t.Fatal(err)
	}
	investigation, err := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	result, err := runner.Run(context.Background(), out.InvestigationRunContext{Investigation: *investigation})
	if err != nil {
		t.Fatal(err)
	}
	if result.Summary != "crash in worker" || result.RelevantFiles[0] != "worker.go" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestRunnerReturnsHumanReviewResultWhenFinalAndRepairOutputAreInvalid(t *testing.T) {
	t.Parallel()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		assertNoResponseFormat(t, r)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": `{"summary":""}`}}},
		})
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewRunner(client, nil, RunnerConfig{})
	if err != nil {
		t.Fatal(err)
	}
	investigation, err := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	result, err := runner.Run(context.Background(), out.InvestigationRunContext{Investigation: *investigation})
	if err != nil {
		t.Fatal(err)
	}
	if result.RootCauseStatus != domain.RootCauseUnknown || result.ResolutionStatus != domain.ResolutionRequiresHuman {
		t.Fatalf("expected human-review result, got %+v", result)
	}
	if calls.Load() != 2 {
		t.Fatalf("expected final attempt plus repair attempt, got %d", calls.Load())
	}
}

func TestRunnerSendsFinalResultWithoutResponseFormat(t *testing.T) {
	t.Parallel()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		assertNoResponseFormat(t, r)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": validResultJSON()}}},
		})
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewRunner(client, nil, RunnerConfig{})
	if err != nil {
		t.Fatal(err)
	}
	investigation, _ := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	result, err := runner.Run(context.Background(), out.InvestigationRunContext{Investigation: *investigation})
	if err != nil {
		t.Fatal(err)
	}
	if result.RootCause != "nil map" || calls.Load() != 1 {
		t.Fatalf("unexpected final result=%+v calls=%d", result, calls.Load())
	}
}

func TestRunnerRepairsInvalidFinalResultWithoutResponseFormat(t *testing.T) {
	t.Parallel()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertNoResponseFormat(t, r)
		if calls.Add(1) == 1 {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": "summary: crash in worker\nroot cause: nil map"}}},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": validResultJSON()}}},
		})
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewRunner(client, nil, RunnerConfig{})
	if err != nil {
		t.Fatal(err)
	}
	investigation, _ := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	result, err := runner.Run(context.Background(), out.InvestigationRunContext{Investigation: *investigation})
	if err != nil {
		t.Fatal(err)
	}
	if result.RootCause != "nil map" || calls.Load() != 2 {
		t.Fatalf("unexpected repaired result=%+v calls=%d", result, calls.Load())
	}
}

func TestRunnerFailsExplicitlyWhenFinalSynthesisUnavailable(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"message": "runtime failed"}})
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewRunner(client, nil, RunnerConfig{})
	if err != nil {
		t.Fatal(err)
	}
	investigation, _ := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	evidenceID := uuid.New()
	evidence, _ := domain.NewEvidence(evidenceID, investigation.ID, domain.EvidenceLogExcerpt, "database error", "failed to initialize database", time.Now().UTC())
	_, err = runner.Run(context.Background(), out.InvestigationRunContext{
		Investigation: *investigation,
		Evidence:      []domain.Evidence{*evidence},
	})
	if !errors.Is(err, domain.ErrQvacUnavailable) {
		t.Fatalf("expected explicit QVAC unavailable error, got %v", err)
	}
}

func assertNoResponseFormat(t *testing.T, r *http.Request) {
	t.Helper()
	var request map[string]any
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		t.Fatal(err)
	}
	if _, exists := request["response_format"]; exists {
		t.Fatalf("request unexpectedly included response_format: %+v", request["response_format"])
	}
	if request["tool_choice"] != "none" {
		t.Fatalf("tool_choice = %v, want none", request["tool_choice"])
	}
}

func TestRunnerFailsExplicitlyWhenToolRoundFailsAfterEvidence(t *testing.T) {
	t.Parallel()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"choices": []map[string]any{{"message": map[string]any{
					"role": "assistant", "tool_calls": []map[string]any{{"id": "call_file", "type": "function", "function": map[string]any{"name": "read_file", "arguments": `{"path":"internal/db.go"}`}}},
				}}},
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"message": "runtime failed"}})
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	evidenceID := uuid.New()
	runner, err := NewRunner(client, nil, RunnerConfig{RepositoryInspector: &fakeRepositoryInspector{}, NewID: func() uuid.UUID { return evidenceID }})
	if err != nil {
		t.Fatal(err)
	}
	investigation, _ := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	result, err := runner.Run(context.Background(), out.InvestigationRunContext{
		Investigation: *investigation,
		Repository:    out.RepositoryScope{Owner: "acme", Name: "api", Branch: "main"},
	})
	if !errors.Is(err, domain.ErrQvacUnavailable) {
		t.Fatalf("expected explicit QVAC unavailable error, got %v", err)
	}
	if len(result.DiscoveredEvidence) != 1 || result.DiscoveredEvidence[0].ID != evidenceID {
		t.Fatalf("expected discovered evidence to survive unavailable QVAC error, got %+v", result)
	}
}

func TestRunnerDoesNotPersistModelCitationForEvidenceOmittedFromBoundedPrompt(t *testing.T) {
	t.Parallel()
	evidenceID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := strings.Replace(validResultJSON(), `"evidence_ids":[]`, fmt.Sprintf(`"evidence_ids":[%q]`, evidenceID.String()), 1)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": payload}}},
		})
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewRunner(client, nil, RunnerConfig{ContextSize: reservedContextTokens})
	if err != nil {
		t.Fatal(err)
	}
	investigation, _ := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	evidence, _ := domain.NewEvidence(evidenceID, investigation.ID, domain.EvidenceLogExcerpt, "real but omitted", "content", time.Now().UTC())
	result, err := runner.Run(context.Background(), out.InvestigationRunContext{Investigation: *investigation, Evidence: []domain.Evidence{*evidence}})
	if err != nil {
		t.Fatal(err)
	}
	if result.RootCauseStatus != domain.RootCauseUnknown || result.ResolutionStatus != domain.ResolutionRequiresHuman {
		t.Fatalf("expected human-review result, got %+v", result)
	}
}

type echoTool struct{}

func (echoTool) Name() string        { return "echo" }
func (echoTool) Description() string { return "echo args" }
func (echoTool) Parameters() json.RawMessage {
	return json.RawMessage(`{"type":"object","properties":{"text":{"type":"string"}}}`)
}
func (echoTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	return string(arguments), nil
}

type largeOutputTool struct{}

func (largeOutputTool) Name() string        { return "large_output" }
func (largeOutputTool) Description() string { return "returns a large read-only payload" }
func (largeOutputTool) Parameters() json.RawMessage {
	return json.RawMessage(`{"type":"object","properties":{"index":{"type":"integer"}}}`)
}
func (largeOutputTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	return fmt.Sprintf(`{"content":%q}`, strings.Repeat("tool-data ", 2048)), nil
}

func TestRunnerToolLoopThenStructuredResult(t *testing.T) {
	t.Parallel()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if n == 1 {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"choices": []map[string]any{{
					"message": map[string]any{
						"role": "assistant",
						"tool_calls": []map[string]any{{
							"id": "call_1", "type": "function",
							"function": map[string]any{"name": "echo", "arguments": `{"text":"hi"}`},
						}},
					},
				}},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": validResultJSON()}}},
		})
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewRunner(client, NewToolRegistry(echoTool{}), RunnerConfig{})
	if err != nil {
		t.Fatal(err)
	}
	investigation, err := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	result, err := runner.Run(context.Background(), out.InvestigationRunContext{Investigation: *investigation})
	if err != nil {
		t.Fatal(err)
	}
	if result.RootCause != "nil map" {
		t.Fatalf("unexpected result: %+v", result)
	}
	if calls.Load() < 2 {
		t.Fatalf("expected tool round + final round, got %d", calls.Load())
	}
}

func TestRunnerStopsToolExplorationBeforeDefaultContextOverflow(t *testing.T) {
	t.Parallel()
	var toolRequests atomic.Int32
	var finalRequests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request chatRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request.ToolChoice == "none" {
			finalRequests.Add(1)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": validResultJSON()}}},
			})
			return
		}
		if len(request.Messages) > 8 {
			t.Fatalf("runner sent an overflow-prone tool round with %d messages", len(request.Messages))
		}
		n := toolRequests.Add(1)
		_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]any{
			"role": "assistant", "tool_calls": []map[string]any{{
				"id": fmt.Sprintf("call_%d", n), "type": "function",
				"function": map[string]any{"name": "large_output", "arguments": fmt.Sprintf(`{"index":%d}`, n)},
			}},
		}}}})
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewRunner(client, NewToolRegistry(largeOutputTool{}), RunnerConfig{MaxToolRounds: 10, MaxToolCalls: 20})
	if err != nil {
		t.Fatal(err)
	}
	result, err := runner.Run(context.Background(), out.InvestigationRunContext{})
	if err != nil {
		t.Fatalf("expected final synthesis after tool budget exhaustion, got %v", err)
	}
	if result.RootCause != "nil map" {
		t.Fatalf("unexpected result after prompt budget exhaustion: %+v", result)
	}
	if toolRequests.Load() == 0 || finalRequests.Load() != 1 {
		t.Fatalf("expected tool exploration followed by one final synthesis, tools=%d final=%d", toolRequests.Load(), finalRequests.Load())
	}
}

func TestRunnerTurnsRepositoryToolDataIntoCitableEvidence(t *testing.T) {
	t.Parallel()
	evidenceID := uuid.New()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request chatRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		switch calls.Add(1) {
		case 1:
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]any{
				"role": "assistant", "tool_calls": []map[string]any{{"id": "call_file", "type": "function", "function": map[string]any{"name": "read_file", "arguments": `{"path":"internal/db.go"}`}}},
			}}}})
		case 2:
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": "ready"}}}})
		default:
			encoded, _ := json.Marshal(request.Messages)
			if !strings.Contains(string(encoded), evidenceID.String()) || !strings.Contains(string(encoded), "UNTRUSTED_DATA_BEGIN") {
				t.Fatalf("tool evidence ID/data was not framed for final QVAC prompt: %s", encoded)
			}
			payload := strings.Replace(validResultJSON(), `"evidence_ids":[]`, fmt.Sprintf(`"evidence_ids":[%q]`, evidenceID.String()), 1)
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": payload}}}})
		}
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	inspector := &fakeRepositoryInspector{}
	runner, err := NewRunner(client, nil, RunnerConfig{RepositoryInspector: inspector, NewID: func() uuid.UUID { return evidenceID }, Now: time.Now})
	if err != nil {
		t.Fatal(err)
	}
	investigation, _ := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	result, err := runner.Run(context.Background(), out.InvestigationRunContext{
		Investigation: *investigation, Repository: out.RepositoryScope{Owner: "acme", Name: "api", Branch: "main"},
	})
	if err != nil || len(result.DiscoveredEvidence) != 1 || len(result.EvidenceIDs) != 1 || result.EvidenceIDs[0] != evidenceID {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}

func TestRunnerDeduplicatesRepeatedRepositoryReadsAndEvidence(t *testing.T) {
	t.Parallel()
	evidenceID := uuid.New()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch calls.Add(1) {
		case 1:
			toolCall := map[string]any{"type": "function", "function": map[string]any{"name": "read_file", "arguments": `{"path":"internal/db.go"}`}}
			first := mapsClone(toolCall)
			first["id"] = "call_1"
			second := mapsClone(toolCall)
			second["id"] = "call_2"
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "tool_calls": []map[string]any{first, second}}}}})
		case 2:
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": "ready"}}}})
		default:
			payload := strings.Replace(validResultJSON(), `"evidence_ids":[]`, fmt.Sprintf(`"evidence_ids":[%q]`, evidenceID.String()), 1)
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": payload}}}})
		}
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	inspector := &fakeRepositoryInspector{}
	runner, err := NewRunner(client, nil, RunnerConfig{RepositoryInspector: inspector, NewID: func() uuid.UUID { return evidenceID }})
	if err != nil {
		t.Fatal(err)
	}
	investigation, _ := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	result, err := runner.Run(context.Background(), out.InvestigationRunContext{
		Investigation: *investigation, Repository: out.RepositoryScope{Owner: "acme", Name: "api", Branch: "main"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(inspector.calls) != 1 || len(result.DiscoveredEvidence) != 1 || len(result.EvidenceIDs) != 1 {
		t.Fatalf("repeated read was not deduplicated: calls=%v result=%+v", inspector.calls, result)
	}
}

func TestRunnerRedactsRepositoryToolOutputBeforeModelAndEvidence(t *testing.T) {
	t.Parallel()
	evidenceID := uuid.New()
	secret := "sentinel-tool-secret"
	var toolMessage string
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request chatRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		switch calls.Add(1) {
		case 1:
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]any{
				"role": "assistant", "tool_calls": []map[string]any{{"id": "call_file", "type": "function", "function": map[string]any{"name": "read_file", "arguments": `{"path":"internal/config.go"}`}}},
			}}}})
		case 2:
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": "ready"}}}})
		default:
			encoded, _ := json.Marshal(request.Messages)
			toolMessage = string(encoded)
			if strings.Contains(toolMessage, secret) {
				t.Fatalf("tool output leaked to QVAC: %s", toolMessage)
			}
			if !strings.Contains(toolMessage, "Authorization: Bearer [REDACTED]") || !strings.Contains(toolMessage, "UNTRUSTED_DATA_BEGIN") {
				t.Fatalf("tool output was not redacted/framed: %s", toolMessage)
			}
			payload := strings.Replace(validResultJSON(), `"evidence_ids":[]`, fmt.Sprintf(`"evidence_ids":[%q]`, evidenceID.String()), 1)
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": payload}}}})
		}
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	inspector := &fakeRepositoryInspector{fileContent: "Authorization: Bearer " + secret + "\npackage config"}
	runner, err := NewRunner(client, nil, RunnerConfig{RepositoryInspector: inspector, NewID: func() uuid.UUID { return evidenceID }})
	if err != nil {
		t.Fatal(err)
	}
	investigation, _ := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	result, err := runner.Run(context.Background(), out.InvestigationRunContext{
		Investigation: *investigation, Repository: out.RepositoryScope{Owner: "acme", Name: "api", Branch: "main"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.DiscoveredEvidence) != 1 || strings.Contains(result.DiscoveredEvidence[0].Content, secret) {
		t.Fatalf("discovered Evidence leaked tool secret: %+v", result.DiscoveredEvidence)
	}
	if toolMessage == "" {
		t.Fatal("test server did not capture the tool round")
	}
}

func mapsClone(value map[string]any) map[string]any {
	clone := make(map[string]any, len(value))
	for key, item := range value {
		clone[key] = item
	}
	return clone
}

func TestRunnerUnknownToolFailsClosed(t *testing.T) {
	t.Parallel()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if n == 1 {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"choices": []map[string]any{{
					"message": map[string]any{
						"role": "assistant",
						"tool_calls": []map[string]any{{
							"id": "call_x", "type": "function",
							"function": map[string]any{"name": "drop_database", "arguments": `{}`},
						}},
					},
				}},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": validResultJSON()}}},
		})
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewRunner(client, NewToolRegistry(echoTool{}), RunnerConfig{})
	if err != nil {
		t.Fatal(err)
	}
	investigation, err := domain.NewInvestigation(uuid.New(), uuid.New(), time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runner.Run(context.Background(), out.InvestigationRunContext{Investigation: *investigation}); !errors.Is(err, ErrUnknownTool) {
		t.Fatalf("expected ErrUnknownTool, got %v", err)
	}
}

func TestRunnerSynthesizesFinalResultWhenToolBudgetIsExhausted(t *testing.T) {
	t.Parallel()
	for _, testCase := range []struct {
		name   string
		config RunnerConfig
		calls  int
	}{
		{name: "rounds", config: RunnerConfig{MaxToolRounds: 1, MaxToolCalls: 24}, calls: 1},
		{name: "calls", config: RunnerConfig{MaxToolRounds: 8, MaxToolCalls: 1}, calls: 2},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if requests.Add(1) > 1 {
					_ = json.NewEncoder(w).Encode(map[string]any{
						"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": validResultJSON()}}},
					})
					return
				}
				toolCalls := make([]map[string]any, 0, testCase.calls)
				for index := 0; index < testCase.calls; index++ {
					toolCalls = append(toolCalls, map[string]any{
						"id": fmt.Sprintf("call_%d", index), "type": "function",
						"function": map[string]any{"name": "echo", "arguments": `{"text":"hi"}`},
					})
				}
				_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "tool_calls": toolCalls}}}})
			}))
			t.Cleanup(server.Close)
			client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
			if err != nil {
				t.Fatal(err)
			}
			runner, err := NewRunner(client, NewToolRegistry(echoTool{}), testCase.config)
			if err != nil {
				t.Fatal(err)
			}
			result, err := runner.Run(context.Background(), out.InvestigationRunContext{})
			if err != nil {
				t.Fatalf("expected final synthesis to complete, got %v", err)
			}
			if result.RootCause != "nil map" {
				t.Fatalf("unexpected result after budget exhaustion: %+v", result)
			}
		})
	}
}

var _ out.InvestigationRunner = (*Runner)(nil)
