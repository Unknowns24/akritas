# Architecture Review

## Summary

La implementación conserva la arquitectura hexagonal, generaliza el concepto
de fuente Dokploy en el dominio y mantiene los detalles del proveedor y Docker
dentro del adapter externo.

## Layering

- Dominio contiene discriminante, invariantes, snapshots e identidad estable
  sin depender de HTTP, GORM runtime ni payloads Dokploy.
- Puertos de entrada exponen discovery de Composes/servicios y comandos Project;
  puertos de salida abstraen resolución remota, persistencia y logs.
- Casos de uso coordinan servidor, repositorios y gateway; handlers sólo
  validan/mapean HTTP.
- Endpoints y labels específicos de Dokploy/Docker permanecen en
  `internal/adapter/external/dokploy`.

## Modularity / SRP

- DTOs y mappers nuevos están separados por recurso para cumplir la regla de
  responsabilidad del profile.
- Discovery, resolución de fuente, persistencia y adquisición de logs mantienen
  fronteras independientes.
- La identidad lógica no incluye container ID; esto evita acoplar checkpoints a
  réplicas efímeras o redeploys.

## OpenAPI consistency

- OpenAPI `2.0.0` conserva `/api/v1` y documenta ambas ramas mediante
  `oneOf`/discriminator.
- Los handlers registrados coinciden con las rutas Compose declaradas.
- El contrato Project ya no expone `dokploy_application`,
  `application_identifier` ni un selector de servidor superior.
- `DokployServer` expone ambos conteos y `refresh` está tipado/validado como
  booleano estricto.

## Findings

No se encontraron hallazgos bloqueantes. El rollback de la migración rechaza de
forma deliberada bases que todavía contengan Projects `compose_service`; antes
de revertir se deben migrar o eliminar esas asociaciones mediante una operación
de producto explícita.

## Result

pass
