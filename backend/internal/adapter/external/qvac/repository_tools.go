package qvac

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func newRepositoryToolRegistry(inspector portsout.RepositoryInspector, scope portsout.RepositoryScope) *ToolRegistry {
	if inspector == nil {
		return NewToolRegistry()
	}
	return NewToolRegistry(
		&searchCodeTool{inspector: inspector, scope: scope},
		&readFileTool{inspector: inspector, scope: scope},
		&listRecentCommitsTool{inspector: inspector, scope: scope},
		&readCommitTool{inspector: inspector, scope: scope},
		&readDiffTool{inspector: inspector, scope: scope},
	)
}

type searchCodeTool struct {
	inspector portsout.RepositoryInspector
	scope     portsout.RepositoryScope
}

func (t *searchCodeTool) Name() string { return "search_code" }
func (t *searchCodeTool) Description() string {
	return "Search code in the Incident Project's configured repository (read-only)."
}
func (t *searchCodeTool) Parameters() json.RawMessage {
	return json.RawMessage(`{"type":"object","additionalProperties":false,"required":["query"],"properties":{"query":{"type":"string","minLength":1,"maxLength":500}}}`)
}
func (t *searchCodeTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	var args struct {
		Query string `json:"query"`
	}
	if err := decodeToolArguments(arguments, &args); err != nil || strings.TrimSpace(args.Query) == "" || len(args.Query) > 500 {
		return "", fmt.Errorf("query is invalid")
	}
	value, err := t.inspector.SearchCode(ctx, t.scope.Account, t.scope.Owner, t.scope.Name, args.Query)
	return marshalToolResult(value, err)
}

type readFileTool struct {
	inspector portsout.RepositoryInspector
	scope     portsout.RepositoryScope
}

func (t *readFileTool) Name() string { return "read_file" }
func (t *readFileTool) Description() string {
	return "Read a file from the configured repository default branch (read-only)."
}
func (t *readFileTool) Parameters() json.RawMessage {
	return json.RawMessage(`{"type":"object","additionalProperties":false,"required":["path"],"properties":{"path":{"type":"string","minLength":1,"maxLength":4096}}}`)
}
func (t *readFileTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	var args struct {
		Path string `json:"path"`
	}
	if err := decodeToolArguments(arguments, &args); err != nil || strings.TrimSpace(args.Path) == "" || len(args.Path) > 4096 {
		return "", fmt.Errorf("path is invalid")
	}
	value, err := t.inspector.ReadFile(ctx, t.scope.Account, t.scope.Owner, t.scope.Name, args.Path, t.scope.Branch)
	return marshalToolResult(value, err)
}

type listRecentCommitsTool struct {
	inspector portsout.RepositoryInspector
	scope     portsout.RepositoryScope
}

func (t *listRecentCommitsTool) Name() string { return "list_recent_commits" }
func (t *listRecentCommitsTool) Description() string {
	return "List recent commits on the configured repository default branch (read-only)."
}
func (t *listRecentCommitsTool) Parameters() json.RawMessage {
	return json.RawMessage(`{"type":"object","additionalProperties":false,"properties":{"limit":{"type":"integer","minimum":1,"maximum":30}}}`)
}
func (t *listRecentCommitsTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	var args struct {
		Limit int `json:"limit"`
	}
	if err := decodeToolArguments(arguments, &args); err != nil {
		return "", fmt.Errorf("limit is invalid")
	}
	if args.Limit == 0 {
		args.Limit = 10
	}
	if args.Limit < 1 || args.Limit > 30 {
		return "", fmt.Errorf("limit is invalid")
	}
	value, err := t.inspector.ListRecentCommits(ctx, t.scope.Account, t.scope.Owner, t.scope.Name, t.scope.Branch, args.Limit)
	return marshalToolResult(value, err)
}

type readCommitTool struct {
	inspector portsout.RepositoryInspector
	scope     portsout.RepositoryScope
}

func (t *readCommitTool) Name() string { return "read_commit" }
func (t *readCommitTool) Description() string {
	return "Read commit metadata in the configured repository (read-only)."
}
func (t *readCommitTool) Parameters() json.RawMessage { return shaParameters() }
func (t *readCommitTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	sha, err := decodeSHA(arguments)
	if err != nil {
		return "", err
	}
	value, err := t.inspector.ReadCommit(ctx, t.scope.Account, t.scope.Owner, t.scope.Name, sha)
	return marshalToolResult(value, err)
}

type readDiffTool struct {
	inspector portsout.RepositoryInspector
	scope     portsout.RepositoryScope
}

func (t *readDiffTool) Name() string { return "read_diff" }
func (t *readDiffTool) Description() string {
	return "Read a commit patch in the configured repository (read-only)."
}
func (t *readDiffTool) Parameters() json.RawMessage { return shaParameters() }
func (t *readDiffTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	sha, err := decodeSHA(arguments)
	if err != nil {
		return "", err
	}
	value, err := t.inspector.ReadDiff(ctx, t.scope.Account, t.scope.Owner, t.scope.Name, sha)
	if err != nil {
		return "", err
	}
	return marshalToolResult(map[string]string{"sha": sha, "diff": value}, nil)
}

func shaParameters() json.RawMessage {
	return json.RawMessage(`{"type":"object","additionalProperties":false,"required":["sha"],"properties":{"sha":{"type":"string","minLength":7,"maxLength":64}}}`)
}

func decodeSHA(arguments json.RawMessage) (string, error) {
	var args struct {
		SHA string `json:"sha"`
	}
	if err := decodeToolArguments(arguments, &args); err != nil {
		return "", fmt.Errorf("sha is invalid")
	}
	args.SHA = strings.TrimSpace(args.SHA)
	if len(args.SHA) < 7 || len(args.SHA) > 64 || strings.ContainsAny(args.SHA, "/\\") || strings.Contains(args.SHA, "..") {
		return "", fmt.Errorf("sha is invalid")
	}
	return args.SHA, nil
}

func decodeToolArguments(arguments json.RawMessage, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(arguments))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("multiple JSON values")
	}
	return nil
}

func marshalToolResult(value any, err error) (string, error) {
	if err != nil {
		return "", err
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
