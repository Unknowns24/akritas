# ADR-013 — Errores propiedad de su capa o adapter

## Estado

Accepted

## Contexto

`domain.Error` es el contrato enriquecido compartido, pero ubicar sentinels REST
o PostgreSQL en el catálogo de dominio mezcla causas de transporte e
infraestructura con reglas del core.

## Decisión

Cada capa declara y cataloga sus propios errores usando `domain.Error`:

- dominio y usecases: `internal/core/domain`;
- REST: `internal/adapter/rest/errors`;
- PostgreSQL: `internal/adapter/db/postgres/errors`;
- adapters externos: dentro del adapter correspondiente.

Un error específico de adapter debe normalizarse antes de cruzar un puerto
cuando el caller sólo necesita semántica de aplicación. Los códigos mantienen
el nibble de capa definido por la policy y continúan registrados en
`docs/errors/aaa-map.md`.

## Consecuencias

- El core no contiene catálogos REST/DB/external.
- Los response mappers siguen aceptando `domain.Error` como contrato común.
- Tests de catálogo se ejecutan en la capa propietaria, evitando dependencias
  invertidas desde dominio hacia adapters.

