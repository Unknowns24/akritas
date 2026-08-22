# Architecture Review

## Veredicto

**Aprobado sin hallazgos bloqueantes.**

## Decisiones revisadas

- Profile `backend_api`; workflow `backend-api-feature`.
- ADR-010 a ADR-013 están aceptados y alineados con policies, memoria y
  arquitectura.
- El composition root vive en `internal/bootstrap/integrations` y sólo realiza
  construcción de dependencias.

## Dirección de dependencias

- Core/usecases no importan adapters, HTTP, GORM, gormigrate ni drivers.
- Los tags GORM del dominio son metadata pasiva; queries, tablas, transacciones
  y errores DB permanecen en PostgreSQL.
- Los usecases dependen de puertos; REST, DB y providers implementan esos
  límites.
- El record cifrado del Credential Store permanece privado al adapter conforme
  ADR-005.

## Modularity y contratos

- Usecases, repositorios y handlers están separados por feature/operación.
- DTOs REST usan sufijo `DTO`, una estructura por archivo y paquetes por
  feature/common; mappers tienen una conversión pública por archivo.
- Los catálogos de errores están en su capa y cruzan límites mediante
  `domain.Error` normalizado.
- Uker es el único contrato/codec de paginación: alias en ports, parsing/firma
  en REST y aplicación en repositorios.

## Persistencia y compatibilidad

- Se eliminaron modelos GORM duplicados y se verificó la compilación del test de
  compatibilidad contra las cinco migraciones existentes.
- No hubo cambios de tablas/columnas ni nuevas migraciones.
- El test PostgreSQL real usa Testcontainers y no sustituye semántica con
  SQLite.

## Runtime y alcance

- Los 16 endpoints H1 corresponden a OpenAPI; no se agregaron rutas de Projects
  ni otros hitos.
- Los callbacks GitHub son públicos por diseño; el resto sólo puede envolverse
  con middleware administrativo no nulo.
- `Build` falla antes de abrir DB cuando falta PB-061..063 y `cmd/main.go` no lo
  invoca.
- El usage reader de Project conserva bloqueo fail-closed hasta PB-010/011.

## Riesgo residual

Docker no estaba activo en este host, por lo que Testcontainers confirmó el
`SKIP` previsto. La suite de integración compiló y el caso debe ejecutarse en CI
o un host con Docker antes de desplegar.
