// Package issuecontent builds deterministic, bounded and sanitized GitHub
// Issue content. It performs no persistence, credential access or I/O.
package issuecontent

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/service/evidencesafety"
	"github.com/google/uuid"
)

const (
	maximumTitleBytes        = 256
	maximumBodyBytes         = 60 << 10
	maximumEvidenceItemBytes = 8 << 10
)

var ErrInvalidInput = errors.New("invalid Issue content input")

type Input = portsout.IssueContentInput

type Builder struct{}

func New() *Builder { return &Builder{} }

func (b *Builder) BuildIssueContent(input portsout.IssueContentInput) (portsout.IssueContent, error) {
	return b.Build(Input(input))
}

func (b *Builder) Build(input Input) (portsout.IssueContent, error) {
	if input.Project.Validate() != nil || input.Incident.Validate() != nil || input.Investigation.Validate() != nil ||
		input.Investigation.Status != domain.InvestigationStatusCompleted || input.Investigation.IncidentID != input.Incident.ID ||
		input.Incident.ProjectID != input.Project.ID {
		return portsout.IssueContent{}, ErrInvalidInput
	}

	title := safeBounded(fmt.Sprintf("[%s] %s", input.Incident.Key, input.Incident.Title), maximumTitleBytes)
	var body strings.Builder
	fmt.Fprintf(&body, "<!-- akritas:investigation_id=%s -->\n", input.Investigation.ID)
	body.WriteString("# Incident " + scalar(input.Incident.Key) + "\n\n")
	body.WriteString("> Automatically documented by Akritas from a completed local QVAC investigation.\n\n")
	body.WriteString("## Project Context\n\n")
	bullet(&body, "Project", input.Project.Name)
	bullet(&body, "Application", input.Project.DokploySource.DisplayName)
	bullet(&body, "Environment", valueOr(input.Project.DokploySource.Environment, "not specified"))
	bullet(&body, "Repository", input.Project.GitHubRepository.FullName)
	bullet(&body, "Default branch", input.Project.GitHubRepository.DefaultBranch)

	body.WriteString("\n## Incident — Observed\n\n")
	bullet(&body, "Incident ID", input.Incident.ID.String())
	bullet(&body, "Fingerprint", input.Incident.Fingerprint)
	bullet(&body, "Severity", string(input.Incident.Severity))
	bullet(&body, "Occurrences", fmt.Sprintf("%d", input.Incident.OccurrenceCount))
	bullet(&body, "First seen", input.Incident.FirstSeenAt.UTC().Format(time.RFC3339))
	bullet(&body, "Last seen", input.Incident.LastSeenAt.UTC().Format(time.RFC3339))
	bullet(&body, "Observed title", input.Incident.Title)

	body.WriteString("\n## Observed Evidence\n\n")
	evidence := orderedEvidence(input.Evidence, input.Investigation.EvidenceIDs)
	omitted := 0
	for index := range evidence {
		section := evidenceSection(evidence[index])
		if body.Len()+len(section)+2048 > maximumBodyBytes {
			omitted = len(evidence) - index
			break
		}
		body.WriteString(section)
	}
	if len(evidence) == 0 {
		body.WriteString("No persisted Evidence was available.\n")
	}
	if omitted > 0 {
		fmt.Fprintf(&body, "\n_%d Evidence items were omitted to keep this Issue bounded._\n", omitted)
	}

	body.WriteString("\n## Investigation — QVAC Analysis\n\n")
	body.WriteString("> The following fields are model-generated conclusions, not independently verified facts.\n\n")
	field(&body, "Root Cause Status", string(*input.Investigation.RootCauseStatus))
	field(&body, "Root Cause / Hypothesis", valueOr(input.Investigation.RootCause, "No root cause was identified."))
	field(&body, "Confidence", fmt.Sprintf("%.4f", *input.Investigation.Confidence))
	field(&body, "Resolution Status", string(*input.Investigation.ResolutionStatus))
	field(&body, "Investigation Summary", input.Investigation.Summary)
	list(&body, "Hypotheses", input.Investigation.Hypotheses)
	list(&body, "Relevant Files", input.Investigation.RelevantFiles)
	list(&body, "Relevant Commits", input.Investigation.RelevantCommits)
	list(&body, "Recommended Actions", input.Investigation.RecommendedActions)

	return portsout.IssueContent{Title: title, Body: safeBounded(body.String(), maximumBodyBytes)}, nil
}

