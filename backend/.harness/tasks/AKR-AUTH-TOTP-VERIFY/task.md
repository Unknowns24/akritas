# AKR-AUTH-TOTP-VERIFY - Enrollment y verificación TOTP

## Estado

complete

## Tipo de tarea

backend-api-feature

## Modo de proyecto

existing_project

## Contexto

`AKR-AUTH-BOOTSTRAP` (PB-061) implementó `GET /auth/setup-status` y `POST /auth/setup`: crea un pending enrollment de corta duración con un secreto TOTP cifrado, sin persistir todavía al `Administrator`. Esta tarea cierra el flujo de bootstrap: verifica el código TOTP contra ese secreto, activa al único `Administrator`, y abre la primera sesión del sistema.

## Objetivo

Implementar `POST /auth/setup/verify` (`verifyAdministratorSetup`): valida `enrollment_id` + `totp_code` contra el pending enrollment correspondiente (RFC 6238, tolerancia ±1 período, ADR-008), y si es válido persiste al `Administrator` (primer `INSERT` real de la tabla), guarda su secreto TOTP cifrado para logins futuros, consume el pending enrollment, crea una `AdministratorSession` opaca y responde `SessionResponse` con `Set-Cookie`.

## Requerimiento funcional

- `POST /auth/setup/verify` valida `TotpEnrollmentVerificationRequest { enrollment_id, totp_code }`.
- El código TOTP se verifica contra el secreto (descifrado sólo en el punto de uso, ADR-005) del pending enrollment indicado por `enrollment_id`, con tolerancia de un período RFC 6238 anterior/posterior.
- `enrollment_id` inexistente, expirado, o `totp_code` incorrecto responden `400` de forma genérica (sin distinguir cuál falló, mismo espíritu que ADR-008 para `/auth/setup`).
- Si ya existe un `Administrator` activo al momento de verificar (carrera entre dos intentos), responde `409` (reutiliza `ErrAdministratorAlreadyExists`).
- Verificación exitosa: persiste `Administrator` (email/display_name del pending enrollment), guarda `password_hash` y el secreto TOTP cifrado (el mismo ciphertext ya generado en `/auth/setup`, sin volver a cifrar), consume/borra el pending enrollment usado, crea `AdministratorSession` (idle 12h / absoluto 168h por defecto, configurables), y responde `200 SessionResponse { data: Session }` con `Set-Cookie` (`HttpOnly`, `Secure` en producción, `SameSite=Lax`, `Path=/`) y `Cache-Control: no-store`.
- Rate limited (reutiliza `out.RateLimiter` tal como en `/auth/setup`) → `429` con `Retry-After`.
- El endpoint es público (`security: []`): la única sesión existente todavía no existe en este punto del flujo.

## Criterios de aceptación

- `go test ./...`, `go vet ./...` y `gofmt -l .` (sin diffs) finalizan correctamente.
- El endpoint devuelve exactamente los schemas y códigos de estado definidos en `docs/openapi.yaml` para `/auth/setup/verify`.
- El secreto TOTP nunca se descifra fuera del punto de uso (verificación) y nunca se devuelve por API.
- El valor de la cookie de sesión nunca se persiste en claro: sólo su hash SHA-256 en `administrator_sessions.token_hash`.
- `internal/core/**` no importa GORM, Chi, `net/http` ni SDKs externos.
- Los errores nuevos cumplen `0x4XXNNNT`; se reutilizan los sentinels existentes donde aplica (no se duplican).
- Los checks de arquitectura, OpenAPI y seguridad del harness pasan sin modificar `docs/openapi.yaml`.
- Las migraciones `20260822_01`/`20260822_02` no se modifican; el secreto TOTP y la tabla de sesiones llegan por migraciones nuevas.

## Restricciones técnicas

- Profile: `.harness/kernel/profiles/go-hexagonal-api.yaml`.
- Workflow: `.harness/kernel/workflows/backend-api-feature.yaml`.
- Representación del token de sesión: token aleatorio de alta entropía en la cookie; sólo su SHA-256 se persiste (`administrator_sessions.token_hash`) — decisión confirmada explícitamente con el usuario.
- TOTP: `github.com/pquerna/otp/totp.ValidateCustom` (ya dependencia del proyecto), `Period: 30, Skew: 1, Digits: 6, Algorithm: SHA1`.
- No modificar `internal/core/domain/administrator.go` ni `administrator_session.go`: se reutilizan tal cual (`NewAdministrator`, `NewAdministratorSession`, `Validate`, `IsActive`, `Revoke`).
- No modificar las migraciones `20260822_01_create_administrators` ni `20260822_02_create_pending_enrollments`.
- No implementar código antes de la aprobación humana de `tdd-test-plan.md`.

## Archivos o zonas probablemente afectadas

- `internal/core/domain/errors.go` (un sentinel nuevo, `0x401008V`).
- `internal/core/ports/out/`: `credential_store.go` (+Decrypt), `administrator_repository.go` (+Create), `pending_enrollment_repository.go` (+FindByID, +Delete), nuevos `totp_verifier.go`, `administrator_session_repository.go`, `session_token_generator.go`.
- `internal/core/ports/in/verify_administrator_setup.go`.
- `internal/usecase/auth/verify_administrator_setup.go`.
- `internal/adapter/db/postgres/model/administrator.go` (+`EncryptedTOTPSecret`), nuevo `model/administrator_session.go`.
- `internal/adapter/db/postgres/migrations/schema/20260822_03_...`, `20260822_04_...`, `migrations/migrate.go`.
- `internal/adapter/db/postgres/repository/administrator/create.go`, `repository/pending_enrollment/{find_by_id,delete}.go`, nuevo `repository/administrator_session/`.
- `internal/adapter/security/`: `credential_store.go` (+Decrypt), nuevos `totp_verifier.go`, `session_token_generator.go`.
- `internal/adapter/rest/dto/auth/`, `internal/adapter/rest/handler/auth/verify_administrator_setup.go`, `internal/adapter/rest/router/router.go`.
- `config/config.go` (+`AKRITAS_SESSION_IDLE_TTL`, `AKRITAS_SESSION_ABSOLUTE_TTL`, `AKRITAS_SESSION_COOKIE_SECURE`, todos con default).
- `cmd/main.go`.
- `docs/errors/aaa-map.md`.

## Fuera de alcance

- Login con password (`POST /auth/login`, PB-063).
- Lectura de sesión actual y logout (`GET`/`DELETE /auth/session`, PB-063).
- Recovery (`POST /auth/recovery*`, PB-064).
- Rate limiting avanzado (PB-065) — alcanza con reutilizar `out.RateLimiter` tal como en `/auth/setup`.
- Persistencia del "último período TOTP aceptado" para detectar reutilización en logins futuros: para un enrollment de un solo uso alcanza con que el consumo del pending enrollment sea la protección contra replay; ese mecanismo (si hace falta contra logins repetidos) es de PB-063, no de esta tarea.

## Instrucción para el harness

Primero generar `implementation-brief.md` y `tdd-test-plan.md`. No implementar código hasta aprobación humana.
