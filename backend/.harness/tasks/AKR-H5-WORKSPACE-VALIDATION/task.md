# AKR-H5-WORKSPACE-VALIDATION - Remediation workspace branch + validation execution + persistence

## Estado

complete

## Tipo de tarea

backend-service-feature

## Modo de proyecto

existing_project

## Contexto

Hito 5 "Autonomous Remediation" necesita, entre otras cosas, que una futura
`Remediation` trabaje sobre una branch dedicada, ejecute validaciones
(tests/build/static analysis) de forma determinística y segura, y persista
esos resultados para auditoría. En paralelo, otro desarrollador está
implementando H4 (`Incident → Investigation → Issue → GitHub Issue →
IssueReference → proyección API`) en otra branch. Esta tarea entrega
capacidades aisladas y testeadas que una orquestación futura (AKR-49, fuera
de alcance) conectará al trigger `Investigation.resolution_status ==
fixable` una vez que H4 se mergee.

El dominio (`domain.Remediation`, `domain.ValidationResult`,
`domain.CodeChange`, `domain.PullRequestReference`) y su catálogo de errores
(`0x406xxx`) ya existen y están testeados. `docs/openapi.yaml` ya contractea
`Remediation`/`ValidationResult`/`CodeChange`/`PullRequestReference`. No
existe infraestructura de Git ni de ejecución de procesos en el repo
(`internal/adapter/external/git/` es un placeholder vacío).

## Objetivo

Implementar, de forma aislada y testeada:

1. **PB-041** — capacidad para crear una branch dedicada para una futura
   Remediation, nunca operando directamente sobre la branch base/default.
2. **PB-044** — sistema determinístico para ejecutar validaciones (tests,
   build, static analysis) sin permitir que input externo (QVAC, Evidence,
   contenido del repositorio) se convierta en un comando de shell arbitrario.
3. **PB-045** — persistencia de `ValidationResult` detrás de un repository
   port, de forma que `Remediation → ValidationResult[]` sea reconstruible
   sin volver a ejecutar nada.

## Requerimiento funcional

- `RepositoryWorkspace.CreateBranch` (adapter Git local vía CLI, argv fijo,
  sin shell) crea y hace checkout de una branch dedicada a partir de una
  base branch recibida como input, rechazando explícitamente operar sobre la
  branch base/protegida.
- `internal/service/validationpolicy.Policy.Plan` detecta el stack (MVP: Go
  vía `go.mod`) y produce un `ValidationPlan` cerrado; si el stack no es
  soportado, `Supported=false` sin fabricar éxito.
- `ValidationRunner.Run` ejecuta un valor cerrado del enum
  `ValidationCommand` (`go_test`/`go_vet`/`go_build`) vía `exec.CommandContext`
  con argv fijo — nunca un string de comando proveniente de input externo.
- `internal/usecase/remediation.CreateRemediationBranch` y
  `.ExecuteRemediationValidations` orquestan lo anterior detrás de ports,
  sin depender de `IssueReference` ni de la orquestación
  `Investigation → Issue → Remediation`.
- `RemediationStore` (Create+Get) y `ValidationResultStore`
  (Create+ListByRemediation) persisten en PostgreSQL vía nuevas tablas
  `remediations` y `validation_results` (migraciones nuevas, append-only).

## Criterios de aceptación

Ver `.harness/tasks/AKR-H5-WORKSPACE-VALIDATION/tdd-test-plan.md` para el
detalle completo. Resumen:

- Crea branch dedicada; usa la branch base correcta; nunca modifica
  directamente la branch base; nombre de branch determinístico
  (`akritas/remediation/<remediation-id>`); error explícito si falla la
  creación; contexto cancelado aborta la operación; comportamiento
  idempotente por `RemediationID`.
- Repositorio soportado → genera un plan; ejecuta test/build/static
  analysis; exit 0 → passed; exit != 0 → failed; timeout → resultado
  explícito distinguible; error del runner → distinguible de fallo de
  validación; no acepta comandos arbitrarios; cancelación de contexto
  detiene la ejecución.
- Persiste `ValidationResult`; múltiples resultados quedan asociados
  correctamente a una Remediation; preserva tipo/status/timestamps/resumen
  auditable; error de PostgreSQL mapea a error de dominio/aplicación; el
  retrieval devuelve resultados correctos en orden estable.

## Restricciones técnicas

- TDD estricto: tests antes que código productivo.
- Ningún cambio a `docs/openapi.yaml` (contrato ya existente y suficiente).
- Ningún cambio a la lógica de negocio de Incident/Investigation, a
  `IssueReference`, a la creación/contenido de GitHub Issues, ni a sus
  proyecciones HTTP.
- Ningún cambio a `cmd/main.go`, al router global, ni a
  `internal/bootstrap/integrations/module.go` (wiring queda documentado,
  no implementado).
- `internal/core` e `internal/usecase` no pueden importar `os/exec`,
  `gorm.io`, `net/http` ni declarar códigos de error REST/DB/adapter
  (gate `check-backend-architecture.sh`).
- No se implementa AKR-49 (trigger), AKR-51/52 (generación de cambios/tests),
  AKR-55 (decisión ante fallo de validación), AKR-56/57 (commit/PR), ni H6.
- Migraciones nuevas son aisladas y pueden requerir renumeración tras un
  rebase con H4 (mismo directorio de migraciones, fecha `20260823`).

## Archivos o zonas probablemente afectadas

Ver el árbol de archivos completo en el plan aprobado
(`/Users/agustinbressan/.claude/plans/quiero-que-implementes-un-cozy-key.md`).
Resumen: `internal/core/ports/{in,out}` (nuevos archivos), nuevo
`internal/service/validationpolicy`, nuevo `internal/usecase/remediation`,
nuevo `internal/adapter/external/git` (implementación del placeholder
existente), nuevo `internal/adapter/external/validationrunner`, nuevos
`internal/adapter/db/postgres/repository/{remediation,validationresult}`,
dos migraciones nuevas, y ediciones aditivas a
`internal/core/domain/errors.go`, `internal/errorcatalog/catalog_test.go`,
`internal/adapter/db/postgres/errors/catalog.go`,
`internal/adapter/db/postgres/migrations/migrate.go`,
`internal/adapter/db/postgres/dbtest/dbtest.go`, `docs/errors/aaa-map.md`.

## Fuera de alcance

PB-040/AKR-49 (trigger de creación de Remediation), PB-042/043 (AKR-51/52,
generación de cambios y tests), PB-046-050 (AKR-55-59: manejo de fallo,
commit, PR, trazabilidad completa, UI), H4 completo, H6, wiring en
`cmd/main.go`/bootstrap, endpoints REST nuevos.

## Instrucción para el harness

Este task ya cuenta con `implementation-brief.md` y `tdd-test-plan.md`
derivados del plan aprobado por el usuario en modo plan (ver referencia
arriba). La aprobación humana del enfoque ya fue otorgada vía
`ExitPlanMode`. Proceder con implementación TDD, validaciones, revisión de
arquitectura, revisión de seguridad y actualización de memoria del harness.
