# Implementation Brief — AKR-H3-FINAL-INTEGRATION

## Task

Finalizar H3 sobre la implementación H2/H3 combinada y restaurar una arquitectura coherente basada en los repositorios reales de Incident y LogEvent.

## Current project context

El catálogo PostgreSQL contiene dos `Catalog()` superpuestos. H3 define un `IncidentReader` temporal con `Exists`, mientras H2 ya ofrece `IncidentStore.Get/List/ListLogEvents`. Evidence sólo contiene metadata de deployment; QVAC recibe count, no contenido; las tools están acopladas a DTOs de adapters; Investigation no tiene FK real a Incident ni registra Evidence citada.

## Proposed approach

Consolidar errores; dividir el port H2 en capacidades pequeñas; eliminar `IncidentReader`; hacer start/run/recovery transaccionales; ensamblar Evidence real y redactada desde LogEvents/Project; entregar contexto bounded a QVAC; expresar inspección GitHub mediante un port read-only; persistir Evidence de tools y citations; endurecer el decoder estructurado; añadir migraciones aditivas y actualizar OpenAPI.

## Architecture impact

Se mantiene `adapter → usecase/service → core ports/domain`. QVAC no consulta PostgreSQL. El repositorio GitHub queda detrás de `RepositoryInspector`. La transición a `publishing_issue` queda exclusivamente en H4.

## API/OpenAPI impact

Sin rutas nuevas. OpenAPI 1.6.0 agrega `created_at` y `evidence_ids` a Investigation y hace `started_at` opcional mientras está pending.

## Data/persistence impact

Dos migraciones aditivas: FKs RESTRICT + índice único parcial activo, y `evidence_ids JSONB NOT NULL DEFAULT '[]'`. No se reescriben migraciones aplicadas.

## Error handling impact

Se conserva la numeración estable. `ErrIncidentNotFound` se cataloga una sola vez. Los fallos técnicos persisten códigos/mensajes públicos y llevan Investigation/Operation/Incident a failed; `root_cause_status=unknown` válido completa normalmente.

## Test strategy

Pruebas unitarias de dominio/usecase/adapters/REST, pruebas de migración/repositorios con PostgreSQL real y un escenario H2→H3 completo con fake QVAC, además de todos los gates y regresiones H1/H2.

## Risks

Carreras al iniciar Investigation, pérdida de Evidence ante fallo externo, prompt injection, fuga de credenciales, overflow de contexto QVAC, alcance cruzado de repositorios y compatibilidad de datos preexistentes con las nuevas restricciones.
