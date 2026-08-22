# Testing Policy

Harness tasks should use tests as the implementation contract.

## Flow

1. Generate test plan.
2. Wait for human approval.
3. Implement code to satisfy approved tests.
4. Run tests and validation scripts.
5. Summarize results.

## Backend

Prefer tests at the lowest useful level:

- Domain tests for domain rules.
- Usecase tests for business flow.
- Handler tests for transport mapping, auth, body parsing and error response behavior.
- Repository tests only when persistence behavior is non-trivial or already covered in the module.

## Frontend

When tests exist, prefer:

- Schema validation tests.
- Service/client tests.
- Component/view tests for important user flows.
- Typecheck and lint are mandatory validation gates.

## Human approval

Do not implement code after generating TDD tests until the human explicitly approves the plan.
