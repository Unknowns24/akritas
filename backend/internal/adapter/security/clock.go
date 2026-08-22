package security

import (
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type systemClock struct{}

func NewClock() out.Clock {
	return systemClock{}
}

func (systemClock) Now() time.Time {
	return time.Now().UTC()
}
