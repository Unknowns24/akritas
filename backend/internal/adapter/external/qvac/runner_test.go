package qvac

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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
	result, err := runner.Run(context.Background(), *investigation)
	if err != nil {
		t.Fatal(err)
	}
	if result.Summary != "crash in worker" || result.RelevantFiles[0] != "worker.go" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestRunnerInvalidOutputErrors(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	_, err = runner.Run(context.Background(), *investigation)
	if !errors.Is(err, ErrInvalidModelOutput) {
		t.Fatalf("expected ErrInvalidModelOutput, got %v", err)
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
	result, err := runner.Run(context.Background(), *investigation)
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

func TestRunnerUnknownToolReturnsControlledErrorPayload(t *testing.T) {
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
	if _, err := runner.Run(context.Background(), *investigation); err != nil {
		t.Fatal(err)
	}
}

var _ out.InvestigationRunner = (*Runner)(nil)
