package evidencesafety

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var secretPatterns = []struct {
	pattern *regexp.Regexp
	replace string
}{
	{regexp.MustCompile(`(?i)(authorization\s*:\s*bearer\s+)[^\s,;]+`), `${1}[REDACTED]`},
	{regexp.MustCompile(`(?i)\b(gh[pousr]_[A-Za-z0-9_]{20,}|github_pat_[A-Za-z0-9_]{20,})\b`), `[REDACTED_GITHUB_TOKEN]`},
	{regexp.MustCompile(`(?i)\b(AKIA|ASIA)[A-Z0-9]{16}\b`), `[REDACTED_ACCESS_KEY]`},
	{regexp.MustCompile(`(?i)\beyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\b`), `[REDACTED_SESSION_TOKEN]`},
	{regexp.MustCompile(`(?is)-----BEGIN [^-\n]*PRIVATE KEY-----.*?-----END [^-\n]*PRIVATE KEY-----`), `[REDACTED_PRIVATE_KEY]`},
	{regexp.MustCompile(`(?i)\b([A-Z][A-Z0-9_]*(?:TOKEN|SECRET|PASSWORD|PASSWD|API_KEY|PRIVATE_KEY|COOKIE|SESSION)[A-Z0-9_]*)\b\s*=\s*[^\s,;]+`), `${1}=[REDACTED]`},
	{regexp.MustCompile(`(?i)\b(https?|postgres(?:ql)?|mysql|redis)://[^/\s:@]+:[^@\s/]+@`), `${1}://[REDACTED]@`},
	{regexp.MustCompile(`(?i)\b(api[_-]?key|token|password|passwd|secret|cookie|session)\b\s*[:=]\s*[^\s,;]+`), `${1}=[REDACTED]`},
}

func Redact(value string) string {
	value = strings.ToValidUTF8(value, "�")
	for _, candidate := range secretPatterns {
		value = candidate.pattern.ReplaceAllString(value, candidate.replace)
	}
	return value
}

func RedactAndLimit(value string, maximumBytes int) string {
	value = Redact(value)
	if maximumBytes <= 0 || len(value) <= maximumBytes {
		return value
	}
	value = value[:maximumBytes]
	for !utf8.ValidString(value) {
		value = value[:len(value)-1]
	}
	return value + "\n[TRUNCATED]"
}
