# Implementation Brief

## Task

`AKR-AUTH-TOTP-VERIFY` (PB-062): implementar `POST /auth/setup/verify`, que cierra el bootstrap — verifica el TOTP, activa al único `Administrator` y abre la primera sesión.

## Current project context

`AKR-AUTH-BOOTSTRAP` dejó: `PendingEnrollment` en el dominio; `AdministratorRepository` (sólo `ExistsActive`), `PendingEnrollmentRepository` (sólo `Save`), `CredentialStore` (sólo `Encrypt`), `TOTPSecretGenerator`, `PasswordHasher`, `BootstrapTokenVerifier`, `RateLimiter`, `Clock` en `ports/out`; los usecases `GetSetupStatus`/`StartAdministratorSetup`; los adapters de seguridad correspondientes; la primera migración Postgres (`administrators` sin columna TOTP, `pending_enrollments`); DTOs/handler/router para `/auth/setup-status` y `/auth/setup`; `config.Load()` con `AKRITAS_BOOTSTRAP_TOKEN`/`AKRITAS_MASTER_KEY`/`AKRITAS_DB_DSN`.

`docs/openapi.yaml` ya especifica `POST /auth/setup/verify` completo: body `TotpEnrollmentVerificationRequest { enrollment_id (uuid), totp_code (pattern ^[0-9]{6}$) }`, respuesta `200 SessionResponse { data: Session { administrator: Administrator, authenticated_at, idle_expires_at, absolute_expires_at } }` con headers `Set-Cookie` (`SessionCookie`: "Secure HttpOnly opaque session cookie; no token is returned in JSON") y `Cache-Control: no-store`; `400`/`409`/`429`/`500` iguales a `/auth/setup`. `domain.AdministratorSession` (con `Validate`/`IsActive`/`Revoke`) y `domain.Administrator` (id/email/display_name/created_at/updated_at, sin campos de secreto) ya existen y no se tocan.

Decisión confirmada explícitamente con el usuario (no estaba fijada por ningún ADR, igual que el motor de DB en PB-061): el valor de la cookie es un token aleatorio de alta entropía; en `administrator_sessions` sólo se persiste su SHA-256 (`token_hash`). El `AdministratorSession.ID` del dominio sigue siendo el id interno de la fila, no el token.

## Proposed approach

