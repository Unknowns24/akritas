# TDD Test Plan

## Scope

Recovery start/verify, rate limiting acotado/configurable, errores genéricos,
cookie consistente y lifecycle/revocación de sesiones.

## Tests to add/update

1. Config: defaults/overrides y rechazo de attempts, window o max keys fuera de
   rango.
2. Limiter: threshold, aislamiento, reset exacto, cleanup, hard cap y fail
   closed sin crecimiento.
3. Recovery start: éxito, misma password, email/token inválidos equivalentes,
   reemplazo atómico y credenciales activas sin cambios antes de confirmar.
4. Recovery verify: enrollment inválido/expirado/TOTP incorrecto equivalentes;
   éxito rota hash/seed, consume período, revoca todas las sesiones anteriores y
   crea una nueva.
5. Rollback: fallos en rotate, secret move, revoke-all, save o consume no dejan
   estado parcial.
6. Login/recovery: un login con el hash observado antes de recovery no crea
   sesión después del commit.
7. Session repository/usecase: refresh activo condicional; expiradas, revocadas
   y random rechazadas; revoke idempotente; revoke-all.
8. REST/router: DTO/validación, rutas, no-store, cookie completa, status/error
   genérico y RateLimitKey basado únicamente en RemoteAddr.
9. PostgreSQL: commit/rollback compartido y carreras determinísticas contra el
   motor real.
10. Regresión: endpoints autenticados normales no usan limiters de auth y todo
    H1/H2 existente continúa verde.

## Expected failing tests before implementation

No existen ports/usecases/handlers de recovery, el limiter no posee hard cap,
la configuración usa constantes en `cmd/main.go`, el refresh de sesión tiene una
ventana lookup/update y no existen rotate/revoke-all/consume atómicos.

## Acceptance criteria covered

Todos los criterios de `task.md`, el plan aprobado por el usuario y los casos
mínimos de PB-064/PB-065.

## Human approval

Aprobado explícitamente por el usuario el 2026-08-23 mediante “PLEASE IMPLEMENT
THIS PLAN”.