func orderedEvidence(items []domain.Evidence, cited []uuid.UUID) []domain.Evidence {
	citations := make(map[uuid.UUID]struct{}, len(cited))
	for _, id := range cited {
		citations[id] = struct{}{}
	}
	seen := make(map[uuid.UUID]struct{}, len(items))
	result := make([]domain.Evidence, 0, len(items))
	for _, item := range items {
		if item.Validate() != nil {
			continue
		}
		if _, duplicate := seen[item.ID]; duplicate {
			continue
		}
		seen[item.ID] = struct{}{}
		result = append(result, item)
	}
	sort.Slice(result, func(left, right int) bool {
		_, leftCited := citations[result[left].ID]
		_, rightCited := citations[result[right].ID]
		if leftCited != rightCited {
			return leftCited
		}
		if !result[left].CreatedAt.Equal(result[right].CreatedAt) {
			return result[left].CreatedAt.Before(result[right].CreatedAt)
		}
		return result[left].ID.String() < result[right].ID.String()
	})
	return result
}

func evidenceSection(item domain.Evidence) string {
	var section strings.Builder
	fmt.Fprintf(&section, "### Evidence %s — %s\n\n", item.ID, scalar(string(item.Type)))
	bullet(&section, "Summary", item.Summary)
	if item.OccurredAt != nil {
		bullet(&section, "Occurred at", item.OccurredAt.UTC().Format(time.RFC3339))
	}
	if item.FilePath != "" {
		bullet(&section, "File", item.FilePath)
	}
	if item.CommitSHA != "" {
		bullet(&section, "Commit", item.CommitSHA)
	}
	excerpt := strings.TrimSpace(strings.Join([]string{item.Content, item.Patch}, "\n"))
	if excerpt != "" {
		section.WriteString("\nObserved excerpt:\n\n")
		for _, line := range strings.Split(bounded(evidencesafety.Redact(excerpt), maximumEvidenceItemBytes), "\n") {
			section.WriteString("    " + line + "\n")
		}
	}
	section.WriteByte('\n')
	return section.String()
}

func field(builder *strings.Builder, heading, value string) {
	builder.WriteString("### " + heading + "\n\n" + scalar(valueOr(value, "not available")) + "\n\n")
}

func list(builder *strings.Builder, heading string, values []string) {
	builder.WriteString("### " + heading + "\n\n")
	values = uniqueSorted(values)
	if len(values) == 0 {
		builder.WriteString("- None\n\n")
		return
	}
	for _, value := range values {
		builder.WriteString("- " + scalar(value) + "\n")
	}
	builder.WriteByte('\n')
}

func uniqueSorted(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(evidencesafety.Redact(value))
		if value == "" {
			continue
		}
		if _, duplicate := seen[value]; duplicate {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func bullet(builder *strings.Builder, label, value string) {
	fmt.Fprintf(builder, "- **%s:** %s\n", label, scalar(valueOr(value, "not available")))
}

func scalar(value string) string {
	value = strings.TrimSpace(evidencesafety.Redact(value))
	value = strings.NewReplacer("\\", "\\\\", "`", "\\`", "\r", " ", "\n", " ").Replace(value)
	return value
}

func valueOr(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func bounded(value string, maximum int) string {
	if maximum <= 0 || len(value) <= maximum {
		return value
	}
	const suffix = "\n[TRUNCATED]"
	limit := maximum - len(suffix)
	if limit < 0 {
		return ""
	}
	value = value[:limit]
	for !utf8.ValidString(value) {
		value = value[:len(value)-1]
	}
	return value + suffix
}

func safeBounded(value string, maximum int) string {
	return bounded(evidencesafety.Redact(bounded(evidencesafety.Redact(value), maximum)), maximum)
}
