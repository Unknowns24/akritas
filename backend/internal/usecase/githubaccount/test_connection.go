package githubaccount

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/google/uuid"
)

func (uc *UseCase) TestConnection(ctx context.Context, id uuid.UUID) (portsin.ConnectionTestResult, error) {
	account, err := uc.store.Get(ctx, id)
	if err != nil {
		return portsin.ConnectionTestResult{}, err
	}
	result, err := uc.gateway.TestConnection(ctx, *account)
	if err != nil {
		return portsin.ConnectionTestResult{}, err
	}
	checkedAt := result.CheckedAt.UTC()
	if checkedAt.IsZero() {
		checkedAt = uc.now().UTC()
	}
	account.LastCheckedAt = &checkedAt
	account.UpdatedAt = checkedAt
	switch result.Status {
	case domain.ConnectionTestConnected:
		account.AuthenticationStatus = domain.IntegrationStatusConnected
	case domain.ConnectionTestAuthenticationFailed:
		account.AuthenticationStatus = domain.IntegrationStatusAuthenticationFailed
	default:
		account.AuthenticationStatus = domain.IntegrationStatusUnavailable
	}
	if err := uc.store.UpdateConnection(ctx, account); err != nil {
		return portsin.ConnectionTestResult{}, err
	}
	latency := result.Latency.Milliseconds()
	if latency < 0 {
		latency = 0
	}
	return portsin.ConnectionTestResult{Status: result.Status, CheckedAt: checkedAt.Format(time.RFC3339Nano), LatencyMS: &latency, UserMessage: result.UserMessage}, nil
}
