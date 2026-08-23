# Architecture Review - AKR-H4-GITHUB-ISSUE

## Veredicto

PASS con una nota de alcance: la reconciliacion remota completa para el caso "GitHub acepto la Issue pero PostgreSQL fallo" queda documentada para H6.

## Dependencias

- El core declara entidades y ports de aplicacion; no importa GORM, HTTP ni SDKs.
- `service/issuecontent` depende de tipos application-level y no conoce GitHub ni PostgreSQL.
- El adapter GitHub implementa el port outbound y reutiliza autenticacion/Credential Store existentes.
- Los repositorios PostgreSQL son SRP y se mantienen bajo `internal/adapter/db/postgres/repository/*`.
- REST consume usecases y mappers; no filtra DTOs hacia el core.

## Transacciones

La llamada externa a GitHub queda fuera de transacciones PostgreSQL. Las transacciones son breves:

1. completar Investigation y mover Incident a `publishing_issue`;
2. persistir IssueReference, actualizar Incident/Operation.

Esto respeta ADR-014 y evita locks durante I/O remoto.

## Idempotencia

La referencia durable por Investigation corta publicaciones repetidas. La unicidad por Incident no se impone porque un mismo Incident puede ser investigado mas de una vez; el detalle REST resuelve la Issue mas reciente por `(incident_id, created_at)`.

## Scope

No se implementaron remediation, branches, writes de codigo, commits, PRs ni deploy.
