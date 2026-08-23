package qvac

import "fmt"

// NewDefaultRunner builds a production InvestigationRunner against the local
// QVAC OpenAI-compatible server using default loopback settings and no tools.
// Tools are attached by higher-level wiring once repository tools are available.
func NewDefaultRunner(tools *ToolRegistry) (*Runner, error) {
	client, err := NewClient(ClientConfig{})
	if err != nil {
		return nil, fmt.Errorf("qvac client: %w", err)
	}
	return NewRunner(client, tools, RunnerConfig{})
}
