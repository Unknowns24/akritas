# TDD Test Plan - AKR-47

## Gate humano

Este plan requiere aprobacion antes de implementar codigo, segun `AGENTS.md` y `testing.md`.

## Tests de redaccion

Agregar tests table-driven en `internal/service/evidencesafety/redact_test.go`.

Cada caso debe verificar:

- el secreto exacto desaparece;
- no queda un fragmento sensible parcial cuando aplique;
- el marcador seguro de redaccion aparece;
- el resultado conserva UTF-8 valido;
- menciones normales como `token rotation needed` o `password policy changed` no se redactan si no tienen valor asociado.

Casos minimos:

- `{"password":"supersecret"}`.
- `{"token": "secret-value"}`.
- JSON con keys case-insensitive, por ejemplo `{"API_KEY":"secret value"}`.
- `PASSWORD="two words"`.
- `TOKEN='secret value'`.
- asignaciones lowercase/mixed-case de `secret`, `api_key`, `cookie`, `session`.
- `Authorization: Bearer <valor>`.
- `Authorization: Basic <valor>`.
- GitHub PAT `github_pat_...`.
- GitHub classic/app-style `ghp_...`, `gho_...`, `ghu_...`, `ghs_...`, `ghr_...`.
- JWT de tres segmentos base64url.
- session token por asignacion/cookie.
- DSN `postgres://user:password@host/db`, `mysql://user:password@host/db` y `redis://:password@host/0`.
- PEM `-----BEGIN PRIVATE KEY----- ... -----END PRIVATE KEY-----`.
- valores con espacios y delimitadores posteriores.

## Tests del IssueContentBuilder

Agregar tests en `internal/service/issuecontent/builder_test.go` que coloquen secretos en:

- titulo del Incident;
- Project/application/environment/repository/default branch;
- Evidence summary/content/patch/file/commit;
- root cause y summary;
- hypotheses;
- relevant files;
- relevant commits;
- recommended actions.

Cada test debe verificar:

- ningun secreto exacto ni fragmento sensible requerido aparece en `Title` o `Body`;
- el titulo y body siguen dentro de `maximumTitleBytes` y `maximumBodyBytes`;
- el marcador `<!-- akritas:investigation_id=... -->` se conserva;
- siguen presentes Project, application, environment, repository, default branch, Incident ID, fingerprint, severity, occurrences, first seen, last seen, Evidence persistida, root cause status, causa/hipotesis, confidence, resolution status, summary, relevant files, relevant commits y recommended actions;
- se mantiene la separacion visible entre `## Observed Evidence` y `## Investigation - QVAC Analysis`.

## Tests PostgreSQL

Agregar/actualizar integration tests del repositorio/migracion para probar:

- se puede persistir una `GitHubIssueReference` cuando `incident_id` coincide con el `incident_id` de la Investigation;
- PostgreSQL rechaza una `GitHubIssueReference` cuyo `incident_id` existe pero no corresponde al `investigation_id`;
- la idempotencia por Investigation sigue devolviendo `ErrGitHubIssueAlreadyPublished` en duplicados.

## Validaciones obligatorias al finalizar

Ejecutar y registrar resultados reales de:

- `go fmt ./...`
- `go build ./...`
- `go test ./...`
- `go test -race ./...`
- `go test -tags=integration ./...`
- `go vet ./...`
- `bash .harness/kernel/scripts/check-backend-architecture.sh`
- `bash .harness/kernel/scripts/check-openapi.sh`
- `bash .harness/kernel/scripts/check-security.sh`
- `git diff --check`

## Documentacion

Actualizar solo cuando corresponda:

- ADR/memoria de publicacion e idempotencia para mencionar la constraint compuesta Incident/Investigation.
- documentacion de dominio/ERD si existe o si se agrega como parte de esta tarea.
- summaries y reviews de AKR-47 con pruebas efectivamente ejecutadas.
- corregir afirmaciones previas de cobertura si siguen siendo inexactas despues de los tests nuevos.

