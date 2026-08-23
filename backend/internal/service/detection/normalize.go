package detection

import (
	"regexp"
	"strings"
)

var (
	rfc3339Pattern    = regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})\b`)
	commonTimePattern = regexp.MustCompile(`\b\d{4}[-/]\d{2}[-/]\d{2}[ T]\d{2}:\d{2}:\d{2}(?:\.\d+)?\b`)
	uuidPattern       = regexp.MustCompile(`(?i)\b[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\b`)
	addressPattern    = regexp.MustCompile(`(?i)\b0x[0-9a-f]{6,}\b`)
	labelledIDPattern = regexp.MustCompile(`(?i)\b(request(?:[_-]?id)?|correlation(?:[_-]?id)?|trace(?:[_-]?id)?|session(?:[_-]?id)?|user(?:[_-]?id)?|order(?:[_-]?id)?)\s*([:=]\s*|\s+)[A-Za-z0-9._:-]+`)
	spacePattern      = regexp.MustCompile(`[ \t]+`)
)

func Normalize(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = rfc3339Pattern.ReplaceAllString(value, "<timestamp>")
	value = commonTimePattern.ReplaceAllString(value, "<timestamp>")
	value = uuidPattern.ReplaceAllString(value, "<uuid>")
	value = addressPattern.ReplaceAllString(value, "<address>")
	value = labelledIDPattern.ReplaceAllString(value, `${1}${2}<id>`)
	lines := strings.Split(value, "\n")
	for index := range lines {
		lines[index] = strings.TrimSpace(spacePattern.ReplaceAllString(lines[index], " "))
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}
