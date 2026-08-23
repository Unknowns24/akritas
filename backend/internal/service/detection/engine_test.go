package detection

import (
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func TestReconstructCombinesCommonStackTraceLines(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	records := []domain.SanitizedLogRecord{
		record(t, now, "ERROR request failed"),
		record(t, now.Add(time.Millisecond), "java.lang.IllegalStateException: boom"),
		record(t, now.Add(2*time.Millisecond), "\tat com.example.App.run(App.java:42)"),
		record(t, now.Add(3*time.Millisecond), "Caused by: java.io.IOException: closed"),
		record(t, now.Add(time.Second), "request complete"),
	}
	events := Reconstruct(records)
	if len(events) != 2 {
		t.Fatalf("events = %d, want 2", len(events))
	}
	if got := events[0].Message(); !strings.Contains(got, "App.java:42") || !strings.Contains(got, "Caused by") {
		t.Fatalf("first event was not reconstructed: %q", got)
	}
}

func TestIgnoredPatternWinsOverPositiveDetection(t *testing.T) {
	engine, err := NewEngine(domain.MonitoringConfiguration{
		Enabled: true, ErrorPatterns: []string{"payment"}, IgnoredPatterns: []string{"healthcheck"},
		GroupingWindow: time.Minute, ContextBefore: 1, ContextAfter: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	event := LogicalEvent{Records: []domain.SanitizedLogRecord{record(t, time.Now().UTC(), "ERROR payment healthcheck failed")}}
	if result := engine.Detect(uuid.New(), event); result != nil {
		t.Fatalf("ignored event detected: %+v", result)
	}
}

func TestBuiltInsAndCustomPatternsAreDeterministic(t *testing.T) {
	engine, err := NewEngine(domain.MonitoringConfiguration{
		Enabled: true, ErrorPatterns: []string{"checkout failed"}, IgnoredPatterns: []string{},
		GroupingWindow: time.Minute, ContextBefore: 0, ContextAfter: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	event := LogicalEvent{Records: []domain.SanitizedLogRecord{record(t, time.Now().UTC(), "FATAL checkout failed: process exited with status 2")}}
	result := engine.Detect(uuid.MustParse("018f4f76-76e3-7f1e-9f4f-6516e81d48b3"), event)
	if result == nil || result.Severity != domain.SeverityCritical {
		t.Fatalf("result = %+v", result)
	}
	want := []string{"fatal_level", "process_crash", "custom:0"}
	if strings.Join(result.Rules, ",") != strings.Join(want, ",") {
		t.Fatalf("rules = %v, want %v", result.Rules, want)
	}
	if !strings.HasPrefix(result.Fingerprint, "sha256:") || len(result.Fingerprint) != 71 {
		t.Fatalf("fingerprint = %q", result.Fingerprint)
	}
}

func TestNoPositiveRuleProducesNoDetection(t *testing.T) {
	engine, err := NewEngine(domain.MonitoringConfiguration{
		Enabled: true, ErrorPatterns: []string{}, IgnoredPatterns: []string{},
		GroupingWindow: time.Minute, ContextBefore: 0, ContextAfter: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := engine.Detect(uuid.New(), LogicalEvent{Records: []domain.SanitizedLogRecord{record(t, time.Now().UTC(), "GET /health 200")}}); got != nil {
		t.Fatalf("unexpected detection: %+v", got)
	}
}

func TestEveryBuiltInRuleHasAnExplicitPositiveExample(t *testing.T) {
	cases := map[string]string{
		"error_level":       "\x1b[31m[error]\x1b[0m failed to initialize database",
		"fatal_level":       "[fatal] startup failed",
		"panic":             "panic: nil pointer",
		"stack_trace":       "Traceback (most recent call last):\nFile \"app.py\", line 42",
		"http_5xx":          "GET /checkout 503",
		"process_crash":     "process exited with status 2",
		"container_restart": "container restart-loop detected",
	}
	engine, err := NewEngine(domain.MonitoringConfiguration{Enabled: true, ErrorPatterns: []string{}, IgnoredPatterns: []string{}, GroupingWindow: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	for rule, message := range cases {
		t.Run(rule, func(t *testing.T) {
			result := engine.Detect(uuid.New(), LogicalEvent{Records: []domain.SanitizedLogRecord{record(t, time.Now().UTC(), message)}})
			if result == nil || !contains(result.Rules, rule) {
				t.Fatalf("rules = %+v", result)
			}
		})
	}
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func TestNormalizationAndFingerprintIgnoreOnlyVolatileValues(t *testing.T) {
	projectID := uuid.MustParse("018f4f76-76e3-7f1e-9f4f-6516e81d48b3")
	a := "2026-08-22T12:00:01.123Z ERROR request_id=018f4f76-76e3-7f1e-9f4f-6516e81d48a2 at App.go:42 status=503 port=8080"
	b := "2026-08-22T12:01:07.456Z ERROR request_id=118f4f76-76e3-7f1e-9f4f-6516e81d48a9 at App.go:42 status=503 port=8080"
	if Fingerprint(projectID, "error_level", Normalize(a)) != Fingerprint(projectID, "error_level", Normalize(b)) {
		t.Fatal("volatile values changed fingerprint")
	}
	c := strings.ReplaceAll(b, "App.go:42", "Other.go:99")
	if Fingerprint(projectID, "error_level", Normalize(a)) == Fingerprint(projectID, "error_level", Normalize(c)) {
		t.Fatal("materially different stack location collapsed")
	}
}

func TestSanitizeRedactsCredentialsAndPEM(t *testing.T) {
	input := "Authorization: Bearer abc.def token=secret password=hunter2\n-----BEGIN PRIVATE KEY-----\nabc\n-----END PRIVATE KEY-----"
	got, redacted := Sanitize(input)
	if !redacted || strings.Contains(got, "abc.def") || strings.Contains(got, "hunter2") || strings.Contains(got, "PRIVATE KEY-----\nabc") {
		t.Fatalf("secret remained in %q", got)
	}
}

func TestSanitizeStripsANSIControlSequences(t *testing.T) {
	got, _ := Sanitize("\x1b[0m\x1b[31m[error]\x1b[0m failed")
	if got != "[error] failed" {
		t.Fatalf("sanitized = %q", got)
	}
}

func record(t *testing.T, timestamp time.Time, message string) domain.SanitizedLogRecord {
	t.Helper()
	value, err := domain.NewSanitizedLogRecord(timestamp, domain.LogStreamUnknown, message)
	if err != nil {
		t.Fatal(err)
	}
	return value
}
