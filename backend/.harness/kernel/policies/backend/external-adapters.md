# Backend External Adapters Policy

Akritas integrates with external systems that must remain infrastructure concerns.

## External systems

Examples include:

- GitHub;
- Dokploy;
- QVAC/local model server;
- Git repositories and worktrees;
- filesystem access;
- process execution such as the Git CLI;
- any future external API or SDK.

## Boundary rule

Core/domain/usecase code must depend on capabilities expressed as output ports, never on concrete infrastructure clients.

Preferred dependency direction:

```text
core/usecase
    ↓
core/ports/out
    ↓
adapter/external/*
    ↓
GitHub / Dokploy / QVAC / Git / filesystem
```

Forbidden examples in core/usecase:

- importing a GitHub SDK;
- constructing an `http.Client` specifically for Dokploy or QVAC;
- invoking `exec.Command("git", ...)`;
- opening or mutating repository files through concrete filesystem logic;
- embedding provider-specific DTOs into domain entities.

## Ports

Ports should describe business capabilities, not vendor APIs.

Prefer capability-oriented interfaces such as:

```text
LogSource
InvestigationEngine
IssuePublisher
RepositoryWorkspace
ChangePublisher
```

Use provider-specific names only when the capability is inherently provider-specific and the domain/documentation explicitly models it that way.

## Mapping

Provider DTOs and API payloads belong in adapters. Translate them at the boundary into domain/application types.

Do not leak SDK errors directly into domain/application behavior. Normalize them into project errors according to the domain error policy.

## Testability

Usecases must be testable with fakes/stubs for output ports without live external services.
