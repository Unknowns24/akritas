package investigationtools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Unknowns24/akritas/backend/internal/adapter/external/qvac"
)

// Registry builds the allowlisted QVAC tools bound to a fixed repository scope.
func Registry(api GitHubAPI, scope Scope) *qvac.ToolRegistry {
	return qvac.NewToolRegistry(
		&searchCodeTool{api: api, scope: scope},
		&readFileTool{api: api, scope: scope},
		&listRecentCommitsTool{api: api, scope: scope},
		&readCommitTool{api: api, scope: scope},
		&readDiffTool{api: api, scope: scope},
	)
}

type searchCodeTool struct {
	api   GitHubAPI
	scope Scope
}

func (t *searchCodeTool) Name() string { return "search_code" }
func (t *searchCodeTool) Description() string {
	return "Search code in the incident's linked GitHub repository (read-only)."
}
func (t *searchCodeTool) Parameters() json.RawMessage {
	return json.RawMessage(`{"type":"object","additionalProperties":false,"required":["query"],"properties":{"query":{"type":"string"}}}`)
}
func (t *searchCodeTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	var args struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil || args.Query == "" {
		return "", fmt.Errorf("query is required")
	}
	matches, err := t.api.SearchCode(ctx, t.scope.Account, t.scope.Owner, t.scope.Name, args.Query)
	if err != nil {
		return "", err
	}
	return marshalToolResult(matches)
}

type readFileTool struct {
	api   GitHubAPI
	scope Scope
}

func (t *readFileTool) Name() string { return "read_file" }
func (t *readFileTool) Description() string {
	return "Read a file from the incident's linked GitHub repository (read-only)."
}
func (t *readFileTool) Parameters() json.RawMessage {
	return json.RawMessage(`{"type":"object","additionalProperties":false,"required":["path"],"properties":{"path":{"type":"string"},"ref":{"type":"string"}}}`)
}
func (t *readFileTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	var args struct {
		Path string `json:"path"`
		Ref  string `json:"ref"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil || args.Path == "" {
		return "", fmt.Errorf("path is required")
	}
	ref := args.Ref
	if ref == "" {
		ref = t.scope.Branch
	}
	file, err := t.api.ReadFile(ctx, t.scope.Account, t.scope.Owner, t.scope.Name, args.Path, ref)
	if err != nil {
		return "", err
	}
	return marshalToolResult(file)
}

type listRecentCommitsTool struct {
	api   GitHubAPI
	scope Scope
}

func (t *listRecentCommitsTool) Name() string { return "list_recent_commits" }
func (t *listRecentCommitsTool) Description() string {
	return "List recent commits on the incident repository default branch (read-only)."
}
func (t *listRecentCommitsTool) Parameters() json.RawMessage {
	return json.RawMessage(`{"type":"object","additionalProperties":false,"properties":{"limit":{"type":"integer","minimum":1,"maximum":30}}}`)
}
func (t *listRecentCommitsTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	var args struct {
		Limit int `json:"limit"`
	}
	_ = json.Unmarshal(arguments, &args)
	commits, err := t.api.ListRecentCommits(ctx, t.scope.Account, t.scope.Owner, t.scope.Name, t.scope.Branch, args.Limit)
	if err != nil {
		return "", err
	}
	return marshalToolResult(commits)
}

type readCommitTool struct {
	api   GitHubAPI
	scope Scope
}

func (t *readCommitTool) Name() string { return "read_commit" }
func (t *readCommitTool) Description() string {
	return "Read commit metadata for a SHA in the incident repository (read-only)."
}
func (t *readCommitTool) Parameters() json.RawMessage {
	return json.RawMessage(`{"type":"object","additionalProperties":false,"required":["sha"],"properties":{"sha":{"type":"string"}}}`)
}
func (t *readCommitTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	var args struct {
		SHA string `json:"sha"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil || args.SHA == "" {
		return "", fmt.Errorf("sha is required")
	}
	detail, err := t.api.ReadCommit(ctx, t.scope.Account, t.scope.Owner, t.scope.Name, args.SHA)
	if err != nil {
		return "", err
	}
	return marshalToolResult(detail)
}

type readDiffTool struct {
	api   GitHubAPI
	scope Scope
}

func (t *readDiffTool) Name() string { return "read_diff" }
func (t *readDiffTool) Description() string {
	return "Read the patch/diff for a commit SHA in the incident repository (read-only)."
}
func (t *readDiffTool) Parameters() json.RawMessage {
	return json.RawMessage(`{"type":"object","additionalProperties":false,"required":["sha"],"properties":{"sha":{"type":"string"}}}`)
}
func (t *readDiffTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	var args struct {
		SHA string `json:"sha"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil || args.SHA == "" {
		return "", fmt.Errorf("sha is required")
	}
	diff, err := t.api.ReadDiff(ctx, t.scope.Account, t.scope.Owner, t.scope.Name, args.SHA)
	if err != nil {
		return "", err
	}
	return marshalToolResult(map[string]string{"sha": args.SHA, "diff": diff})
}

func marshalToolResult(value any) (string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
