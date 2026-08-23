# Implementation brief — AKR-INVESTIGATION-LIFECYCLE

## Estado inicial

`internal/core/domain/investigation.go` está completo (estados, invariantes,
transiciones) y sin tags gorm. No existe `Operation`, ni ports, ni usecases, ni
persistencia, ni REST para ninguno de los dos. `internal/core/domain/incident.go`
está completo pero aislado: no hay tabla `incidents` ni forma real de resolver
un `incident_id`.

## Estrategia

```text
REST/Chi -> StartIncidentInvestigation -> IncidentReader (stub hasta H2)
                                        -> InvestigationStore (crea pending)
                                        -> OperationStore (crea queued)
                                        -> InvestigationDispatcher (goroutine)
                                             -> RunInvestigationUseCase
                                                  -> InvestigationRunner (stub QVAC)
                                                  -> InvestigationStore/OperationStore (persisten transición final)
```

`Operation` es dominio + puerto + persistencia genéricos, reutilizables luego
por remediation/pull_request (H5); esta tarea solo los ejercita desde
Investigation. La ejecución asíncrona corre en background con su propio
contexto (no el de la request) y siempre termina en una transición discreta
persistida (`succeeded`/`failed`), sin transacciones abiertas sobre llamadas
externas.

## Dominio

- `Investigation`: se agregan tags `gorm` a los campos existentes (columna +
  tipo, `serializer:json`+`jsonb` para los 4 slices), sin tocar validaciones ni
  transiciones.
- `Operation` (nuevo, `internal/core/domain/operation.go`): `ID, Type
  (OperationType), Status (OperationStatus), ResourceType, ResourceID *uuid.UUID,
  UserMessage, FailureCode *string, IdempotencyKey *string, CreatedAt, UpdatedAt,
  FinishedAt *time.Time`, con `NewOperation(...)`, `Start(at)`, `Succeed(at,
  message)`, `Fail(at, message, code)` y `Validate()` por estado (mismo estilo
  que `Investigation`: `queued` sin `FinishedAt`, terminal con `FinishedAt >=
  CreatedAt`). Nuevos errores de dominio en el grupo `407`.

## Ports

- `out.IncidentReader { Exists(ctx, id uuid.UUID) (bool, error) }` — mínimo
  necesario; implementación real pendiente del merge de H2.
- `out.InvestigationStore` — crear/actualizar, `FindByID`, `ListByIncident`
  (paginado, `paging.Params`/`paging.Slice`), `ExistsActiveForIncident` (para
  el 409 de investigación duplicada).
- `out.OperationStore` — crear/actualizar, `FindByID`, `FindByIdempotencyKey`
  (replay de `Idempotency-Key`).
- `out.InvestigationRunner { Run(ctx, domain.Investigation)
  (InvestigationRunResult, error) }` — hoy implementado por un stub en
  `internal/adapter/external/qvac/` que siempre devuelve error; la integración
  real con QVAC queda pendiente de PB-028+.
- `out.InvestigationDispatcher { Dispatch(ctx, investigationID, operationID
  uuid.UUID) }` — implementado en `internal/service/investigationdispatch/`
  con una goroutine simple (no worker pool/queue; documentado como mínimo
  intencional hasta que el volumen lo justifique).

## Usecases

- `internal/usecase/investigation/`: `StartIncidentInvestigation` (valida
  incident vía `IncidentReader`, rechaza duplicado activo con 409, resuelve
  replay por `Idempotency-Key` devolviendo la misma Operation, crea
  Investigation+Operation, dispara el dispatcher, devuelve la Operation),
  `GetInvestigation`, `ListIncidentInvestigations`, y `RunInvestigation`
  (orquesta la transición async: pending->running->completed/failed sobre
  Investigation y Operation, invocado por el dispatcher, testeable de forma
  síncrona con fakes).
- `internal/usecase/operation/`: `GetOperation`.

## Persistencia

Migraciones `20260822_10_add_investigations` y `20260822_11_add_operations`
(mismo formato que `..._09_add_projects`): `investigations.incident_id` sin FK
(no existe la tabla `incidents` todavía — gap documentado, se agrega en una
migración de seguimiento cuando H2 mergee); jsonb con default `'[]'::jsonb`
para los 4 arrays; `operations.idempotency_key` con índice único (nulls
múltiples permitidos); `resource_id` sin FK (polimórfico). Errores de
persistencia nuevos en `204` (investigation) y `205` (operation).

## REST

DTOs `InvestigationDTO`/`InvestigationResponseDTO`/`InvestigationListResponseDTO`
y `OperationDTO`/`OperationResponseDTO` reflejan los schemas del OpenAPI campo a
campo; `InvestigationDTO.started_at` usa el fallback a `CreatedAt` descrito
arriba. `startIncidentInvestigation` exige y valida `Idempotency-Key` (UUID,
400 si falta o es inválida) y responde 202 con `Operation`. Las 4 rutas
requieren sesión (`Admin`); la ruta POST además exige Origin permitido, como el
resto de las mutaciones (aunque el spec no liste 403 en ese path puntual, es
una política transversal ya aplicada a otras mutaciones). Errores de negocio:
`incident_id` inexistente -> 404; investigación activa duplicada -> 409;
`operation_id`/`investigation_id` inexistente -> 404.

## Wiring de producción

Se extiende `internal/bootstrap/integrations/module.go` (mismo composition
root que ya wirea Project) con los nuevos repos/usecases/handlers. El
`IncidentReader` de producción es el stub que siempre deniega (`Exists=false`):
las 4 rutas quedan vivas y auditables end-to-end, pero crear una Investigation
da 404 hasta que H2 aporte un reader real — nunca autoriza de más.