1. Dominio: un sentinel nuevo en `errors.go` (`ErrInvalidTotpEnrollmentVerification`, `0x401008V`, 400) para el caso genérico "enrollment_id inexistente/expirado o totp_code incorrecto" — un único código, sin distinguir la causa (mismo criterio que `ErrInvalidBootstrapToken` en PB-061). Se reutiliza `ErrAdministratorAlreadyExists` (`0x401007C`, 409) para la carrera de doble activación. No se toca `pending_enrollment.go` ni `administrator_session.go`.
2. Ports out: agregar `Decrypt` a `CredentialStore`; agregar `Create` a `AdministratorRepository`; agregar `FindByID`/`Delete` a `PendingEnrollmentRepository`; nuevos `TOTPVerifier` (separado de `TOTPSecretGenerator` por SRP), `AdministratorSessionRepository`, `SessionTokenGenerator`.
3. Port in: `VerifyAdministratorSetupUseCase`.
4. Usecase `VerifyAdministratorSetup`: rate limit → parsear `enrollment_id` → `FindByID` (nil = no encontrado) → `IsExpired` → `Decrypt` el secreto → `TOTPVerifier.Verify` (±1 período) → `ExistsActive` (carrera) → `NewAdministrator` → `AdministratorRepository.Create` (con `password_hash` y el mismo ciphertext ya guardado en el pending enrollment, sin re-cifrar) → `SessionTokenGenerator.Generate` → `NewAdministratorSession` → `AdministratorSessionRepository.Save` → `PendingEnrollmentRepository.Delete`. Cualquier fallo en un paso previo no ejecuta los siguientes; cualquier error de infraestructura se propaga tal cual (mapea a 500 en el handler).
5. Adapters de persistencia: agregar `EncryptedTOTPSecret []byte` a `model.Administrator` (columna nueva vía migración `AddColumn`, no `AutoMigrate`); nuevo `model.AdministratorSession` (`administrator_sessions`, con `token_hash` único); dos migraciones nuevas (`20260822_03`, `20260822_04`); repositorios `administrator.Create`, `pending_enrollment.{FindByID,Delete}`, nuevo `repository/administrator_session`.
6. Adapters de seguridad: `Decrypt` en `aesGCMCredentialStore` (abre el mismo formato `nonce||ciphertext` que ya escribe `Encrypt`); nuevo `totpVerifier` (`pquerna/otp/totp.ValidateCustom`, `Period:30, Skew:1, Digits:6, Algorithm:SHA1`); nuevo `sessionTokenGenerator` (32 bytes random → base64url para la cookie, SHA-256 hex para `token_hash`).
7. REST: DTOs `TotpEnrollmentVerificationRequest`, `Administrator`, `Session`, `SessionResponse` espejando el OpenAPI; handler que valida shape (`enrollment_id` UUID, `totp_code` `^[0-9]{6}$`), invoca el usecase, y en éxito escribe `Set-Cookie` (`akritas_session`, `HttpOnly`, `Secure` según config, `SameSite=Lax`, `Path=/`, expira en `absolute_expires_at`) + `Cache-Control: no-store` + body `SessionResponse`; en error reutiliza el mapeo de errores ya existente (`response.WriteDomainError`/`WriteInternalError`, más el caso rate-limited).
8. Router: registrar `POST /api/v1/auth/setup/verify`.
9. Config: agregar `AKRITAS_SESSION_IDLE_TTL` (default `12h`), `AKRITAS_SESSION_ABSOLUTE_TTL` (default `168h`), `AKRITAS_SESSION_COOKIE_SECURE` (default `true`) — opcionales con default, no entran en la validación de "faltantes" (a diferencia de `AKRITAS_BOOTSTRAP_TOKEN`/`AKRITAS_MASTER_KEY`/`AKRITAS_DB_DSN`), tal como `docs/configuration.md` los presenta bajo "Session defaults" en vez de "Required values".
10. Wiring en `cmd/main.go`.

## Architecture impact

```text
internal/core/ports/out/{credential_store,administrator_repository,pending_enrollment_repository}.go   (extendidos)
internal/core/ports/out/{totp_verifier,administrator_session_repository,session_token_generator}.go     (nuevos)
internal/core/ports/in/verify_administrator_setup.go
internal/usecase/auth/verify_administrator_setup.go
internal/adapter/db/postgres/model/{administrator,administrator_session}.go
internal/adapter/db/postgres/migrations/schema/20260822_0{3,4}_...go
internal/adapter/db/postgres/repository/{administrator/create,pending_enrollment/{find_by_id,delete},administrator_session}.go
internal/adapter/security/{credential_store,totp_verifier,session_token_generator}.go
internal/adapter/rest/{dto,handler,router}/auth/...
```

Dirección de dependencias sin cambios: `adapter/rest|db → usecase → core`. `internal/core/**` sigue sin GORM/Chi/`net/http`/SDKs.

## API/OpenAPI impact

No se modifica `docs/openapi.yaml`: `/auth/setup/verify` ya está completamente especificado desde `AKR-OPENAPI-MVP`. Se implementa contra ese contrato existente.

## Data/persistence impact

Dos migraciones nuevas, sin tocar `20260822_01`/`20260822_02`:

- `20260822_03_add_totp_secret_to_administrators`: `ALTER TABLE administrators ADD COLUMN encrypted_totp_secret bytea NOT NULL` vía `tx.Migrator().AddColumn(&model.Administrator{}, "EncryptedTOTPSecret")` (no `AutoMigrate`). Seguro porque `administrators` está garantizada vacía hasta que esta tarea corra su primer `Create` — no hace falta backfill.
- `20260822_04_create_administrator_sessions`: tabla `administrator_sessions` (`id`, `administrator_id`, `token_hash` único, `authenticated_at`, `idle_expires_at`, `absolute_expires_at`, `revoked_at` nullable).

