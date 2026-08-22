# Backend Background Processes Policy

Akritas contains long-running and retryable workflows such as log collection, detection, investigation and remediation.

## Scope

This policy applies to backend-service tasks including:

- Dokploy log polling/collection;
- detection and incident grouping;
- QVAC investigations;
- GitHub issue publication;
- remediation workflows;
- branch/change/PR orchestration.

## Idempotency

Operations that may be retried must be designed to avoid duplicate externally visible effects.

Examples:

- do not create multiple incidents for the same grouping decision merely because a poll is retried;
- do not create duplicate GitHub issues for the same publication attempt when an idempotency/deduplication mechanism exists in the application design;
- do not create multiple remediation branches/PRs for the same successful workflow unless explicitly requested.

Idempotency keys, fingerprints and persisted execution state must follow documented domain decisions; do not invent new domain concepts silently.

## Transaction boundaries

Do not keep database transactions open while performing slow or unreliable external calls to GitHub, Dokploy, QVAC, Git or the filesystem.

Prefer:

```text
read/persist state
→ close transaction
→ call external system
→ persist resulting state
```

Use a stronger atomicity mechanism only when the architecture explicitly requires one.

## Retryable orchestration

Separate orchestration state from infrastructure implementation so a failed external step can be retried safely when the product design permits it.

Do not make correctness depend on an HTTP request remaining alive.

## Cancellation and timeouts

External operations must accept context/cancellation and use bounded timeouts according to existing project conventions.

## Concurrency

When multiple workers/processes can observe the same incident or project, protect state transitions using the persistence/concurrency mechanism chosen by the project. Do not rely on in-memory mutexes for cross-process correctness.

## Observability

Processes should emit enough structured context to diagnose failures without logging secrets, credentials or full sensitive payloads.
