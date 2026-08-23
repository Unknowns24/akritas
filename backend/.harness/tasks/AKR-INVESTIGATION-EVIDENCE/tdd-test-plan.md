# TDD test plan — AKR-INVESTIGATION-EVIDENCE

Plan sujeto a aprobación explícita antes de escribir tests o implementación.

## Dominio

- Sin tests nuevos de comportamiento (`Evidence` no cambia); solo confirmar
  que las tags gorm no rompen los tests existentes.

## Ports / fakes

- fake `IncidentReader` extendido con `Get` controlable (incidente
  existente/`ErrIncidentNotFound`/error de infra).
- fake `EvidenceStore` (`Create`, `ListByInvestigation`) y fake
  `EvidenceAssembler` (resultado controlable) para `RunUseCase`.

## Servicio — evidenceassembly

- incident/project encontrados -> una Evidence `deployment_metadata` con
  contenido derivado de campos reales de `Project` (nombre, monitoring/
  health status, repo owner/name/branch, aplicación Dokploy/entorno/estado);
  nunca credenciales.
- `IncidentReader.Get` devuelve `ErrIncidentNotFound` -> slice vacío, sin
  error.
- `ProjectStore.Get` devuelve `ErrProjectNotFound` -> slice vacío, sin error.
- error de infraestructura en cualquiera de los dos pasos -> se propaga tal
  cual (no se traga).

## Usecase — RunUseCase (extensión)

- happy path: `Assemble` devuelve 1+ Evidence -> cada una se persiste via
  `EvidenceStore.Create`, `Investigation.EvidenceCount` queda persistido con
  la cantidad real antes de invocar al runner.
- el runner falla después de ensamblar evidencia -> la Investigation queda
  `failed` pero la Evidence ya persistida y el `EvidenceCount` no se
  revierten (se conservan).
- `Assemble` devuelve slice vacío sin error -> `EvidenceCount = 0`, el
  pipeline continúa normalmente hacia el runner.
- `Assemble` devuelve error de infraestructura -> `Execute` lo propaga sin
  invocar al runner ni a `failInvestigation` (mismo criterio que un fallo de
  `Update`/`Create` ya existente en el pipeline).

## Usecase — evidence

- `ListInvestigationEvidence`: investigación inexistente -> propaga
  `ErrInvestigationNotFound`; happy path delega en
  `EvidenceStore.ListByInvestigation` y devuelve la página tal cual.

## PostgreSQL

- migración y rollback de `evidence`, FK real a `investigations` (ON DELETE
  CASCADE), checks de `type`/`redacted`/coherencia de líneas;
- `Create` persiste todos los campos opcionales (`file_path`, `line_start`/
  `line_end`, `commit_sha`, `patch`, `occurred_at`) correctamente nulos o
  seteados;
- `ListByInvestigation` scoped por `investigation_id`, paginación real y
  filtro `type_in`.

## REST

- inventario exacto de la ruta nueva y su envelope;
- `type_in` con un valor fuera del enum -> 400;
- 401 sin sesión, 404 con `investigation_id` inexistente;
- `EvidenceDTO` nunca expone `redacted: false` (el dominio ya lo garantiza,
  se confirma que el mapper no lo pisa) ni contenido con forma de credencial.

## Validación final

- `go test ./...`
- `go vet ./...`
- `gofmt -l .` (sin diffs)
- `.harness/kernel/scripts/check-backend-architecture.sh`
- `.harness/kernel/scripts/check-openapi.sh`
- `.harness/kernel/scripts/check-security.sh`
- Confirmar antes de correr nada que `TestCreatePersistsAdministrator` y
  `TestSavePersistsAdministratorSession` siguen siendo los únicos dos tests
  preexistentes rotos (no relacionados a esta tarea) y que nada nuevo se
  rompió.
- prueba manual end-to-end contra Postgres local: `listInvestigationEvidence`
  con `investigation_id` inexistente (404) y con una investigación creada a
  mano en la base (200, lista vacía dado que el pipeline async no corre sin
  un `IncidentReader` real en producción).
