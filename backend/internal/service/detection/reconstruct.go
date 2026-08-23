package detection

import (
	"regexp"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

var (
	continuationPattern = regexp.MustCompile(`(?i)^(?:at\s+|caused by:|suppressed:|\.\.\. \d+ more$|file \"|traceback \(|goroutine \d+ \[|created by |runtime\.|[a-z0-9_./-]+\([^)]*\)$|exception in thread)`)
	exceptionPattern    = regexp.MustCompile(`(?i)(?:panic:|traceback \(most recent call last\)|(?:exception|error)(?::|$)|unhandled rejection)`)
	framePattern        = regexp.MustCompile(`(?i)(?:\bat\s+[^ ]+\([^)]*:\d+\)|file \"[^\"]+\", line \d+|\.go:\d+|\.js:\d+|\.ts:\d+)`)
)

func Reconstruct(records []domain.SanitizedLogRecord) []LogicalEvent {
	result := make([]LogicalEvent, 0, len(records))
	for _, record := range records {
		if len(result) == 0 || !continues(result[len(result)-1], record.Message) {
			result = append(result, LogicalEvent{Records: []domain.SanitizedLogRecord{record}})
			continue
		}
		current := &result[len(result)-1]
		if len(current.Message()) < maximumLogicalEventBytes {
			current.Records = append(current.Records, record)
		}
	}
	return result
}

func continues(current LogicalEvent, next string) bool {
	next = strings.TrimSpace(next)
	if next == "" {
		return true
	}
	currentMessage := current.Message()
	if continuationPattern.MatchString(next) || framePattern.MatchString(next) {
		return exceptionPattern.MatchString(currentMessage) || framePattern.MatchString(currentMessage)
	}
	return (exceptionPattern.MatchString(currentMessage) || errorLevelPattern.MatchString(currentMessage)) && exceptionPattern.MatchString(next)
}