Ambas con `Rollback` explícito (`DropColumn`/`DropTable`).

## Error handling impact

Un sentinel nuevo: `ErrInvalidTotpEnrollmentVerification` (`0x401008V`, 400) — cubre `enrollment_id` inexistente, expirado, y `totp_code` incorrecto, sin distinguir la causa en la respuesta (mismo criterio ADR-008 que ya se aplicó a `ErrInvalidBootstrapToken`). Se reutiliza `ErrAdministratorAlreadyExists` (409) para la carrera de doble activación, y el sentinel de rate limit ya existente (`auth.ErrSetupRateLimited`, 429). Documentado en `docs/errors/aaa-map.md`.

## Test strategy

- Dominio: test de catálogo del nuevo sentinel (extiende `errors_test.go`).
- Usecase: fakes de los ports out (rate limiter, pending enrollment repo, credential store, totp verifier, administrator repo, session token generator, clock, administrator session repo) cubriendo camino feliz y cada rama de error, con aserciones de que ningún paso posterior se ejecuta tras un fallo temprano.
- Adapters: `Decrypt` es inverso de `Encrypt` (round-trip); `TOTPVerifier` acepta el período actual y ±1, rechaza fuera de tolerancia y códigos incorrectos (vectores RFC 6238 conocidos); `SessionTokenGenerator` genera tokens distintos y su hash es el SHA-256 correcto del token devuelto.
- Repositorios (Postgres real, local, como en PB-061): `AdministratorRepository.Create` persiste y respeta el único email (segundo `Create` con el mismo email falla); `PendingEnrollmentRepository.FindByID`/`Delete`; `AdministratorSessionRepository.Save`.
- REST: handler con usecase fake — 200 + `Set-Cookie` con los atributos correctos + `Cache-Control: no-store` + body; 400 (uuid malformado, código con formato inválido, verificación fallida); 409 (administrator ya existe); 429 (+`Retry-After`); ninguna respuesta de error expone `totp_code`, la cookie ni el secreto TOTP.
- Migraciones: la 03 agrega la columna esperada; la 04 crea la tabla esperada; ambas con `Rollback` verificable.
- Validaciones finales: `go test ./...`, `go vet ./...`, `gofmt -l .`, los tres scripts del harness.

## Risks

- Sin transacción cruzando `AdministratorRepository.Create` + `AdministratorSessionRepository.Save` + `PendingEnrollmentRepository.Delete`: si el proceso falla entre el `Create` y el `Delete`, queda un pending enrollment huérfano (inofensivo — cualquier reintento de `/auth/setup/verify` u `/auth/setup` va a fallar con 409/400 igual, porque `ExistsActive` ya es `true`). Se documenta como límite aceptado en vez de introducir un Unit of Work nuevo para esta tarea.
- Dos verificaciones casi simultáneas con el mismo código válido pueden pasar ambas el chequeo `ExistsActive` antes de que cualquiera haga `Create`; el índice único de `email` en `administrators` evita una segunda fila, pero la request perdedora recibe un 500 genérico en vez de un 409 limpio. Aceptado dado el escenario improbable (un humano tipeando un código una vez) y fuera del alcance de rate limiting avanzado (PB-065).
- No se persiste "último período TOTP aceptado": para este enrollment de un solo uso, el consumo (`Delete`) del pending enrollment ya impide reutilizar el mismo `enrollment_id` una segunda vez. Si PB-063 (login) necesita protección contra reutilización de período en verificaciones repetidas contra el mismo secreto, se agrega ahí — agregarlo ahora sería diseñar para una necesidad que esta tarea no tiene.

## Files likely to change

Ver "Archivos o zonas probablemente afectadas" en `task.md`.
