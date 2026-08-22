# Security Review

## Summary

Cierra el ciclo de autenticación con las mismas garantías ya establecidas: comparaciones en tiempo constante, secretos descifrados sólo en el punto de uso, respuestas genéricas que no permiten inferir la causa de un fallo, y ahora protección real contra reutilización de código TOTP (pendiente desde PB-062 a propósito).

## Auth / permissions

- `POST /auth/login` es público (`security: []`) — correcto, todavía no hay sesión en ese punto.
- `GET`/`DELETE /auth/session` exigen sesión válida vía el nuevo middleware `RequireSession`, que resuelve la cookie a través de `AuthenticateSessionUseCase` — primera vez que el proyecto aplica `cookieAuth` de verdad.
- `PasswordHasher.Verify` compara en tiempo constante (`crypto/subtle.ConstantTimeCompare`) contra la clave re-derivada, no contra el hash completo como string.
- `totpVerifier.Verify` compara cada candidato en tiempo constante (heredado de PB-063, sin cambios funcionales en esta tarea salvo exponer qué período matcheó).
- Rate limiting independiente por IP y por cuenta para login, con **instancia propia** de `RateLimiter` en `cmd/main.go` — separada de las de `/auth/setup` y `/auth/setup/verify` (lección aplicada desde el inicio esta vez, no como corrección posterior).
- Reutilización de un período TOTP ya aceptado por ese `Administrator`: rechazada explícitamente comparando contra `last_accepted_totp_period`, verificado manualmente contra Postgres real (mismo código aceptado dos veces seguidas → la segunda vez `401` genérico).
- Login **nunca** revoca sesiones anteriores del mismo `Administrator` — verificado tanto por test (`TestLoginAdministratorDoesNotRevokeOtherSessions`) como manualmente (3 sesiones coexistiendo en la base tras verify + 2 logins, ninguna revocada por las siguientes).

## Input validation

- `LoginRequest` valida `email`/`password`/`totp_code` en el handler (mismo criterio de shape que endpoints anteriores) antes de invocar el usecase.
- El usecase revalida en profundidad: email inexistente, password incorrecta, TOTP inválido o reutilizado — los cuatro colapsan al mismo `401 ErrInvalidCredentials`, sin distinguir la causa en la respuesta (ADR-008).

## Data exposure

- El token de sesión en claro sólo viaja en `Set-Cookie`, nunca en el body JSON — verificado por test y manualmente (comparando el valor de la cookie contra `token_hash` en la base: no coinciden, como se espera de un hash).
- `password_hash` y `encrypted_totp_secret` nunca salen de `AdministratorCredentials` (interno a la capa de infraestructura/usecase); ningún DTO REST los expone.
- `last_accepted_totp_period` tampoco se expone en ningún DTO ni forma parte de `domain.Administrator` — es puramente un dato de persistencia interno, consistente con el mismo criterio ya aplicado a `password_hash`/`encrypted_totp_secret` desde PB-061/062.
- El logger de GORM sigue en modo `Silent` — sin SQL con valores en logs, verificado manualmente inspeccionando el log del proceso tras el flujo completo (login, replay rechazado, session, logout) sin encontrar password/TOTP/token en claro.

## Error leakage

- `writeLoginError`/`writeSessionError` siguen el mismo patrón ya establecido: sentinels de dominio mapean por sufijo de código, rate limit mapea a `429` con `Retry-After`, cualquier otro error cae al `500` genérico sin exponer la causa real.
- El middleware nunca ecoa el valor de la cookie recibida en su respuesta de error — verificado por test (`TestRequireSessionWithInvalidCookie`).

## Findings

Ningún hallazgo bloqueante nuevo. Un gap ya conocido y explícitamente fuera de alcance se mantiene: `DELETE /auth/session` (mutación autenticada) no valida `Origin`, porque `AKRITAS_PUBLIC_URL`/`AKRITAS_ALLOWED_ORIGINS` todavía no existen en `config.go`. Documentado en detalle en `final-summary.md` (Remaining risks) como corresponde a PB-065 (session hardening), no como omisión silenciosa.

## Result

pass
