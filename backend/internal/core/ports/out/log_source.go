package out

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type RawLogRecord struct {
	Timestamp   time.Time
	Ordinal     int
	ContentHash string
	Message     string
}

type LogFetchRequest struct {
	Server domain.DokployServer
	Source domain.DokploySource
	Tail   int
	Since  string
}

type LogSource interface {
	FetchLogs(context.Context, LogFetchRequest) ([]RawLogRecord, error)
}
