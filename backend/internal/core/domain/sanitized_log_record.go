package domain

import (
	"strings"
	"time"
)

type LogStream string

const (
	LogStreamStdout  LogStream = "stdout"
	LogStreamStderr  LogStream = "stderr"
	LogStreamUnknown LogStream = "unknown"
)

func (s LogStream) Validate() error {
	switch s {
	case LogStreamStdout, LogStreamStderr, LogStreamUnknown:
		return nil
	default:
		return ErrInvalidSanitizedLogRecord.Wrap(validationCause("log stream"))
	}
}

type SanitizedLogRecord struct {
	Timestamp time.Time
	Stream    LogStream
	Message   string
	Redacted  bool
}

func NewSanitizedLogRecord(timestamp time.Time, stream LogStream, message string) (SanitizedLogRecord, error) {
	record := SanitizedLogRecord{Timestamp: timestamp, Stream: stream, Message: strings.TrimSpace(message), Redacted: true}
	if err := record.Validate(); err != nil {
		return SanitizedLogRecord{}, err
	}
	return record, nil
}

func (r SanitizedLogRecord) Validate() error {
	if !validTime(r.Timestamp) || r.Stream.Validate() != nil || !nonBlank(r.Message) || len(r.Message) > 20000 || !r.Redacted {
		return ErrInvalidSanitizedLogRecord.Wrap(validationCause("sanitized log record"))
	}
	return nil
}
