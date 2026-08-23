# Implementation brief — AKR-H1-PROJECTS

## Estado inicial

`feat/project-handling` parte de `main` y duplica configuración, conexión DB,
migraciones, repositorios de integraciones, sesión, router, errores y paginación.
La rama destino ya posee reemplazos autoritativos para esas responsabilidades.

## Estrategia

Se preserva ancestry mediante un merge `ours` y se reconstruye únicamente la
capacidad Project sobre la arquitectura vigente:

```text
REST/Chi + Uker -> ProjectUseCase -> ProjectStore + provider gateways
                                      -> PostgreSQL / GitHub / Dokploy
```

Las llamadas externas terminan antes de persistir. Project guarda snapshots no
secretos y nunca credenciales. El repositorio Project también responde el port
de uso de integraciones para que account/server no puedan eliminarse mientras
sean referenciados.

## Contrato

Se implementan los paths Project ya publicados y se agrega DELETE con respuesta
204. OpenAPI pasa a 1.4.0. Los métodos inseguros requieren sesión y Origin
permitido. La lista usa cursores Uker firmados, filtros allowlisted y orden
determinístico.

## Persistencia

La migración `20260822_09_add_projects` crea la tabla explícitamente, con FKs a
GitHub/Dokploy, JSONB para patrones, checks de estados/límites, nombre único
case-insensitive y exclusividad de la aplicación Dokploy.
