# Implementation Brief — AKR-H4-GITHUB-ISSUE

## Task

Implementar AKR-45..48 sobre el estado combinado H1–H3.

## Current project context

H3 completa Investigation y Operation pero deja Incident en `investigating`.
`GitHubIssueReference` sólo es un value object embebido en un JSONB histórico;
no existe publisher, repository dedicado ni timeline REST. OpenAPI 1.6.0 ya
define Incident detail con latest Investigation/Issue y el endpoint timeline.

## Proposed approach

Persistir una IssueReference normalizada por Investigation, construir Markdown
sanitizado en un servicio aislado, publicar mediante un nuevo capability port
implementado por el cliente GitHub existente y extender el runner H3 para
completar H4 antes de finalizar la Operation.

## Architecture impact

Se mantiene `adapter → usecase/service → core`. El core no conoce HTTP/SDK ni
credenciales. Builder, publisher, repository, orchestration y REST permanecen
separados por responsabilidad.

## API/OpenAPI impact

Sin rutas ni schemas nuevos. Se implementan la proyección Incident detail y
`GET /incidents/{incident_id}/timeline` ya publicados en OpenAPI 1.6.0.

## Data/persistence impact

Migración aditiva `20260823_06_add_github_issue_references`, FK RESTRICT a
Incident/Investigation, PK/unique por Investigation, unique repository+number e
índice por Incident/created_at. El JSONB histórico queda sin nuevos usos.

## Error handling impact

Errores nuevos pertenecen a GitHub externo y PostgreSQL; se registran una vez
en sus catálogos y en el mapa AAA. GitHub/provider bodies y causas internas no
se exponen.

## Test strategy

Tests primero en dominio, builder, adapter GitHub, usecase, repositories,
migración, REST/timeline y escenario H2→H4 con PostgreSQL real + HTTP controlado.

## Risks

Existe una ventana no atómica entre crear la Issue y persistir la referencia.
H4 evita re-publicar cuando ya hay referencia y agrega un marcador durable para
reconciliación futura, pero la reconciliación remota completa pertenece a H6.
