# TDD test plan — AKR-H2-DETECTION-INCIDENTS

Plan aprobado explícitamente por el pedido de implementación del 2026-08-22.

## Adquisición y cursor

- parsear logs Dokploy con timestamps RFC3339Nano, orden estable y stream unknown;
- `from_now` por defecto y backfill one-shot `last_10000`;
- checkpoint durable entre ejecuciones y overlap sin reprocesamiento;
- timestamps repetidos usan ordinal/hash y rotación se resuelve fail-closed;
- fetch, cálculo o persistencia fallidos no avanzan el cursor.

## Detection Engine

- reconstruir stack traces, panic y cadenas de excepción como un evento;
- `ignored_patterns` precede a cualquier regla positiva;
- cubrir cada built-in, regex custom case-sensitive y no-match;
- redactar secretos antes de persistir;
- normalizar sólo valores volátiles documentados;
- equivalentes producen el mismo fingerprint y errores materiales uno distinto;
- contexto before/after acotado, durable y finalizado por conteo o 30 segundos.

## Persistencia y concurrencia

- migraciones y rollback explícitos;
- occurrence key única e idempotencia de retry;
- primer evento crea Incident;
- evento dentro de la ventana inclusiva actualiza usando `last_seen_at`;
- evento fuera de ventana crea otro Incident;
- contador, timestamps y promoción de severidad correctos;
- cursor, LogEvent e Incident se confirman o revierten juntos;
- agrupación concurrente queda serializada por locks transaccionales;
- FKs restringen borrar Projects con historia.

## REST y ejecución background

- enum/default de `initial_log_ingestion` y rechazo de valores inválidos;
- lista, detalle y LogEvents con filtros, Uker, 401/404/500;
- serialización omite opcionales futuros y no expone secretos;
- sólo Projects habilitados adquieren logs;
- pendientes pueden finalizar sin nueva adquisición y shutdown respeta contexto.

## Validación final

- `go fmt ./...`
- `go test ./...`
- `go test -race ./...`
- `go test -tags=integration ./...`
- `go vet ./...`
- gates de arquitectura, OpenAPI y seguridad;
- `git diff --check` e inspección de secretos/artefactos.
