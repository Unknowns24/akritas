# Implementation Summary

## Resultado

Se implementó y corrigió el incremento H1 de integraciones para AKR-5..12 y
AKR-21, sin incorporar persistencia o mutaciones de `Project` ni autenticación
PB-061..063.

## Decisiones y estructura

- ADR-010 formaliza PostgreSQL antes de usarlo como dependencia permanente.
- ADR-011 centraliza configuración en `config/config.go` mediante Viper, con una
  única estructura etiquetada, defaults y validación fail-closed.
- ADR-012 permite tags GORM pasivos en entidades persistibles sin imports ni
  comportamiento GORM en core.
- ADR-013 asigna los errores REST, DB y external a sus capas propietarias.
- El antiguo `internal/app/integrations` era el composition root: construía DB,
  repositorios, adapters, usecases y router. Se movió sin cambiar esa
  responsabilidad a `internal/bootstrap/integrations`, nombre previsto por el
  harness.

## Persistencia y secretos

- PostgreSQL/GORM y las cinco migraciones gormigrate reversibles se mantienen
  sin cambios de schema ni `AutoMigrate` global.
- `GitHubAccount`, `DokployServer`, `GitHubAppRegistration` y
  `GitHubAppBinding` se persisten directamente desde dominio.
- Se eliminó `internal/adapter/db/postgres/models`; sólo
  `integration_credentials` conserva un record privado del adapter DB.
- Credential Store usa AES-256-GCM, nonce aleatorio y AAD por owner, clase y
  versión. PAT, API key, private key y webhook secret nunca aparecen en dominio
  público o responses.

## Paginación Uker

- Se incorporó `github.com/unknowns24/uker` v1.2.2 como única implementación de
  parsing, cursores y aplicación de paginación.
- `internal/core/ports/paging` expone el alias de `pagination.Params` y un
  resultado genérico con items, total y boundaries sin codec propio.
- REST usa `ParseWithSecurity` y `BuildPageSigned`; PostgreSQL usa `Apply` y
  `ApplyFilters`; GitHub/Dokploy sólo traducen boundaries firmados a
  página/offset del proveedor.
- Allowlists, scope firmado por integración, TTL, rechazo de overrides y orden
  total `created_at,id` quedan aplicados.
- OpenAPI avanzó a `1.2.0` y el límite por defecto pasó de 20 a 25. Los cursores
  del codec eliminado no son compatibles; el módulo nunca estuvo montado.

## REST y errores

- Todos los contratos REST tienen sufijo `DTO`, una estructura por archivo y
  están agrupados en `dto/github`, `dto/dokploy` o `dto/common`.
- Las conversiones viven en `internal/adapter/rest/mapper`, una responsabilidad
  por archivo.
- Los errores REST viven en `internal/adapter/rest/errors`, los PostgreSQL en
  `internal/adapter/db/postgres/errors` y el error interno GitHub en su adapter.
- El gate arquitectónico valida tags/imports, ownership de códigos, sufijos DTO
  y SRP de mappers.

## Funcionalidad H1 conservada

- GitHubAccount PAT: CRUD seguro, validación/rotación, connection test y
  discovery.
- GitHub App Manifest: conversión, states de uso único, instalación verificada,
  binding y installation tokens efímeros.
- DokployServer: CRUD seguro, health validation, connection test, discovery y
  protección SSRF/DNS rebinding.
- Todos los endpoints previstos tienen handlers, sin implementar Projects ni
  hitos posteriores.

## Límites deliberados

- `cmd/main.go` no monta integraciones hasta disponer del middleware
  administrativo/Origin de PB-061..063.
- Bootstrap verifica ese middleware antes de abrir PostgreSQL o migrar.
- El reader de referencias Project continúa fallando cerrado hasta PB-010/011;
  por ello los deletes no quedan productivamente habilitables todavía.
