# Architecture Review

## Summary

La integración respeta la dirección hexagonal vigente y elimina las
duplicaciones estructurales de la rama auth. Auth e integraciones comparten
infraestructura sin acoplar sus usecases entre sí.

## Layering

- Core/usecases no importan adapters, GORM, `net/http` ni Chi.
- Los tags GORM del dominio son declarativos y no incorporan comportamiento de
  persistencia, conforme ADR-012.
- Queries, columnas técnicas, transacción GORM y cifrado permanecen en adapters.
- `out.Transactor` es el boundary application-level; la propagación del handle
  es privada de PostgreSQL.
- El bootstrap posee una sola conexión, registry de migraciones y Credential
  Store para todos los módulos.

## Modularity

- Seguridad se divide por responsabilidad y el Store duplicado fue eliminado.
- REST conserva feature DTO folders, una estructura por archivo y mappers con
  una responsabilidad exportada por archivo.
- El router `net/http` monta auth, callbacks públicos e integraciones protegidas
  sin introducir un segundo stack HTTP.

## Persistence and migrations

- Entidades auth se persisten directamente; los hashes son proyecciones/columnas
  técnicas y los seeds permanecen en `credentials`.
- Migraciones ordenadas `01..08`, SQL explícito y rollback; no hay
  `AutoMigrate` global ni dependiente del modelo Go actual.
- ADR-014 documenta cuándo usar y cuándo no usar el transactor.

## Findings

- Recovery sigue documentado en el contrato MVP pero permanece fuera de este
  router por alcance explícito; requiere su propia tarea antes de habilitarse.
- El limiter en memoria es apropiado para el despliegue MVP single-instance,
  pero no constituye coordinación distribuida.

## Result

pass
