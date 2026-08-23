# Final summary — AKR-INVESTIGATION-LIFECYCLE

PB-026 quedó implementado sobre `feat/investigation-with-qvac`: los 4
endpoints de Investigation/Operation, el dominio `Operation` nuevo, y el
pipeline asíncrono completo (crear → encolar → poll), con la frontera hacia
H2 y PB-028+ resuelta con puertos + stubs explícitos en vez de dejarla
implícita.

## Qué está implementado (código real, no fake)

- Dominio `Operation`, tags GORM en `Investigation`, todos los ports
  (`IncidentReader`, `InvestigationStore`, `OperationStore`,
  `InvestigationRunner`, `InvestigationDispatcher`), usecases
  (`StartIncidentInvestigation`, `GetInvestigation`,
  `ListIncidentInvestigations`, `RunInvestigation`, `GetOperation`),
  persistencia PostgreSQL real (migraciones + repositorios), REST completo
  (DTOs, mapper, handlers, rutas, `Idempotency-Key`), y el dispatcher
  asíncrono (`internal/service/investigationdispatch`) — todo esto corre tal
  cual en producción.

## Qué está probado con fake/stub (documentado, no un gap oculto)

- `out.IncidentReader`: los usecases se probaron con un fake controlable
  (`Exists` true/false/error). La implementación real (PostgreSQL, contra la
  tabla `incidents` de H2) queda pendiente de que H2 mergee. En producción se
  conecta `stub.DenyAllIncidentReader` (siempre `Exists=false`), verificado
  manualmente contra un server real: las 4 rutas responden 401/403/404 según
  corresponda, nunca autorizan de más.
- `out.InvestigationRunner`: se probó con un fake controlable
  (éxito/fallo) en `RunInvestigation`. En producción corre
  `qvac.StubRunner`, que siempre falla con un mensaje explícito
  ("QVAC integration is not implemented yet; see PB-028+"), marcando la
  Investigation y la Operation como `failed`. La integración real con QVAC
  (llamadas al modelo, armado de evidencia, tool calling, clasificación de
  root_cause/resolution) es PB-028 a PB-035, explícitamente fuera de
  alcance de esta tarea.

## Decisiones de producto aplicadas (aprobadas por vos)

- `InvestigationDTO.started_at` usa `CreatedAt` como fallback mientras el
  estado es `pending`; el dominio no cambió.
- `IncidentReader` de producción es un stub deny-by-default, no rutas sin
  registrar: las 4 rutas quedan vivas y auditables end-to-end aunque crear
  una Investigation siempre 404 hasta que H2 aporte un reader real.

## Riesgos remanentes / gaps documentados

- `investigations.incident_id` no tiene FK a una tabla `incidents` (no
  existe todavía); se agrega en una migración de seguimiento cuando H2
  mergee.
- `InvestigationDispatcher` es una goroutine simple sin reintentos ni
  persistencia de cola: si el proceso se reinicia mientras una Investigation
  está `running`, queda huérfana (ni completa ni recuperada automáticamente).
  Documentado como mínimo intencional hasta que el volumen de PB-028+
  justifique un worker real.

## Hallazgo no relacionado, no corregido (fuera de alcance)

`TestCreatePersistsAdministrator` y `TestSavePersistsAdministratorSession`
(paquetes `repository/administrator` y `repository/administrator_session`)
fallan en este branch de forma preexistente: un `db.First(&record, ...)`
sobre un struct anónimo no resuelve nombre de tabla. Confirmado con
`git stash` que el fallo es idéntico sin ningún cambio de esta tarea — no se
tocó ningún archivo de esos paquetes. Se deja señalado para quien lleve auth,
no se corrigió por estar fuera del alcance de PB-026.

## Validaciones

- `go build ./...`: correcto.
- `go vet ./...`: correcto.
- `gofmt -l .`: sin diffs.
- `go test ./...`: correcto, salvo los dos tests preexistentes y no
  relacionados arriba.
- `check-backend-architecture.sh`: correcto.
- `check-openapi.sh`: correcto, OpenAPI 1.4.0 sin cambios, 60 operaciones y
  112 schemas (el contrato ya publicado se implementó tal cual).
- `check-security.sh`: correcto.
- Prueba manual end-to-end contra Postgres local real (`akritas_e2e`):
  bootstrap → TOTP verify → sesión → `listIncidentInvestigations` (404),
  `startIncidentInvestigation` sin `Idempotency-Key` (400), sin Origin (403),
  con ambos (404 vía el stub), `getInvestigation`/`getOperation` con ID
  inexistente (404), `getOperation` sin sesión (401).

No quedan validaciones pendientes ni excepciones al harness para el alcance
de esta tarea.
