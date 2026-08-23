package request

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

var ErrMissingIdempotencyKey = errors.New("missing or invalid Idempotency-Key header")

// IdempotencyKey reads and parses the required Idempotency-Key header used by
// commands that queue asynchronous work.
func IdempotencyKey(r *http.Request) (uuid.UUID, error) {
	key, err := uuid.Parse(r.Header.Get("Idempotency-Key"))
	if err != nil {
		return uuid.Nil, ErrMissingIdempotencyKey
	}
	return key, nil
}
