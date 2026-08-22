# Backend Modularity and SRP Policy

Backend code MUST follow strict Single Responsibility Principle (SRP) and feature-oriented modularity.

The purpose of this policy is to prevent horizontal packages, monolithic files, and generic modules that accumulate unrelated operations as the codebase grows.

## General rules

Each package MUST represent one cohesive responsibility or feature.

Each file MUST have one clear responsibility.

When implementing:

- a usecase operation;
- a repository operation;
- a service operation;
- a REST handler operation;
- a port operation;

there MUST be at most one public operation per file.

Struct definitions, constructors, shared interfaces, and narrowly scoped helpers MAY live in dedicated files.

Do NOT group multiple CRUD or business operations into a single implementation file.

Forbidden examples:

```text
handlers.go
campaign_handlers.go
repositories.go
campaign_repository.go   # containing every repository operation
campaign_service.go      # containing every service operation
usecases.go
crud.go
```

when those files contain multiple independent operations.

## Feature-oriented packages

Packages MUST be organized around a cohesive feature or responsibility.

Do NOT use horizontal packages as dumping grounds for unrelated implementations.

For example, REST handlers MUST NOT be structured like:

```text
internal/adapter/rest/handler/
  create_campaign.go
  update_campaign.go
  get_company.go
  create_lead.go
  list_leads.go
```

Instead, handlers MUST be grouped by feature:

```text
internal/adapter/rest/handler/
  campaign/
    handler.go
    create.go
    get.go
    list.go
    update.go

  company/
    handler.go
    get.go
    list.go

  lead/
    handler.go
    create.go
    get.go
    list.go
```

The same principle applies to repositories, services, usecases, and other feature implementations.

## Usecase packages

Each usecase MUST live inside its feature package:

```text
internal/usecase/<feature>/
```

A feature with multiple operations MUST be split by operation:

```text
internal/usecase/campaign/
  uc.go
  create.go
  get.go
  list.go
  update.go
  activate.go
  pause.go
```

`uc.go` MUST contain only the usecase structure, constructor, and package-level wiring directly related to constructing the usecase.

Example:

```go
package campaign

import "project/internal/core/ports/in"

type CampaignUseCase struct {
    // dependencies
}

func NewCampaignUseCase(/* deps */) in.CampaignUseCase {
    return &CampaignUseCase{}
}
```

Each operation MUST be implemented in its own file:

```go
// create.go
func (uc *CampaignUseCase) Create(/* params */) (/* returns */) {
    // implementation
}
```

```go
// activate.go
func (uc *CampaignUseCase) Activate(/* params */) error {
    // implementation
}
```

A file such as:

```text
campaign.go
```

MUST NOT contain `Create`, `Get`, `List`, `Update`, `Activate`, `Pause`, and other independent operations together.

## Repository packages

Repositories MUST be organized by feature:

```text
internal/adapter/db/<technology>/repository/<feature>/
```

Example:

```text
internal/adapter/db/postgres/repository/campaign/
  repo.go
  create.go
  get_by_id.go
  list.go
  update.go
  delete.go
```

`repo.go` MUST contain only:

- the repository struct;
- its constructor;
- dependencies required by the repository;
- narrowly scoped package-level declarations required for construction.

Each repository operation MUST live in its own file.

Example:

```go
// get_by_id.go
func (r *Repository) GetByID(/* params */) (/* returns */) {
    // implementation
}
```

Repository packages MUST NOT become generic database access packages containing operations from multiple domain features.

Forbidden:

```text
repository/
  repository.go
```

with methods for campaigns, companies, leads, contacts, and unrelated entities.

## REST handlers

REST handlers MUST be grouped by feature or cohesive operational concern:

```text
internal/adapter/rest/handler/<feature>/
```

Example:

```text
internal/adapter/rest/handler/campaign/
  handler.go
  create.go
  get.go
  list.go
  update.go
  activate.go
  pause.go
```

`handler.go` MUST contain only:

- the handler struct;
- its constructor;
- direct dependencies needed by the handler.

Each HTTP operation MUST live in its own file.

For example:

```go
// create.go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    // implementation
}
```

Operational endpoints MAY share one cohesive package when they belong to the same concern.

For example:

```text
internal/adapter/rest/handler/operational/
  handler.go
  health.go
  readiness.go
```

The root `handler` package MUST NOT contain unrelated endpoint implementations directly.

Therefore this is forbidden:

```text
internal/adapter/rest/handler/
  health.go
  readiness.go
  create_campaign.go
  get_campaign.go
  create_lead.go
```

