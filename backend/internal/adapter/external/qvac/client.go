package qvac

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	EndpointURL   string
	APIKey        string
	BasicUsername string
	BasicPassword string
	HTTPClient    *http.Client
	Model         string
	Timeout       time.Duration
	ContextSize   int
}

type Client struct {
	base          string
	apiKey        string
	basicUsername string
	basicPassword string
	httpClient    *http.Client
	model         string
	contextSize   int
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
		timeout := config.Timeout
		if timeout <= 0 {
			timeout = 2 * time.Minute
		}
		httpClient = &http.Client{Timeout: timeout}
	}
	model := strings.TrimSpace(config.Model)
	if model == "" {
		model = defaultModel
	}
	contextSize := config.ContextSize
	if contextSize <= 0 {
		contextSize = domain.DefaultQvacContextSize
	}
	return &Client{
		base:          strings.TrimRight(base.String(), "/"),
		apiKey:        config.APIKey,
		basicUsername: strings.TrimSpace(config.BasicUsername),
		basicPassword: config.BasicPassword,
		httpClient:    httpClient,
		model:         model,
		contextSize:   contextSize,
	}, nil
}

func (c *Client) Model() string    { return c.model }
func (c *Client) Endpoint() string { return c.base }
func (c *Client) ContextSize() int { return c.contextSize }

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

type chatOptions struct {
	NumCtx int `json:"num_ctx,omitempty"`
}

type chatRequest struct {
	Model           string           `json:"model"`
	Messages        []chatMessage    `json:"messages"`
	Tools           []toolDefinition `json:"tools,omitempty"`
	ToolChoice      string           `json:"tool_choice,omitempty"`
	Temperature     *float64         `json:"temperature,omitempty"`
	MaxTokens       *int             `json:"max_tokens,omitempty"`
	ReasoningBudget *bool            `json:"reasoning_budget,omitempty"`
	Options         *chatOptions     `json:"options,omitempty"`
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
	if request.Options == nil && c.contextSize > 0 {
		request.Options = &chatOptions{NumCtx: c.contextSize}
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
	} else if c.basicUsername != "" {
		httpRequest.SetBasicAuth(c.basicUsername, c.basicPassword)
	}
	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		if ctx.Err() != nil {
			log.Printf("qvac: chat completion failed endpoint=%s model=%s context_size=%d error=%v context_error=%v", c.base, request.Model, c.contextSize, err, ctx.Err())
			return chatResponse{}, domain.ErrQvacUnavailable.Wrap(fmt.Errorf("%w: %v", ErrUnavailable, ctx.Err()))
		}
		log.Printf("qvac: chat completion failed endpoint=%s model=%s context_size=%d error=%v", c.base, request.Model, c.contextSize, err)
		return chatResponse{}, domain.ErrQvacUnavailable.Wrap(fmt.Errorf("%w: %v", ErrUnavailable, err))
	}
	defer response.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(response.Body, maximumBodyBytes+1))
	if err != nil {
		return chatResponse{}, domain.ErrQvacUnavailable.Wrap(err)
	}
	if len(payload) > maximumBodyBytes {
		return chatResponse{}, domain.ErrQvacUnavailable.Wrap(ErrUnavailable)
	}
	var decoded chatResponse
	if response.StatusCode == http.StatusNotFound {
		return chatResponse{}, domain.ErrQvacUnavailable.Wrap(ErrModelUnavailable)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		msg := strings.TrimSpace(string(payload))
		if err := json.Unmarshal(payload, &decoded); err == nil && decoded.Error != nil && decoded.Error.Message != "" {
			msg = decoded.Error.Message
		}
		if isContextOverflowResponse(msg) {
			log.Printf("qvac: chat completion context overflow endpoint=%s model=%s context_size=%d status=%d", c.base, request.Model, c.contextSize, response.StatusCode)
			return chatResponse{}, domain.ErrQvacContextOverflow.Wrap(fmt.Errorf("%w: status %d", ErrContextOverflow, response.StatusCode))
		}
		log.Printf("qvac: chat completion unavailable endpoint=%s model=%s context_size=%d status=%d response=%q", c.base, request.Model, c.contextSize, response.StatusCode, msg)
		return chatResponse{}, domain.ErrQvacUnavailable.Wrap(fmt.Errorf("%w: status %d: %s", ErrUnavailable, response.StatusCode, msg))
	}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return chatResponse{}, domain.ErrQvacUnavailable.Wrap(err)
	}
	if len(decoded.Choices) == 0 {
		return chatResponse{}, domain.ErrQvacUnavailable.Wrap(ErrInvalidModelOutput)
	}
	return decoded, nil
}

func isContextOverflowResponse(message string) bool {
	normalized := strings.ToLower(message)
	return strings.Contains(normalized, "context_overflow") ||
		strings.Contains(normalized, "context overflow") ||
		strings.Contains(normalized, "prompt exceeds") ||
		strings.Contains(normalized, "context window")
}

func (c *Client) Ping(ctx context.Context) (time.Duration, error) {
	started := time.Now()
	_, err := c.chatCompletions(ctx, chatRequest{
		Messages:        []chatMessage{{Role: "user", Content: "Return OK."}},
		MaxTokens:       intPtr(8),
		ReasoningBudget: boolPtr(false),
	})
	return time.Since(started), err
}
