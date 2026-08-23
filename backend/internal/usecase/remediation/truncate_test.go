package remediation

import "testing"

func TestTruncateWithMarkerNoOpBelowLimit(t *testing.T) {
	got := truncateWithMarker("short", 100)
	if got != "short" {
		t.Fatalf("expected untouched string, got %q", got)
	}
}

func TestTruncateWithMarkerCapsAtLimit(t *testing.T) {
	input := make([]byte, 1000)
	for i := range input {
		input[i] = 'x'
	}
	got := truncateWithMarker(string(input), 100)
	if len(got) > 100 {
		t.Fatalf("expected result capped at 100 bytes, got %d", len(got))
	}
	if got == string(input)[:100] {
		t.Fatal("expected a visible truncation marker, not a bare cut")
	}
}
