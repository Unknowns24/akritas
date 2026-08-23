# AKR-DOKPLOY-COMPOSE-SOURCES

Extender la integración Dokploy para descubrir aplicaciones y Composes, listar
los servicios declarados por cada Compose y permitir que un Project monitoree
una aplicación o un servicio Compose concreto.

## Alcance acordado

- mantener el listado paginado de aplicaciones existente;
- agregar listado paginado de Composes y listado de servicios por Compose;
- permitir refresco explícito de servicios mediante `refresh=true`;
- reemplazar en Project `dokploy_application` por `dokploy_source`;
- soportar fuentes `application` y `compose_service`;
- leer logs de una réplica activa del servicio Compose seleccionado;
- migrar Projects, checkpoints y eventos existentes como fuentes application;
- actualizar OpenAPI, arquitectura, persistencia, monitoreo y evidencia.

## Decisiones de producto

- el cambio del contrato Project es inmediato y no conserva el request/response
  anterior;
- el contrato público continúa bajo `/api/v1`, con OpenAPI `2.0.0`;
- un selector Compose identifica exactamente `servidor + compose + servicio`;
- la misma fuente no puede estar asociada a más de un Project;
- servicios distintos del mismo Compose sí pueden pertenecer a Projects
  distintos;
- se soportan runtimes `docker-compose` y `stack`;
- se monitorea una única réplica activa, elegida determinísticamente;
- no se persisten IDs de contenedor porque cambian durante redeploys;
- la evidencia histórica persistida no se reescribe.

## Estado

Completada el 2026-08-23. El usuario aprobó explícitamente el plan TDD antes de
la implementación. Tests, `go vet`, validación de arquitectura, OpenAPI y
seguridad finalizaron correctamente.
