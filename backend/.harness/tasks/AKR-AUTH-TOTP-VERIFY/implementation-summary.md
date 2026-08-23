# Implementation Summary

## Implemented

- Un sentinel nuevo en `internal/core/domain/errors.go`: `ErrInvalidTotpEnrollmentVerification` (`0x401008V`, 400) — cubre enrollment inexistente/expirado y código incorrecto sin distinguir la causa, registrado en `DomainErrors()` y `docs/errors/aaa-map.md`.
- Ports extendidos: `CredentialStore.Decrypt`, `AdministratorRepository.Create` (mapea la violación del índice único de `email` a `ErrAdministratorAlreadyExists`), `PendingEnrollmentRepository.{FindByID,Delete}`. Ports nuevos: `TOTPVerifier`, `AdministratorSessionRepository`, `SessionTokenGenerator`, `Transactor`, y el in port `VerifyAdministratorSetupUseCase`.
- Usecase `VerifyAdministratorSetup` en `internal/usecase/auth`: rate limit → parsear `enrollment_id` → `FindByID` → `IsExpired` → `Decrypt` → `TOTPVerifier.Verify` (±1 período RFC 6238) → `ExistsActive` (chequeo temprano) → `NewAdministrator` → **transacción** (`AdministratorRepository.Create` + `AdministratorSessionRepository.Save`) → `PendingEnrollmentRepository.Delete`. El mismo ciphertext ya guardado en el pending enrollment se persiste tal cual en `administrators.encrypted_totp_secret`, sin volver a cifrar.
- Persistencia: `model.Administrator` gana `EncryptedTOTPSecret`; nuevo `model.AdministratorSession` (`administrator_sessions`, `token_hash` único). Migraciones `20260822_03` (idempotente — ver Risks) y `20260822_04`, sin tocar `01`/`02`. Repositorios `administrator.Create`, `pending_enrollment.{FindByID,Delete}`, nuevo `repository/administrator_session`. `txcontext` + `Transactor` (GORM) propagan una transacción compartida entre `Create` y `Save`.
- Adapters de seguridad: `Decrypt` en `aesGCMCredentialStore` (inverso de `Encrypt`); `totpVerifier` (`pquerna/otp/totp.ValidateCustom`, `Period:30, Skew:1, Digits:6, Algorithm:SHA1`); `sessionTokenGenerator` (32 bytes random → base64url para la cookie, SHA-256 hex para `token_hash` — decisión de representación de sesión confirmada explícitamente con el usuario).
- REST: DTOs `TotpEnrollmentVerificationRequest`, `Administrator`, `Session`, `SessionResponse`; handler `VerifyAdministratorSetup` (valida `enrollment_id` UUID + `totp_code` `^[0-9]{6}$`, setea `Set-Cookie` `akritas_session` HttpOnly/Secure-según-config/SameSite=Lax/Path=/ expirando en `absolute_expires_at`, `Cache-Control: no-store`); ruta `POST /api/v1/auth/setup/verify`.
- Config: `AKRITAS_SESSION_IDLE_TTL` (default 12h), `AKRITAS_SESSION_ABSOLUTE_TTL` (default 168h), `AKRITAS_SESSION_COOKIE_SECURE` (default true) — opcionales con default, per `docs/configuration.md`.
- Wiring completo en `cmd/main.go`.

## Deliberately not implemented

- Login con password (`POST /auth/login`) — PB-063.
- Lectura de sesión actual y logout (`GET`/`DELETE /auth/session`) — PB-063.
- Recovery (`POST /auth/recovery*`) — PB-064.
- Rate limiting avanzado más allá de reutilizar `out.RateLimiter` — PB-065.
- Persistencia de "último período TOTP aceptado": el consumo del pending enrollment ya protege contra reutilización del mismo `enrollment_id` en este flujo de un solo uso; decisión aprobada explícitamente por el usuario.

## Validation result

- `go test ./...`: pass.
- `go test -race ./...`: pass.
- Coverage: domain 82.2%, usecase/auth 91.5%, adapter/security 88.9%, adapter/db/postgres (Transactor) 100%, migrations 100%, migrations/schema 42.9% (las dos migraciones sin código condicional, 01/02, no ejercen ramas nuevas; la rama idempotente de 03 sí está cubierta), repository/administrator 85.7%, repository/administrator_session 100%, repository/pending_enrollment 85.7%, adapter/rest/handler/auth 98.6%, adapter/rest/response 86.7%, adapter/rest/router 100%, config 91.7%.
- `go vet ./...`: pass.
- `gofmt -l .`: sin diferencias.
- `check-backend-architecture.sh`, `check-openapi.sh` (59 operaciones, 112 schemas, sin cambios al spec), `check-security.sh`: pass.
- Verificación manual end-to-end contra Postgres local (no simulada): `setup` → generación real del código TOTP a partir del `manual_entry_key` devuelto → `setup/verify` con código incorrecto (400), `enrollment_id` inexistente (400), código correcto (200 + `Set-Cookie` con todos los atributos correctos + body `SessionResponse`), `setup-status` post-activación (`false`/`false`), reintento de verify con el mismo enrollment ya consumido (400). Inspección directa de la base: `administrators` con 1 fila (`encrypted_totp_secret` de 60 bytes, no el secreto en claro), `pending_enrollments` vacía, `administrator_sessions` con 1 fila cuyo `token_hash` es el SHA-256 del token de la cookie, no el token mismo. Logs del proceso sin secretos.
- Esta verificación manual encontró dos bugs reales:
  1. En una base nueva, la migración `20260822_03` fallaba porque `20260822_01`'s `AutoMigrate` ya crea la columna nueva (refleja el struct Go actual, no uno congelado). Corregido haciendo `Migrate` idempotente (`HasColumn` antes de `AddColumn`), con test de regresión (`TestMigration03IsIdempotentWhenColumnAlreadyExists`).
  2. `cmd/main.go` pasaba la misma instancia de `RateLimiter` a `StartAdministratorSetup` y a `VerifyAdministratorSetup`, compartiendo un único budget de 5/15min por IP entre ambos endpoints — unos pocos reintentos de `verify` (typos, desfasaje de reloj) podían agotar el budget que necesitaba `setup`, y viceversa. Corregido con dos instancias independientes; confirmado manualmente que agotar el budget de `verify` deja `setup`/`setup-status` totalmente utilizables.
