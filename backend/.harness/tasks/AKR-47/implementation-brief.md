# Implementation Brief - AKR-47

## Profile y workflow

- Profile activo: `backend_service`.
- Workflow: `.harness/kernel/workflows/backend-service-feature.yaml`.
- Policies cargadas: estructura backend, wiring, modularidad/SRP, errores de dominio, migraciones, adapters externos, background processes, Uker, testing, security y architecture decisions.

## Contexto tecnico

- `evidencesafety.Redact` aplica una lista de expresiones regulares y `RedactAndLimit` redacciona antes de truncar, preservando UTF-8.
- `issuecontent.Builder` ya aplica redaccion por campo y una redaccion defensiva final del body, pero los patrones actuales no cubren suficientemente JSON strings, valores entre comillas con espacios, `Authorization: Basic`, cookies, varios tokens de GitHub App, DSN completas ni todos los campos del builder.
- `github_issue_references` tiene FKs independientes a `incidents(id)` e `investigations(id)`, pero PostgreSQL aun permite que una referencia combine un Incident valido con una Investigation valida de otro Incident.
- OpenAPI no deberia cambiar porque no hay cambio de payload, endpoint ni error publico.

## Cambios propuestos

- Ampliar `internal/service/evidencesafety` con redaccion case-insensitive para:
  - secretos JSON con valores string;
  - asignaciones con comillas simples/dobles y valores con espacios;
  - `Authorization: Bearer` y `Authorization: Basic`;
  - GitHub PAT y tokens de GitHub App;
  - JWT y session tokens;
  - passwords, secrets, API keys, cookies y tokens;
  - DSN con usuario y password;
  - private keys PEM.
- Evitar redaccion de texto normal que solo mencione palabras como `token` o `password` sin valor asociado.
- Mantener una unica forma segura de reemplazo sin preservar prefijos o fragmentos del valor secreto.
- Aplicar redaccion defensiva final sobre titulo y body completos en el builder.
- Agregar migracion append-only para constraint compuesta `(investigation_id, incident_id)` contra `investigations(id, incident_id)`.
- Actualizar documentacion de dominio/ADR/memoria/summaries solo con afirmaciones respaldadas por tests y validaciones ejecutadas.

## Riesgos

- Las regex deben redactar de forma defensiva sin destruir contenido auditable normal.
- El truncado debe seguir siendo deterministico y UTF-8 valido despues de la redaccion.
- La migracion debe ser append-only y no reescribir migraciones historicas ya aplicadas.

