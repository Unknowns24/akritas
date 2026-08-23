package detection

import "regexp"

var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\b(?:authorization\s*:\s*)?(?:bearer|basic)\s+[A-Za-z0-9._~+/=-]+`),
	regexp.MustCompile(`(?i)\b(authorization|token|api[_-]?key|password|passwd|secret|credential|cookie)\s*[:=]\s*[^\s,;]+`),
	regexp.MustCompile(`\b(?:gh[pousr]_[A-Za-z0-9]{20,}|github_pat_[A-Za-z0-9_]{20,})\b`),
	regexp.MustCompile(`(?s)-----BEGIN [^-]*PRIVATE KEY-----.*?-----END [^-]*PRIVATE KEY-----`),
}

func Sanitize(value string) (string, bool) {
	redacted := false
	for index, pattern := range secretPatterns {
		replacement := "[REDACTED]"
		if index == 1 {
			replacement = `${1}=[REDACTED]`
		}
		next := pattern.ReplaceAllString(value, replacement)
		redacted = redacted || next != value
		value = next
	}
	if len(value) > maximumLogicalEventBytes {
		value = value[:maximumLogicalEventBytes]
	}
	return value, redacted
}
