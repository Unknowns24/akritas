package qvac

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

const (
	defaultEndpoint  = "http://127.0.0.1:11434/v1"
	defaultModel     = "akritas"
	maximumBodyBytes = 4 << 20
)

type ClientConfig struct {
	EndpointURL string
	APIKey      string
	HTTPClient  *http.Client
	Model       string
}

type Client struct {
	base       string
	apiKey     string
	httpClient *http.Client
	model      string
}

func NewClient(config ClientConfig) (*Client, error) {
	endpoint := strings.TrimSpace(config.EndpointURL)
	if endpoint == "" {
		endpoint = defaultEndpoint
	}
	base, err := validateEndpoint(endpoint)
	if err != nil {
		return nil, err
	}
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 2 * time.Minute}
	}
	model := strings.TrimSpace(config.Model)
	if model == "" {
		model = defaultModel
	}
	return &Client{
		base:       strings.TrimRight(base.String(), "/"),
		apiKey:     config.APIKey,
		httpClient: httpClient,
		model:      model,
	}, nil
}

func (c *Client) Model() string { return c.model }

type chatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []toolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

type toolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function toolCallFunction `json:"function"`
}

type toolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type toolDefinition struct {
	Type     string             `json:"type"`
	Function toolFunctionSchema `json:"function"`
}

type toolFunctionSchema struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

type responseFormat struct {
	Type       string          `json:"type"`
	JSONSchema *jsonSchemaSpec `json:"json_schema,omitempty"`
}

type jsonSchemaSpec struct {
	Name   string          `json:"name"`
	Strict bool            `json:"strict"`
	Schema json.RawMessage `json:"schema"`
}

type chatRequest struct {
	Model          string           `json:"model"`
	Messages       []chatMessage    `json:"messages"`
	Tools          []toolDefinition `json:"tools,omitempty"`
	ResponseFormat *responseFormat  `json:"response_format,omitempty"`
	Temperature    *float64         `json:"temperature,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message      chatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Client) chatCompletions(ctx context.Context, request chatRequest) (chatResponse, error) {
	if request.Model == "" {
		request.Model = c.model
	}
	body, err := json.Marshal(request)
	if err != nil {
		return chatResponse{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return chatResponse{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		httpRequest.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		if ctx.Err() != nil {
			return chatResponse{}, domain.ErrIntegrationUnavailable.Wrap(fmt.Errorf("%w: %v", ErrUnavailable, ctx.Err()))
		}
		return chatResponse{}, domain.ErrIntegrationUnavailable.Wrap(fmt.Errorf("%w: %v", ErrUnavailable, err))
	}
	defer response.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(response.Body, maximumBodyBytes+1))
	if err != nil {
		return chatResponse{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	if len(payload) > maximumBodyBytes {
		return chatResponse{}, domain.ErrIntegrationUnavailable.Wrap(ErrUnavailable)
	}
	var decoded chatResponse
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return chatResponse{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	if response.StatusCode == http.StatusNotFound {
		return chatResponse{}, domain.ErrIntegrationUnavailable.Wrap(ErrModelUnavailable)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		msg := strings.TrimSpace(string(payload))
		if decoded.Error != nil && decoded.Error.Message != "" {
			msg = decoded.Error.Message
		}
		return chatResponse{}, domain.ErrIntegrationUnavailable.Wrap(fmt.Errorf("%w: status %d: %s", ErrUnavailable, response.StatusCode, msg))
	}
	if len(decoded.Choices) == 0 {
		return chatResponse{}, domain.ErrIntegrationUnavailable.Wrap(ErrInvalidModelOutput)
	}
	return decoded, nil
}
