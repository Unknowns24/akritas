# TDD test plan — AKR-INVESTIGATION-LIFECYCLE

Plan sujeto a aprobación explícita antes de escribir tests o implementación.

## Dominio

- `Operation`: validar transiciones `queued->running->succeeded/failed`,
  invariantes de tiempo (`FinishedAt >= CreatedAt`, terminal requiere
  `FinishedAt`), estados/tipos inválidos.
- `Investigation`: sin tests nuevos de comportamiento (el dominio no cambia);
  solo confirmar que las tags gorm no rompen los tests existentes.

## Ports / fakes

- fake `IncidentReader` controlable (`Exists` true/false/error) para usecases.
- fake `InvestigationRunner` controlable (éxito con `InvestigationRunResult`
  completo, y error) para `RunInvestigation`.

## Usecases — investigation

- `StartIncidentInvestigation`: incident inexistente -> 404; investigación
  activa ya existente -> 409; `Idempotency-Key` repetida -> devuelve la misma
  Operation sin crear una nueva; happy path crea Investigation `pending` +
  Operation `queued` y dispara el dispatcher exactamente una vez.
- `GetInvestigation`/`ListIncidentInvestigations`: not-found, paginación
  (cursor/límite/orden), scoping por `incident_id`.
- `RunInvestigation`: runner exitoso completa Investigation+Operation con los
  mismos datos; runner con error falla ambas con mensaje explícito; el estado
  intermedio `running` se persiste antes de invocar al runner.

## Usecases — operation

- `GetOperation`: happy path y not-found.

## Servicio — investigationdispatch

- `Dispatch` invoca `RunInvestigation` de forma asíncrona (no bloquea al
  caller) y no propaga sus errores al flujo síncrono.

## PostgreSQL

- migraciones y rollback de `investigations`/`operations`;
- persistencia completa de slices jsonb y punteros nullable;
- `ExistsActiveForIncident` y `ListByIncident` con paginación real;
- `FindByIdempotencyKey` distingue hit/miss;
- errores not-found y de persistencia normalizados.

## REST

- inventario exacto de las 4 rutas, envelopes y código 202 con header
  `Retry-After` en `startIncidentInvestigation`;
- `Idempotency-Key` ausente o inválida -> 400;
- 401 sin sesión, 403 en POST sin Origin permitido, 404 incident/investigation/
  operation inexistente, 409 investigación duplicada;
- `InvestigationDTO.started_at` presente siempre (fallback a `created_at` en
  `pending`), slices nunca `null`;
- stub de `IncidentReader` en el router de producción: `startIncidentInvestigation`
  siempre 404 hasta que H2 aporte un reader real.

## Validación final

- `go test ./...`
- `go vet ./...`
- `gofmt -l .` (sin diffs)
- `.harness/kernel/scripts/check-backend-architecture.sh`
- `.harness/kernel/scripts/check-openapi.sh`
- `.harness/kernel/scripts/check-security.sh`
- prueba manual end-to-end contra Postgres local: `startIncidentInvestigation`
  contra un `incident_id` inexistente (404 esperado, dado el stub de
  producción), y poll de `getOperation` hasta estado terminal usando el fake
  de test.
