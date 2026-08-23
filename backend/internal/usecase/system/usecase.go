package system

import (
	"context"
	"errors"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

var ErrInvalidUseCase = errors.New("invalid system use case configuration")

type qvacStatusReader interface {
	GetStatus(context.Context) (portsin.QvacRuntimeStatus, error)
	GetConfiguration(context.Context) (domain.QvacConfiguration, error)
}

type UseCase struct {
	projections portsout.ProjectionStore
	operations  portsout.OperationStore
	qvac        qvacStatusReader
	newID       func() uuid.UUID
	now         func() time.Time
}

func New(projections portsout.ProjectionStore, operations portsout.OperationStore, qvac qvacStatusReader, newID func() uuid.UUID, now func() time.Time) (portsin.SystemUseCase, error) {
	if projections == nil || operations == nil || qvac == nil || newID == nil || now == nil {
		return nil, ErrInvalidUseCase
	}
	return &UseCase{projections: projections, operations: operations, qvac: qvac, newID: newID, now: now}, nil
}

func (uc *UseCase) GetStatus(ctx context.Context) (domain.SystemStatus, error) {
	githubAccounts, err := uc.projections.CountGitHubAccounts(ctx)
	if err != nil {
		return domain.SystemStatus{}, err
	}
	dokployServers, err := uc.projections.CountDokployServers(ctx)
	if err != nil {
		return domain.SystemStatus{}, err
	}
	config, err := uc.qvac.GetConfiguration(ctx)
	if err != nil {
		return domain.SystemStatus{}, err
	}
	qvacStatus, err := uc.qvac.GetStatus(ctx)
	if err != nil {
		return domain.SystemStatus{}, err
	}
	componentStatus := domain.ComponentHealthUnavailable
	if qvacStatus.ConnectionStatus == domain.IntegrationStatusConnected {
		componentStatus = domain.ComponentHealthHealthy
	}
	checkedAt := uc.now().UTC()
	lastDiagnostics, err := uc.projections.FindLastSystemDiagnostics(ctx)
	if err != nil {
		return domain.SystemStatus{}, err
	}
	var lastDiagnosticsAt *time.Time
	if lastDiagnostics != nil {
		lastDiagnosticsAt = &lastDiagnostics.UpdatedAt
	}
	return domain.SystemStatus{
		GitHubAccountCount: githubAccounts,
		DokployServerCount: dokployServers,
		QvacEndpoint:       config.EndpointURL,
		Components: []domain.ComponentHealth{
			{Component: "github", Status: configuredStatus(githubAccounts), CheckedAt: &checkedAt},
			{Component: "dokploy", Status: configuredStatus(dokployServers), CheckedAt: &checkedAt},
			{Component: "qvac", Status: componentStatus, CheckedAt: &checkedAt},
			{Component: "investigator", Status: componentStatus, CheckedAt: &checkedAt},
		},
		LastDiagnosticsAt: lastDiagnosticsAt,
	}, nil
}

func (uc *UseCase) RunDiagnostics(ctx context.Context, idempotencyKey uuid.UUID) (*domain.Operation, error) {
	key := idempotencyKey.String()
	existing, err := uc.operations.FindByIdempotencyKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	now := uc.now().UTC()
	resourceType := domain.OperationResourceSystem
	operation, err := domain.NewOperation(uc.newID(), domain.OperationTypeSystemDiagnostics, &resourceType, &idempotencyKey, &key, "El diagnóstico del sistema fue encolado.", now)
	if err != nil {
		return nil, err
	}
	if err := uc.operations.Create(ctx, operation); err != nil {
		return nil, err
	}
	go uc.executeDiagnostics(context.Background(), operation.ID)
	return operation, nil
}

func (uc *UseCase) executeDiagnostics(ctx context.Context, operationID uuid.UUID) {
	operation, err := uc.operations.FindByID(ctx, operationID)
	if err != nil {
		return
	}
	if err := operation.Start(uc.now().UTC()); err != nil {
		return
	}
	_ = uc.operations.Update(ctx, operation)
	status, err := uc.GetStatus(ctx)
	if err != nil {
		code := domain.ErrIntegrationUnavailable.Code
		_ = operation.Fail(uc.now().UTC(), "No se pudieron completar los diagnósticos.", &code)
		_ = uc.operations.Update(ctx, operation)
		return
	}
	message := "Diagnóstico completado."
	for _, component := range status.Components {
		if component.Status == domain.ComponentHealthUnavailable {
			message = "Diagnóstico completado con componentes no disponibles."
			break
		}
	}
	_ = operation.Succeed(uc.now().UTC(), message)
	_ = uc.operations.Update(ctx, operation)
}

func configuredStatus(count int) domain.ComponentHealthStatus {
	if count > 0 {
		return domain.ComponentHealthHealthy
	}
	return domain.ComponentHealthNotConfigured
}
