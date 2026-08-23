package qvac

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Tool is a read-only capability the local model may invoke during an investigation.
type Tool interface {
	Name() string
	Description() string
	Parameters() json.RawMessage
	Execute(ctx context.Context, arguments json.RawMessage) (string, error)
}

// ToolRegistry is an explicit allowlist of tools. Unknown names never execute.
type ToolRegistry struct {
	tools map[string]Tool
}

func NewToolRegistry(tools ...Tool) *ToolRegistry {
	registry := &ToolRegistry{tools: make(map[string]Tool, len(tools))}
	for _, tool := range tools {
		if tool == nil || strings.TrimSpace(tool.Name()) == "" {
			continue
		}
		registry.tools[tool.Name()] = tool
	}
	return registry
}

func (r *ToolRegistry) definitions() []toolDefinition {
	if r == nil || len(r.tools) == 0 {
		return nil
	}
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	sort.Strings(names)
	defs := make([]toolDefinition, 0, len(names))
	for _, name := range names {
		tool := r.tools[name]
		params := tool.Parameters()
		if len(params) == 0 {
			params = json.RawMessage(`{"type":"object","properties":{}}`)
		}
		defs = append(defs, toolDefinition{
			Type: "function",
			Function: toolFunctionSchema{
				Name: tool.Name(), Description: tool.Description(), Parameters: params,
			},
		})
	}
	return defs
}

func (r *ToolRegistry) execute(ctx context.Context, name string, arguments string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("%w: %s", ErrUnknownTool, name)
	}
	tool, ok := r.tools[name]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnknownTool, name)
	}
	raw := json.RawMessage(strings.TrimSpace(arguments))
	if len(raw) == 0 {
		raw = json.RawMessage(`{}`)
	}
	if !json.Valid(raw) {
		return "", fmt.Errorf("invalid arguments for tool %s", name)
	}
	return tool.Execute(ctx, raw)
}

// DefinitionsForTest exposes allowlisted tool names for tests.
func (r *ToolRegistry) DefinitionsForTest() []string {
	defs := r.definitions()
	names := make([]string, 0, len(defs))
	for _, def := range defs {
		names = append(names, def.Function.Name)
	}
	return names
}

// ExecuteForTest executes an allowlisted tool for tests.
func (r *ToolRegistry) ExecuteForTest(ctx context.Context, name, arguments string) (string, error) {
	return r.execute(ctx, name, arguments)
}
