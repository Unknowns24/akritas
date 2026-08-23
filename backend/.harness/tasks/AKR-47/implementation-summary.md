# Implementation Summary - AKR-47

## Resultado

AKR-47 fortalece la redaccion de Evidence y GitHub Issue content sin cambiar el contrato HTTP publico.

## Redaccion

- `evidencesafety.Redact` cubre JSON string secrets, asignaciones quoted con espacios, `Authorization: Bearer`, `Authorization: Basic`, cookies, GitHub PAT/App tokens, JWT/session tokens, DSN con credenciales y PEM private keys.
- `RedactAndLimit` mantiene UTF-8 valido y ahora incluye el sufijo de truncado dentro del limite solicitado.
- La redaccion evita modificar texto normal que solo menciona `token`, `password` o `cookie` sin valor asociado.

## IssueContentBuilder

- El builder aplica `safeBounded` sobre titulo y body: redacta, acota, redacta de nuevo el contenido final y vuelve a acotar.
- Los tests colocan secretos en titulo, Project/application/environment/repository/default branch, Evidence summary/content/patch/file/commit, root cause, summary, hypotheses, relevant files, relevant commits y recommended actions.
- El contenido auditable conserva Project, application/environment, repository/default branch, Incident ID/fingerprint/severity/occurrences/first seen/last seen, Evidence persistida, root cause status, causa/hipotesis, confidence, resolution status, summary, relevant files, relevant commits y recommended actions.

## Persistencia

- Nueva migracion append-only `20260823_07_enforce_issue_reference_investigation_incident`.
- PostgreSQL agrega `UNIQUE (id, incident_id)` en `investigations` y FK compuesta desde `github_issue_references(investigation_id, incident_id)` hacia `investigations(id, incident_id)`.
- El integration test de `githubissuereference` prueba que una IssueReference con Incident e Investigation existentes pero cruzados es rechazada por PostgreSQL.

## Documentacion

- Se actualizo memoria del harness y decisiones de publicacion/idempotencia.
- Se agregaron `docs/domain.md` y `docs/erd.md` para documentar el vinculo Incident/Investigation/Evidence/IssueReference.
- No se modifico OpenAPI ni mapa de errores porque no hubo cambio de contrato HTTP ni errores nuevos.
- Se agrego una correccion posterior a summaries/reviews de H4 para no sobreafirmar la cobertura original de redaccion.

