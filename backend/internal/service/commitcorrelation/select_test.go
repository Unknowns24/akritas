package commitcorrelation

import (
	"strings"
	"testing"
	"time"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func TestSelectIsDeterministicBoundedDeduplicatedAndSanitized(t *testing.T) {
	t.Parallel()

	firstSeen := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	commits := []portsout.RepositoryCommitSummary{
		{SHA: "bbb", Date: firstSeen.Add(-2 * time.Hour).Format(time.RFC3339), Author: "dev", Message: "TOKEN=secret-value fix users", URL: "https://example.test/bbb"},
		{SHA: "aaa", Date: firstSeen.Add(-2 * time.Hour).Format(time.RFC3339), Author: "dev", Message: "safe", URL: "https://example.test/aaa"},
		{SHA: "bbb", Date: firstSeen.Add(-time.Hour).Format(time.RFC3339), Author: "dupe", Message: "duplicate should be ignored"},
		{SHA: "old", Date: firstSeen.Add(-72 * time.Hour).Format(time.RFC3339), Author: "dev", Message: "too old"},
	}

	got := Select(IncidentWindow{FirstSeenAt: firstSeen, LastSeenAt: firstSeen.Add(time.Minute)}, commits, 2)
	if len(got) != 2 {
		t.Fatalf("len=%d got=%+v", len(got), got)
	}
	if got[0].SHA != "aaa" || got[1].SHA != "bbb" {
		t.Fatalf("unexpected deterministic order: %+v", got)
	}
	for _, commit := range got {
		if strings.Contains(commit.Message, "secret-value") {
			t.Fatalf("commit message leaked secret: %+v", commit)
		}
		if !strings.Contains(commit.Reason, "potencialmente relacionado") {
			t.Fatalf("commit reason must avoid causal language: %+v", commit)
		}
	}
}
