# TDD test plan — AKR-H1-PROJECTS

Plan aprobado por el pedido de implementación del 2026-08-22.

## Dominio

- validar límites, regex, tiempos, estados y coherencia enabled/status;
- activar o reconfigurar lleva a `starting`; desactivar conserva configuración;
- readiness exige snapshots válidos y configuración habilitada;
- asociaciones no pueden cambiar con monitoring activo.

## Gateways

- resolver repositorio por ID opaco y por `owner/name` con metadata real;
- rechazar owner o default branch diferentes;
- resolver aplicación por identificador exacto;
- distinguir not found, credencial rechazada e integración no disponible.

## Use case

- create/get/list/update/delete y get/put monitoring;
- nombre duplicado y aplicación exclusiva producen conflicto;
- asociaciones/borrado activos producen conflicto;
- activación revalida ambos proveedores y no persiste ante fallo;
- updates concurrentes no pisan cambios previos.

## PostgreSQL

- migración, rollback, FKs, índices y checks;
- persistencia completa de snapshots/configuración;
- filtros/sorts Uker y conteos de uso de integraciones;
- errores not-found, duplicate y concurrent-update normalizados.

## REST

- inventario exacto de rutas, envelopes y códigos 201/200/204;
- UUID/body/query inválidos, 401, 403, 404 y 409;
- cursores firmados y filtros estables;
- respuestas Project sin tokens, passwords, API keys ni secretos.

## Validación final

- `go test ./...`
- `go fmt ./...`
- `.harness/kernel/scripts/check-backend-architecture.sh`
- `.harness/kernel/scripts/check-openapi.sh`
- `.harness/kernel/scripts/check-security.sh`
