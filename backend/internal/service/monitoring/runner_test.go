package monitoring

import (
	"errors"
	"testing"
	"time"
)

func TestRunnerEnforcesConfiguredConcurrencyBoundary(t *testing.T) {
	service := &Service{}
	if runner, err := NewRunner(service, 10*time.Second, 4); err != nil || runner == nil {
		t.Fatalf("valid runner = %+v, %v", runner, err)
	}
	if runner, err := NewRunner(service, 10*time.Second, 5); runner != nil || !errors.Is(err, ErrInvalidRunnerConfiguration) {
		t.Fatalf("excess concurrency = %+v, %v", runner, err)
	}
}
