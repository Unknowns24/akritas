# Architecture Review

## Summary

La implementación respeta el profile `backend_api`, el `tdd-test-plan.md` aprobado (incluyendo el ajuste post-aprobación de envolver `Create`+`Save` en una transacción) y la dirección de dependencias hexagonal ya establecida en PB-061. Cierra el flujo de bootstrap: es la primera vez que el proyecto hace un `INSERT` real en `administrators` y la primera vez que persiste una sesión.

## Layering

- `internal/core/domain` y `internal/core/ports/{in,out}` siguen sin importar `internal/adapter`, GORM, Chi, `net/http` ni `os/exec` (`check-backend-architecture.sh` pasa).
- `internal/usecase/auth` sólo depende de `internal/core/domain` y `internal/core/ports/{in,out}` — incluido el nuevo `Transactor`, que es un port de dominio (`func(ctx) error` genérico), no una abstracción de GORM.
- El nuevo paquete `txcontext` (propagación de `*gorm.DB` transaccional vía `context.Context`) vive enteramente en `internal/adapter/db/postgres` — el usecase nunca ve un `*gorm.DB`, sólo el `context.Context` que ya recibía.

## Modularity / SRP

- Un archivo por operación pública se mantiene: `create.go`, `find_by_id.go`, `delete.go`, `save.go` en cada paquete de repositorio; `totp_verifier.go`, `session_token_generator.go` como adapters de una sola responsabilidad cada uno.
- `Transactor` y `txcontext` son dos piezas separadas con responsabilidades distintas (orquestar la transacción vs. propagarla) en vez de una sola abstracción mixta.
- Nota de diseño, no un defecto: de los tres repositorios de este dominio, `AdministratorRepository` y `AdministratorSessionRepository` son "transaction-aware" (usan `txcontext.From`), mientras que `PendingEnrollmentRepository` no participa nunca de una transacción (su `Delete` corre siempre fuera, por decisión ya aprobada). Es intencional y está documentado en `implementation-brief.md`, pero un lector nuevo del código debería confirmar esa asimetría en vez de asumir que los tres repos se comportan igual.

## OpenAPI consistency

- No se modificó `docs/openapi.yaml`; `check-openapi.sh` confirma 59 operaciones y 112 schemas sin cambios.
- `SessionResponse`/`Session`/`Administrator` espejan el schema exactamente (`Session` no incluye `administrator_id` ni `revoked_at`, que son campos internos de `domain.AdministratorSession` — se dejan fuera del DTO deliberadamente).
- Verificado manualmente: el body y los headers (`Set-Cookie`, `Cache-Control`) coinciden con el contrato para 200; 400/409/429 se comportan igual que en PB-061.

## Findings

- (Heredado de PB-061, se repite aquí porque el mismo patrón se usó de nuevo) `internal/adapter/rest/handler/auth` importa el paquete concreto `internal/usecase/auth` (no sólo `ports/in`) para reconocer `authusecase.ErrSetupRateLimited` en el handler de verify también. Mismo no-bloqueante que en PB-061; permitido por `project-structure.md`.
- Bug real encontrado y corregido durante la verificación manual (no un hallazgo pendiente): la migración `20260822_03` no era idempotente frente a una base nueva, porque `20260822_01`'s `AutoMigrate` ya refleja el struct Go actual de `model.Administrator` (incluyendo el campo agregado por esta tarea). Corregido con un chequeo `HasColumn` y cubierto por `TestMigration03IsIdempotentWhenColumnAlreadyExists`. Se documenta acá porque es una lección de diseño reutilizable: cualquier migración futura que reutilice `AutoMigrate` sobre un struct que otra migración posterior vaya a extender debe considerar esta misma trampa.

## Result

pass
