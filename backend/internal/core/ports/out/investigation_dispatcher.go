package out

import "github.com/google/uuid"

// InvestigationDispatcher queues an already-persisted, pending Investigation
// for asynchronous execution. Dispatch does not block the caller, does not
// propagate execution errors back to the HTTP request that triggered it, and
// deliberately takes no context: the work it starts must outlive the request
// that returned 202, so the implementation owns its own bounded lifetime
// instead of inheriting the request's. Callers observe progress by polling
// the returned Operation.
type InvestigationDispatcher interface {
	Dispatch(investigationID, operationID uuid.UUID)
}
