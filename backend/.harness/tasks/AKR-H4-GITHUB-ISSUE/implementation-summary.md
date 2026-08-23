# Implementation Summary - AKR-H4-GITHUB-ISSUE

## Resultado

H4 implementa el flujo `Investigation completed -> GitHub Issue -> IssueReference -> REST` sin incluir remediation, branches, commits, Pull Requests ni deploy.

## Dominio y puertos

- `GitHubIssueReference` queda como entidad persistible con `incident_id`, `investigation_id`, `issue_number`, `issue_url`, `repository` y `created_at`.
- `Incident` conserva la proyeccion publica singular de Issue y agrega latest Investigation para el detalle REST.
- Se agregaron ports cohesivos para `IssuePublisher`, `IssueContentBuilder`, `GitHubIssueReferenceStore`, timeline e Investigation latest-by-Incident.
- Se registro `ErrGitHubIssueAlreadyPublished` como conflicto idempotente por Investigation.

## Persistencia

- Nueva migracion append-only `20260823_06_add_github_issue_references`.
- La tabla usa FK `RESTRICT` a `incidents` e `investigations`, una referencia maxima por Investigation, unique `(repository, issue_number)` e indice `(incident_id, created_at DESC)`.
- Se agrego un repositorio PostgreSQL SRP con `Create`, `FindByInvestigation` y `FindLatestByIncident`.
- `IncidentWorkflowStore.Update` persiste phase, terminal outcome, summary y clasificaciones QVAC; IssueReference queda en su repositorio dedicado.

## GitHub y contenido

- `IssueContentBuilder` queda aislado de GitHub/PostgreSQL y produce titulo/body deterministas, secciones estables, Evidence separada de QVAC, limites y segunda redaccion defensiva.
- El cliente GitHub existente implementa `PublishIssue` con `POST /repos/{owner}/{repo}/issues`, reutilizando `accountToken`, Credential Store, headers, timeout y normalizacion de errores existentes.
- El adapter devuelve solo numero, URL y timestamp.

## Orquestacion y REST

- `RunInvestigationUseCase` mantiene la Operation running hasta confirmar GitHub y persistir IssueReference.
- El flujo usa una transaccion corta para completar Investigation y mover Incident a `publishing_issue`, cierra la transaccion, publica en GitHub y abre otra transaccion para IssueReference, Incident y Operation.
- `requires_human` termina en `completed/requires_human`; `fixable` queda en `publishing_issue` esperando H5.
- Fallos GitHub no crean referencia y fallan Incident/Operation con mensaje publico seguro.
- `GET /incidents/{incident_id}` proyecta Incident, latest Investigation y latest IssueReference.
- `GET /incidents/{incident_id}/timeline` devuelve eventos derivados de registros persistidos con paginacion Uker por `occurred_at,id`.

## Tests agregados/actualizados

- Orquestacion H4 para fixable/requires_human, fallo GitHub, fallo de persistencia y referencia existente.
- Builder de contenido, redaccion y determinismo.
- Cliente GitHub `PublishIssue` contra `httptest`.
- Repositorio IssueReference con PostgreSQL/Testcontainers.
- Routing/handler/mapper REST y catlogos de errores.
