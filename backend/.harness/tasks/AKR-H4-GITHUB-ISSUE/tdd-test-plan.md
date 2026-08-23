# TDD Test Plan — AKR-H4-GITHUB-ISSUE

## Scope

AKR-45, AKR-46, AKR-47 y soporte backend AKR-48, sin H5.

## Tests to add/update

- Dominio: IssueReference vinculada, lifecycle fixable/requires_human e
  invariantes de duplicación.
- Builder: Project/application/repository, Incident/timestamps/occurrences,
  clasificaciones/confidence, Evidence separada, files/commits/actions,
  determinismo, límites y secret redaction.
- GitHub: POST al repository correcto, PAT/App mediante Credential Store,
  mapping de response y fallos 401/403/404/422/5xx/network seguros.
- Orchestration: ambos resolution statuses publican; fallos previos no
  publican; GitHub/persistence fallan explícitamente; Operation sólo sucede
  después de referencia durable; referencia existente evita duplicación.
- PostgreSQL: create/find/latest, constraints/FKs/índices, clean/H3 migration,
  rollback/reapply y flujo H2→H3→H4.
- REST: Incident detail con latest Investigation/Issue, timeline `issue_created`,
  Evidence/Investigation conservadas y JSON sin secretos.
- Regresiones H1–H3 y catálogos globales.

## Expected failing tests before implementation

No existen los nuevos tipos/ports/builder/repository/publisher/timeline y el
runner finaliza Operation antes de H4, por lo que los tests nuevos inicialmente
no compilarán o fallarán.

## Acceptance criteria covered

AKR-45, AKR-46, AKR-47, AKR-48 / PB-036..PB-039.

## Human approval

Aprobado explícitamente por el usuario el 2026-08-23 junto con estas decisiones:
una Issue por Investigation, Issue más reciente en Incident, `fixable` espera H5
en `publishing_issue` y timeline incluido en H4.
