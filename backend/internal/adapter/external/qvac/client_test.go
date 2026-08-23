package qvac

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClientRejectsPublicEndpoint(t *testing.T) {
	t.Parallel()
	_, err := NewClient(ClientConfig{EndpointURL: "https://api.openai.com/v1"})
	if !errors.Is(err, ErrInvalidEndpoint) {
		t.Fatalf("expected ErrInvalidEndpoint, got %v", err)
	}
}

func TestNewClientAllowsDockerHostEndpoint(t *testing.T) {
	t.Parallel()
	client, err := NewClient(ClientConfig{EndpointURL: "http://host.docker.internal:11434/v1"})
	if err != nil {
		t.Fatal(err)
	}
	if client.Endpoint() != "http://host.docker.internal:11434/v1" {
		t.Fatalf("endpoint = %q", client.Endpoint())
	}
}

func TestChatCompletionsHappyPath(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": `{"ok":true}`}}},
		})
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.chatCompletions(context.Background(), chatRequest{
		Messages: []chatMessage{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if response.Choices[0].Message.Content != `{"ok":true}` {
		t.Fatalf("content = %q", response.Choices[0].Message.Content)
	}
}

func TestChatCompletionsSendsConfiguredContextSize(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request chatRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request.Options == nil || request.Options.NumCtx != 65536 {
			t.Fatalf("num_ctx = %+v", request.Options)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": `{"ok":true}`}}},
		})
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client(), ContextSize: 65536})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.chatCompletions(context.Background(), chatRequest{
		Messages: []chatMessage{{Role: "user", Content: "hi"}},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestChatCompletionsMapsModelMissing(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"message":"model not found"}}`))
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.chatCompletions(context.Background(), chatRequest{Messages: []chatMessage{{Role: "user", Content: "hi"}}})
	if err == nil || !errors.Is(err, ErrModelUnavailable) {
		t.Fatalf("expected ErrModelUnavailable, got %v", err)
	}
}
