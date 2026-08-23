# ADR-012 — Tags GORM en entidades de dominio persistibles

## Estado

Accepted

## Contexto

ADR-010 confinó GORM al adapter PostgreSQL y la primera implementación creó
estructuras de persistencia que duplicaban campo por campo a `GitHubAccount`,
`DokployServer` y otros conceptos persistibles. Esa duplicación aumentó el
costo de mantenimiento sin representar un schema histórico diferente.

## Decisión

Las entidades de dominio que representan directamente una tabla vigente pueden
declarar tags estructurales `gorm` y ser leídas/escritas por el adapter DB.

- `internal/core` no importa GORM, gormigrate, drivers ni tipos de base de datos.
- Los tags sólo describen columnas, claves e índices; no introducen queries,
  hooks ni comportamiento de persistencia en dominio.
- El nombre de tabla y las consultas siguen perteneciendo al adapter.
- Los structs auxiliares de migraciones históricas permanecen en migraciones.
- Los registros puramente infraestructurales que no pertenecen al dominio,
  especialmente ciphertext/nonce del Credential Store, permanecen privados en
  su adapter conforme ADR-005.

Esta decisión reemplaza únicamente la exigencia de ADR-010 de mantener todo el
metadata GORM fuera de dominio. GORM como dependencia ejecutable continúa
confinado al adapter PostgreSQL.

## Consecuencias

- Se elimina la carpeta de modelos duplicados del adapter PostgreSQL.
- Las entidades y el schema deben evolucionar coordinadamente mediante
  migraciones versionadas; los tags nunca habilitan `AutoMigrate` global.
- El gate arquitectónico permite tags GORM pero continúa rechazando imports de
  infraestructura desde core/usecases.

