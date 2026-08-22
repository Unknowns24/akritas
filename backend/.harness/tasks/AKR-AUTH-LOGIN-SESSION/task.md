# AKR-AUTH-LOGIN-SESSION - Login, sesión opaca y logout

## Estado

complete

## Tipo de tarea

backend-api-feature

## Modo de proyecto

existing_project

## Contexto

`AKR-AUTH-BOOTSTRAP` (PB-061) y `AKR-AUTH-TOTP-VERIFY` (PB-062) dejaron al único `Administrator` activable con su secreto TOTP cifrado persistido. Falta el resto del ciclo de autenticación: iniciar sesión con el segundo factor, leer la sesión actual, y cerrarla. Esta es además la primera tarea que necesita middleware de autenticación real (hasta ahora `internal/adapter/rest/middleware/` está vacío).

## Objetivo

Implementar `POST /auth/login`, `GET /auth/session` y `DELETE /auth/session`: login con email+password+TOTP (respuesta genérica sin distinguir causa de fallo, 401), lectura de la sesión actual vía cookie, y logout que revoca la sesión y expira la cookie.

## Requerimiento funcional

- `POST /auth/login` valida `LoginRequest { email, password, totp_code }`. Falla (email inexistente, password incorrecta, TOTP incorrecto o reutilizado) → siempre la misma respuesta `401` genérica. Éxito → `200 SessionResponse` + `Set-Cookie` (nueva sesión; login NO revoca sesiones anteriores) + `Cache-Control: no-store`.
- Rate limiting independiente por IP y por cuenta (dos llamadas a `RateLimiter.Allow` con keys distintas) → `429` con `Retry-After` si se excede cualquiera de los dos.
- El código TOTP se verifica con tolerancia RFC 6238 de un período, y se **rechaza la reutilización de un período ya aceptado** por ese `Administrator` (a diferencia de PB-062, acá sí aplica: el mismo secreto se verifica repetidamente a través del tiempo).
- `GET /auth/session` requiere sesión activa vía cookie; devuelve `200 SessionResponse` con la proyección segura de la sesión actual. Cada request autenticada exitosa desliza (renueva) el idle TTL — decisión confirmada explícitamente con el usuario.
- `DELETE /auth/session` requiere sesión activa; revoca la sesión, expira la cookie, responde `204`.
- Sesión inexistente/expirada/revocada en `GET`/`DELETE /auth/session` → `401` genérico (mismo criterio que login).

## Criterios de aceptación

- `go test ./...`, `go vet ./...` y `gofmt -l .` (sin diffs) finalizan correctamente.
- Los tres endpoints devuelven exactamente los schemas y códigos de estado de `docs/openapi.yaml`.
- El idle TTL se renueva en cada request autenticada exitosa, sin superar nunca el TTL absoluto.
- Se rechaza la reutilización del último período TOTP aceptado por login.
- El secreto TOTP y el password hash nunca salen del `Administrator` recuperado por email; ninguna respuesta los expone.
- `internal/core/**` no importa GORM, Chi, `net/http` ni SDKs externos.
- Los errores nuevos cumplen `0x4XXNNNT`; se reutilizan sentinels existentes donde aplica.
- Los checks de arquitectura, OpenAPI y seguridad del harness pasan sin modificar `docs/openapi.yaml`.
- Las migraciones `20260822_01` a `04` no se modifican.

## Restricciones técnicas

- Profile: `.harness/kernel/profiles/go-hexagonal-api.yaml`.
- Workflow: `.harness/kernel/workflows/backend-api-feature.yaml`.
- `TOTPVerifier.Verify` se extiende para devolver el período RFC 6238 que efectivamente matcheó (necesario para rechazar reutilización); esto toca el único call site existente en `internal/usecase/auth/verify_administrator_setup.go` (PB-062), que simplemente ignora el nuevo valor de retorno — no se cambia su comportamiento.
- Reutilizar `domain.AdministratorSession.Validate()/IsActive()/Revoke()` tal cual; se agrega un método nuevo (`ExtendIdle`) sin tocar los existentes.
- Nueva migración `20260822_05` (idempotente vía `HasColumn`, aprendiendo de la corrección aplicada en PB-062) para `administrators.last_accepted_totp_period`. No tocar `01`-`04`.
- Rate limiter de login: instancia propia en `cmd/main.go`, separada de las de `/auth/setup` y `/auth/setup/verify` (mismo criterio que la corrección aplicada en PB-062 — nunca compartir una instancia de `RateLimiter` entre endpoints con presupuestos que deben ser independientes).
- No implementar código antes de la aprobación humana de `tdd-test-plan.md`.

## Archivos o zonas probablemente afectadas

- `internal/core/domain/errors.go` (`ErrInvalidCredentials`, `0x401009U`), `internal/core/domain/administrator_session.go` (+`ExtendIdle`).
- `internal/core/ports/out/`: `administrator_repository.go` (+`FindByID`, +`FindByEmail`, +`UpdateLastAcceptedTOTPPeriod`), `administrator_session_repository.go` (+`FindByTokenHash`, +`UpdateIdleExpiry`, +`Revoke`), `password_hasher.go` (+`Verify`), `totp_verifier.go` (firma extendida), `session_token_generator.go` (+`Hash`).
- `internal/core/ports/in/`: `login_administrator.go`, `authenticate_session.go`, `get_current_session.go`, `logout_administrator.go`.
- `internal/usecase/auth/`: `login_administrator.go`, `authenticate_session.go`, `get_current_session.go`, `logout_administrator.go`; ajuste de una línea en `verify_administrator_setup.go`.
- `internal/adapter/db/postgres/model/administrator.go` (+`LastAcceptedTOTPPeriod`); migración `20260822_05_...`; `repository/administrator/{find_by_id,find_by_email,update_last_accepted_totp_period}.go`; `repository/administrator_session/{find_by_token_hash,update_idle_expiry,revoke}.go`.
- `internal/adapter/security/{password_hasher,totp_verifier,session_token_generator}.go` (extensiones).
- `internal/adapter/rest/middleware/auth.go` (primer middleware real del proyecto).
- `internal/adapter/rest/dto/auth/login_request.go` (nuevo; `administrator.go`/`session.go` se reutilizan tal cual).
- `internal/adapter/rest/handler/auth/{login,get_current_session,logout}.go`; `internal/adapter/rest/response/cookie.go` (nombre de cookie compartido entre handler y middleware).
- `internal/adapter/rest/router/router.go`.
- `cmd/main.go`.
- `docs/errors/aaa-map.md`.

## Fuera de alcance

- Recovery (`POST /auth/recovery*`) — PB-064.
- Revocar sesiones anteriores en login (ADR-008 reserva eso para recovery).
- Rate limiting avanzado/hardening adicional — PB-065.
- Validación de `Origin` en mutaciones autenticadas (`DELETE /auth/session`): requiere `AKRITAS_PUBLIC_URL`/`AKRITAS_ALLOWED_ORIGINS`, todavía no implementados en `config.go`; queda como gap documentado, no como decisión silenciosa.

## Instrucción para el harness

Primero generar `implementation-brief.md` y `tdd-test-plan.md`. No implementar código hasta aprobación humana.
