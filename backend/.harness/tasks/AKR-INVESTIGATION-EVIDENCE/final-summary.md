# Final summary — AKR-INVESTIGATION-EVIDENCE

PB-027 quedó implementado sobre `feat/investigation-with-qvac`: Evidence
tiene persistencia real, listado REST paginado, y un tipo (`deployment_metadata`)
se ensambla de verdad dentro del pipeline asíncrono de PB-026, sin
placeholders para los tipos que todavía no tienen fuente real.

## Qué está implementado (código real, no fake)

- Tags GORM en `Evidence`, extensión de `out.IncidentReader` con `Get`,
  ports nuevos (`EvidenceStore`, `EvidenceAssembler`), usecase
  `ListInvestigationEvidence`, extensión de `RunUseCase.Execute` para
  ensamblar/persistir evidencia y actualizar `EvidenceCount` real,
  persistencia PostgreSQL (migración + repositorio con FK real a
  `investigations`), REST completo (DTO, mapper, handler, ruta), y el
  servicio `evidenceassembly` que arma `deployment_metadata` a partir de
  `Project` — todo esto corre tal cual en producción.

## Qué evidencia quedó real vs. pendiente (no es un gap oculto)

- **Real, ensamblada de verdad**: `deployment_metadata`, a partir de
  `Investigation.IncidentID` → `Incident.ProjectID` → `Project` (nombre,
  monitoring/health status, snapshot de GitHubRepository/DokployApplication
  — nunca credenciales).
- **Pendiente de H2**: `log_excerpt`, `stack_trace` — necesitan
  `LogEvent`/`Incident` real, que hoy solo existen como dominio.
- **Pendiente de PB-030/PB-031**: `code_location`, `commit`, `diff` —
  necesitan lectura real de repositorio/commits, no implementada en esta
  rama.
- **Pendiente de H5 (Remediation)**: `validation_result` — pertenece a un
  dominio que esta tarea no toca.

El `EvidenceAssembler` simplemente no produce Evidence de estos cuatro
tipos; no hay datos inventados en ningún caso.

## Decisiones de diseño aplicadas

- Un fallo de ensamblado por "no encontrado" (incidente o proyecto
  inexistente) no aborta la investigación: devuelve slice vacío. Un fallo
  real de infraestructura se propaga tal cual, igual que ya propagan los
  `Update`/`Create` intermedios del pipeline de PB-026 (mismo criterio, sin
  inventar una excepción nueva).
- La evidencia ensamblada y el `EvidenceCount` se persisten antes de invocar
  al runner, así que sobreviven aunque la investigación termine en `failed`
  — confirmado con test dedicado
  (`TestRunInvestigationKeepsAssembledEvidenceWhenRunnerLaterFails`).

## Riesgos remanentes / gaps documentados

- Igual que en PB-026, `investigations.incident_id` sigue sin FK real
  (pendiente de H2). `evidence.investigation_id` sí tiene FK real porque
  `investigations` ya existe en este repo.
- El `IncidentReader` de producción sigue siendo deny-by-default, así que
  `EvidenceAssembler` es código real pero, igual que `InvestigationRunner`,
  no se ejerce en producción hasta que H2 mergee — sí se ejerce completo en
  tests con fakes y se verificó manualmente insertando datos a mano en
  Postgres.

## Hallazgo no relacionado, confirmado sin cambios (fuera de alcance)

Antes de implementar se confirmó que `TestCreatePersistsAdministrator` y
`TestSavePersistsAdministratorSession` seguían siendo los únicos dos tests
preexistentes rotos (bug de `main` no relacionado, ya señalado en el
`final-summary.md` de PB-026). Se repitió la comprobación al final: siguen
siendo exactamente los mismos dos, nada nuevo se rompió.

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
  bootstrap → TOTP verify → sesión → `listInvestigationEvidence` con
  investigación inexistente (404), con `type_in` inválido (400), sin sesión
  (401), con una investigación insertada a mano sin evidencia (200, lista
  vacía), y con una Evidence `deployment_metadata` insertada a mano (200,
  contenido real serializado correctamente en la respuesta).

No quedan validaciones pendientes ni excepciones al harness para el alcance
de esta tarea.
