package evidencesafety

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type RedactionResult struct {
	Value     string
	Redacted  bool
	Truncated bool
}

var secretPatterns = []struct {
	pattern *regexp.Regexp
	replace string
}{
	{regexp.MustCompile(`(?is)-----BEGIN [^-\n]*PRIVATE KEY-----.*?-----END [^-\n]*PRIVATE KEY-----`), `[REDACTED_PRIVATE_KEY]`},
	{regexp.MustCompile(`(?i)(authorization\s*:\s*(?:bearer|basic)\s+)[^\r\n,;]+`), `${1}[REDACTED]`},
	{regexp.MustCompile(`(?i)(set-cookie\s*:\s*)[^\r\n]+`), `${1}[REDACTED]`},
	{regexp.MustCompile(`(?i)(cookie\s*:\s*)[^\r\n]+`), `${1}[REDACTED]`},
	{regexp.MustCompile(`(?i)\b(github_pat_[A-Za-z0-9_]{20,}|gh[pousr]_[A-Za-z0-9_]{20,})\b`), `[REDACTED_GITHUB_TOKEN]`},
	{regexp.MustCompile(`(?i)\b(AKIA|ASIA)[A-Z0-9]{16}\b`), `[REDACTED_ACCESS_KEY]`},
	{regexp.MustCompile(`(?i)\beyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\b`), `[REDACTED_SESSION_TOKEN]`},
	{regexp.MustCompile(`(?i)\b((?:postgres(?:ql)?|mysql|mariadb|redis|mongodb(?:\+srv)?|amqp|rabbitmq|sqlserver)://)(?:[^/@\s:]*:)?[^/@\s]+@`), `${1}[REDACTED]@`},
	{regexp.MustCompile(`(?i)("[A-Za-z0-9_-]*(?:token|password|passwd|secret|api[_-]?key|private[_-]?key|cookie|session)[A-Za-z0-9_-]*"\s*:\s*)"(?:\\.|[^"\\])*"`), `${1}"[REDACTED]"`},
	{regexp.MustCompile(`(?i)\b([A-Za-z0-9_-]*(?:token|password|passwd|secret|api[_-]?key|private[_-]?key|cookie|session)[A-Za-z0-9_-]*)\b\s*[:=]\s*"(?:\\.|[^"\\])*"`), `${1}="[REDACTED]"`},
	{regexp.MustCompile(`(?i)\b([A-Za-z0-9_-]*(?:token|password|passwd|secret|api[_-]?key|private[_-]?key|cookie|session)[A-Za-z0-9_-]*)\b\s*[:=]\s*'[^']*'`), `${1}='[REDACTED]'`},
	{regexp.MustCompile(`(?i)\b([A-Za-z0-9_-]*(?:token|password|passwd|secret|api[_-]?key|private[_-]?key|cookie|session)[A-Za-z0-9_-]*)\b\s*[:=]\s*[^'"\[\s,;]+`), `${1}=[REDACTED]`},
}

func Redact(value string) string {
	return RedactWithReport(value).Value
}

func RedactWithReport(value string) RedactionResult {
	value = strings.ToValidUTF8(value, "\uFFFD")
	redacted := false
	for _, candidate := range secretPatterns {
		next := candidate.pattern.ReplaceAllString(value, candidate.replace)
		if next != value {
			redacted = true
		}
		value = next
	}
	return RedactionResult{Value: value, Redacted: redacted}
}

func RedactAndLimit(value string, maximumBytes int) string {
	return RedactAndLimitWithReport(value, maximumBytes).Value
}

func RedactAndLimitWithReport(value string, maximumBytes int) RedactionResult {
	result := RedactWithReport(value)
	value = result.Value
	if maximumBytes <= 0 || len(value) <= maximumBytes {
		return result
	}
	const suffix = "\n[TRUNCATED]"
	limit := maximumBytes - len(suffix)
	if limit < 0 {
		result.Value = suffix[:maximumBytes]
		result.Truncated = true
		return result
	}
	value = value[:limit]
	for !utf8.ValidString(value) {
		value = value[:len(value)-1]
	}
	result.Value = value + suffix
	result.Truncated = true
	return result
}
