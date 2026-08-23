package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type OperationType string

const (
	OperationTypeSystemDiagnostics OperationType = "system_diagnostics"
	OperationTypeInvestigation     OperationType = "investigation"
	OperationTypeRemediation       OperationType = "remediation"
	OperationTypePullRequest       OperationType = "pull_request"
)

func (t OperationType) Validate() error {
	switch t {
	case OperationTypeSystemDiagnostics, OperationTypeInvestigation, OperationTypeRemediation, OperationTypePullRequest:
		return nil
	default:
		return ErrInvalidOperationType.Wrap(validationCause("operation type"))
	}
}

type OperationStatus string

const (
	OperationStatusQueued    OperationStatus = "queued"
	OperationStatusRunning   OperationStatus = "running"
	OperationStatusSucceeded OperationStatus = "succeeded"
	OperationStatusFailed    OperationStatus = "failed"
)

func (s OperationStatus) Validate() error {
	switch s {
	case OperationStatusQueued, OperationStatusRunning, OperationStatusSucceeded, OperationStatusFailed:
		return nil
	default:
		return ErrInvalidOperationStatus.Wrap(validationCause("operation status"))
	}
}

type OperationResourceType string

const (
	OperationResourceSystem        OperationResourceType = "system"
	OperationResourceInvestigation OperationResourceType = "investigation"
	OperationResourceRemediation   OperationResourceType = "remediation"
	OperationResourcePullRequest   OperationResourceType = "pull_request"
)

func (t OperationResourceType) Validate() error {
	switch t {
	case OperationResourceSystem, OperationResourceInvestigation, OperationResourceRemediation, OperationResourcePullRequest:
		return nil
	default:
		return ErrInvalidOperationResourceType.Wrap(validationCause("operation resource type"))
	}
}

// Operation is generic async-command infrastructure: it tracks a queued
// command through completion regardless of which resource it acts on, so
// investigation, remediation and pull_request flows can share it.
type Operation struct {
	ID             uuid.UUID              `gorm:"column:id;type:uuid;primaryKey"`
	Type           OperationType          `gorm:"column:type"`
	Status         OperationStatus        `gorm:"column:status"`
	ResourceType   *OperationResourceType `gorm:"column:resource_type"`
	ResourceID     *uuid.UUID             `gorm:"column:resource_id;type:uuid"`
	UserMessage    string                 `gorm:"column:user_message"`
	FailureCode    *string                `gorm:"column:failure_code"`
	IdempotencyKey *string                `gorm:"column:idempotency_key"`
	CreatedAt      time.Time              `gorm:"column:created_at"`
	UpdatedAt      time.Time              `gorm:"column:updated_at"`
	FinishedAt     *time.Time             `gorm:"column:finished_at"`
}

func NewOperation(
	id uuid.UUID,
	operationType OperationType,
	resourceType *OperationResourceType,
	resourceID *uuid.UUID,
	idempotencyKey *string,
	userMessage string,
	createdAt time.Time,
) (*Operation, error) {
	operation := &Operation{
		ID: id, Type: operationType, Status: OperationStatusQueued,
		ResourceType: resourceType, ResourceID: resourceID, IdempotencyKey: idempotencyKey,
		UserMessage: strings.TrimSpace(userMessage), CreatedAt: createdAt, UpdatedAt: createdAt,
	}
	if err := operation.Validate(); err != nil {
		return nil, err
	}
	return operation, nil
}

func (o Operation) Validate() error {
	if o.ID == uuid.Nil || o.Type.Validate() != nil || o.Status.Validate() != nil || !validTime(o.CreatedAt) || o.UpdatedAt.Before(o.CreatedAt) {
		return ErrInvalidOperation.Wrap(validationCause("operation"))
	}
	if o.ResourceType != nil && o.ResourceType.Validate() != nil {
		return ErrInvalidOperation.Wrap(validationCause("operation resource type"))
	}
	if (o.ResourceType == nil) != (o.ResourceID == nil) {
		return ErrInvalidOperation.Wrap(validationCause("operation resource pair"))
	}
	switch o.Status {
	case OperationStatusQueued, OperationStatusRunning:
		if o.FinishedAt != nil {
			return ErrInvalidOperation.Wrap(validationCause("open operation times"))
		}
	case OperationStatusSucceeded, OperationStatusFailed:
		if o.FinishedAt == nil || o.FinishedAt.Before(o.CreatedAt) {
			return ErrInvalidOperation.Wrap(validationCause("terminal operation times"))
		}
		if o.Status == OperationStatusFailed && !nonBlank(o.UserMessage) {
			return ErrInvalidOperation.Wrap(validationCause("failed operation message"))
		}
	}
	return nil
}

func (o *Operation) Start(at time.Time) error {
	if o == nil || o.Status != OperationStatusQueued || at.Before(o.CreatedAt) {
		return ErrOperationTransition.Wrap(validationCause("start operation"))
	}
	o.Status = OperationStatusRunning
	o.UpdatedAt = at
	return nil
}

func (o *Operation) Succeed(at time.Time, userMessage string) error {
	if o == nil || o.Status != OperationStatusRunning || at.Before(o.CreatedAt) {
		return ErrOperationTransition.Wrap(validationCause("succeed operation"))
	}
	o.Status = OperationStatusSucceeded
	o.UserMessage = strings.TrimSpace(userMessage)
	o.UpdatedAt = at
	o.FinishedAt = &at
	return nil
}

func (o *Operation) Fail(at time.Time, userMessage string, failureCode *string) error {
	if o == nil || o.Status != OperationStatusRunning || at.Before(o.CreatedAt) || !nonBlank(userMessage) {
		return ErrOperationTransition.Wrap(validationCause("fail operation"))
	}
	o.Status = OperationStatusFailed
	o.UserMessage = strings.TrimSpace(userMessage)
	o.FailureCode = failureCode
	o.UpdatedAt = at
	o.FinishedAt = &at
	return nil
}