## Services

Services MUST follow the same feature-oriented and operation-oriented structure:

```text
internal/service/<feature>/
  service.go
  operation_a.go
  operation_b.go
```

`service.go` MUST contain the service struct and constructor.

Independent behaviors MUST be split into separate files.

Services MUST NOT be used as generic containers for unrelated business logic.

## Ports

Ports MUST remain cohesive and domain-oriented.

A port MUST NOT become a generic interface containing unrelated capabilities.

When a port grows to represent multiple independent responsibilities, it MUST be split.

Example:

Preferred:

```go
type CampaignReader interface {
    GetByID(...)
    List(...)
}

type CampaignWriter interface {
    Create(...)
    Update(...)
}
```

instead of an unrelated catch-all interface.

When port implementations contain multiple methods, each implemented operation MUST still live in its own implementation file.

## Domain

Domain code MUST be organized around explicit domain concepts.

Domain entities, value objects, domain errors, and domain behavior MUST NOT be hidden inside adapters, handlers, repositories, or bootstrap code.

Examples:

```text
internal/core/domain/
  campaign.go
  company.go
  lead.go
  contact.go
  signal.go
  opportunity.go
```

A feature requiring a new domain concept MUST define that concept in the domain layer before adapters begin representing it through persistence or transport-specific structures.

HTTP DTOs MUST NOT replace domain entities. Duplicate database model structs that mirror domain fields MUST NOT replace domain entities either: persisted aggregates are the domain types (with `gorm` tags). REST DTOs remain required at the HTTP boundary.

## Bootstrap and wiring

Bootstrap code MUST only compose and start the application.

It MAY:

- load already-resolved dependencies;
- construct adapters;
- construct repositories;
- construct services;
- construct usecases;
- construct routers/servers;
- coordinate application startup and shutdown.

Bootstrap code MUST NOT contain:

- domain logic;
- repository logic;
- HTTP handler logic;
- persistence queries;
- feature-specific business rules.

When bootstrap lifecycle logic becomes independently meaningful, it MUST be split into focused files.

Example:

```text
internal/bootstrap/
  bootstrap.go
  dependencies.go
  run.go
  shutdown.go
```

A single bootstrap file MUST NOT grow indefinitely merely because all code relates to application startup.

## Shared helpers

Small helper files are allowed only when they represent a cohesive technical responsibility.

Examples:

```text
mapper.go
filters.go
errors.go
validator.go
response.go
```

A helper file MUST NOT become a dumping ground for unrelated functions.

Generic names such as:

```text
utils.go
common.go
helpers.go
misc.go
```

SHOULD be avoided.

If helpers belong only to one operation or concept, keep them close to that operation instead of creating a shared helper file prematurely.

## File naming

Go file names MUST use `snake_case.go`.

Operation files MUST describe their operation.

Examples:

```text
create.go
get_by_id.go
list.go
activate.go
pause.go
calculate_score.go
mark_as_contacted.go
```

Avoid ambiguous names such as:

```text
logic.go
methods.go
operations.go
manager.go
stuff.go
```

unless the name represents a precise and cohesive concept.

## Tests

Tests MUST remain close to the package under test.

Test files SHOULD mirror the operation being tested when practical:

```text
create.go
create_test.go

activate.go
activate_test.go
```

Large generic test files containing unrelated feature behaviors SHOULD be split.

## Review requirements

Architecture review MUST treat violations of the mandatory rules in this policy as blocking findings.

The reviewer MUST explicitly inspect:

1. whether handlers are grouped by feature or cohesive concern;
2. whether repository implementations are grouped by feature;
3. whether usecases and services are grouped by feature;
4. whether files contain multiple independent public operations;
5. whether generic horizontal packages are accumulating unrelated behavior;
6. whether domain concepts have incorrectly leaked into adapters;
7. whether bootstrap contains feature or business logic;
8. whether helper files are becoming dumping grounds.

Passing dependency-direction checks alone is NOT sufficient to pass architecture review.

A solution that respects layer imports but violates SRP or feature modularity MUST fail architecture review.

## Exceptions

Exceptions are allowed only when splitting code would reduce clarity without improving responsibility boundaries.

An exception MUST satisfy all of the following:

- the code represents one cohesive responsibility;
- the file does not contain independent business or CRUD operations;
- the exception does not create a horizontal dumping ground;
- the reason is obvious from the code or documented in the task implementation summary.

File size alone is NOT the deciding factor.

The deciding factor is whether the module has one clear reason to change.
