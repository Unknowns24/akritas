# Implementation Brief

## Task

`AKR-AUTH-LOGIN-SESSION` (PB-063): implementar `POST /auth/login`, `GET /auth/session` y `DELETE /auth/session` — el resto del ciclo de autenticación sobre el `Administrator` ya activado por PB-062.

## Current project context

PB-062 dejó: `Administrator` persistido con `password_hash` y `encrypted_totp_secret`; `AdministratorSession` (dominio con `Validate`/`IsActive`/`Revoke`) persistida en `administrator_sessions` con `token_hash`; `TOTPVerifier`, `PasswordHasher`, `CredentialStore`, `SessionTokenGenerator`, `RateLimiter`, `Transactor`, `Clock` en `ports/out`; `docs/openapi.yaml` ya especifica `POST /auth/login` y `GET`/`DELETE /auth/session` completos (schemas `LoginRequest`, `SessionResponse`, `Session`, `Administrator`, header `SessionCookie`, security scheme `cookieAuth` con nombre de cookie `akritas_session`). `internal/adapter/rest/middleware/` está vacío — esta es la primera tarea que necesita autenticación real de un request entrante.

Decisión confirmada explícitamente con el usuario: el idle TTL **se desliza** en cada request autenticada exitosa (se extiende a `now + idle_ttl`, sin superar nunca `absolute_expires_at`).

## Proposed approach

1. **Dominio**: un sentinel nuevo (`ErrInvalidCredentials`, `0x401009U`, 401) para el fallo genérico de login (ADR-008: no distinguir email/password/TOTP). Nuevo método `AdministratorSession.ExtendIdle(now, idleTTL) error` (aditivo, no toca `Validate`/`IsActive`/`Revoke`): valida que la sesión siga activa, calcula `now+idleTTL` y lo capa a `AbsoluteExpiresAt` si lo supera.
2. **Ports out extendidos**:
   - `AdministratorRepository`: `FindByID` (solo lectura, para `GET /auth/session`); `FindByEmail` (devuelve un `AdministratorCredentials{Administrator, PasswordHash, EncryptedTOTPSecret, LastAcceptedTOTPPeriod}` — mismo criterio que `Create`: el dominio no carga secretos); `UpdateLastAcceptedTOTPPeriod`.
   - `AdministratorSessionRepository`: `FindByTokenHash` (nil,nil si no existe); `UpdateIdleExpiry`; `Revoke`.
   - `PasswordHasher`: `Verify(password, hash) (bool, error)`.
   - `SessionTokenGenerator`: `Hash(token) string` (determinístico, sin error — reutiliza el mismo algoritmo que ya usa `Generate` internamente, necesario para que el middleware pueda hashear el valor de la cookie de forma independiente).
   - `TOTPVerifier.Verify` gana un tercer valor de retorno: el período RFC 6238 que matcheó (`(valid bool, period int64, err error)`). Es el único cambio de firma de un port existente en esta tarea; el call site de PB-062 (`verify_administrator_setup.go`) sólo necesita ignorar el nuevo valor (`valid, _, err := ...`), sin cambiar su comportamiento — PB-062 decidió explícitamente no perseguir protección de reutilización de período porque el pending enrollment es de un solo uso; esta tarea sí la necesita porque el `Administrator` persiste entre logins.
3. **Ports in nuevos**: `LoginAdministratorUseCase`; `AuthenticateSessionUseCase` (resuelve+valida+desliza una sesión a partir del token crudo de la cookie — lo usa el middleware); `GetCurrentSessionUseCase` (toma la sesión ya resuelta, arma la proyección); `LogoutAdministratorUseCase` (toma la sesión ya resuelta, la revoca).
4. **Usecases**:
   - `LoginAdministrator`: `RateLimiter.Allow` con key `"ip:"+RateLimitKey` y con key `"account:"+Email` (ambas deben pasar) → `FindByEmail` (nil ⇒ `ErrInvalidCredentials`) → `PasswordHasher.Verify` (false ⇒ `ErrInvalidCredentials`) → `CredentialStore.Decrypt` → `TOTPVerifier.Verify` (false ⇒ `ErrInvalidCredentials`) → si `period == LastAcceptedTOTPPeriod` ⇒ `ErrInvalidCredentials` (reutilización) → `SessionTokenGenerator.Generate` → `NewAdministratorSession` (sesión **nueva**, no se tocan sesiones previas) → dentro de una transacción: `UpdateLastAcceptedTOTPPeriod` + `AdministratorSessionRepository.Save`.
   - `AuthenticateSession`: token vacío ⇒ `ErrInactiveAdministratorSession` sin tocar la DB → `SessionTokenGenerator.Hash` → `FindByTokenHash` (nil ⇒ `ErrInactiveAdministratorSession`) → `IsActive(now)` (false ⇒ `ErrInactiveAdministratorSession`) → `ExtendIdle` + `UpdateIdleExpiry` (desliza) → devuelve la sesión ya extendida.
   - `GetCurrentSession`: toma la sesión resuelta por el middleware, `AdministratorRepository.FindByID` para armar `Administrator`, arma la proyección con los timestamps ya deslizados.
   - `LogoutAdministrator`: toma la sesión resuelta, `session.Revoke(now)`, `AdministratorSessionRepository.Revoke`.
5. **Persistencia**: `model.Administrator` gana `LastAcceptedTOTPPeriod int64` (`default:0`); migración `20260822_05` idempotente (`HasColumn`) igual que la corrección de PB-062. Repositorios nuevos/extendidos correspondientes; los de solo-lectura (`FindByID`, `FindByEmail`, `FindByTokenHash`) no participan de transacciones; `UpdateLastAcceptedTOTPPeriod` sí (dentro de la transacción de login, junto a `Save`), igual patrón que PB-062.
6. **Seguridad**: `PasswordHasher.Verify` usa `argon2.IDKey` con los mismos parámetros embebidos en el hash (ya versionados desde PB-061) y compara en tiempo constante; `totpVerifier.Verify` deja de delegar en `totp.ValidateCustom` (que no expone qué período matcheó) y en su lugar prueba los tres períodos candidatos (`-1,0,+1`) con `totp.GenerateCodeCustom`, comparando cada candidato con `crypto/subtle.ConstantTimeCompare` (mismo criterio anti-timing que el bootstrap token de PB-061), devolviendo el contador del que matcheó.
7. **Middleware**: `internal/adapter/rest/middleware/auth.go` — `RequireSession(in.AuthenticateSessionUseCase) func(http.Handler) http.Handler`: lee la cookie `response.SessionCookieName` (se extrae esa constante, hoy privada en `handler/auth`, a `internal/adapter/rest/response` para que handler y middleware la compartan sin acoplarse entre sí), invoca el usecase, en error escribe la respuesta (mismo mapeo `WriteDomainError`/`WriteInternalError` ya establecido) y corta la cadena; en éxito inyecta `domain.AdministratorSession` en el contexto (`SessionFromContext`).
8. **REST**: DTO `LoginRequest` nuevo; `Administrator`/`Session`/`SessionResponse` se reutilizan sin cambios. Handlers `Login`, `GetCurrentSession`, `Logout`. Router: `POST /auth/login` público; `GET`/`DELETE /auth/session` envueltos en `RequireSession` vía `r.Group`.
9. **Wiring**: `cmd/main.go` arma un `RateLimiter` propio para login (nunca compartir instancia entre endpoints, lección de la corrección de PB-062), construye los 4 usecases nuevos y pasa `AuthenticateSessionUseCase` al router para el middleware.

## Architecture impact

```text
internal/core/ports/out/{administrator_repository,administrator_session_repository,password_hasher,totp_verifier,session_token_generator}.go   (extendidos)
internal/core/ports/in/{login_administrator,authenticate_session,get_current_session,logout_administrator}.go                                  (nuevos)
internal/usecase/auth/{login_administrator,authenticate_session,get_current_session,logout_administrator}.go
internal/adapter/db/postgres/model/administrator.go   (extendido)
internal/adapter/db/postgres/migrations/schema/20260822_05_...go
internal/adapter/db/postgres/repository/{administrator,administrator_session}/*.go   (nuevos archivos, un método por archivo)
internal/adapter/security/{password_hasher,totp_verifier,session_token_generator}.go (extendidos)
internal/adapter/rest/middleware/auth.go
internal/adapter/rest/{dto,handler}/auth/...
```

Dirección de dependencias sin cambios: `adapter/rest|db → usecase → core`. El middleware depende de `ports/in` (no de un usecase concreto), igual que los handlers.

## API/OpenAPI impact

No se modifica `docs/openapi.yaml`: los tres endpoints ya están completamente especificados desde `AKR-OPENAPI-MVP`. `Session` en la respuesta no incluye `administrator_id` ni `revoked_at` (campos internos de `domain.AdministratorSession`), igual que en PB-062.

## Data/persistence impact

Una migración nueva, sin tocar `01`-`04`:

- `20260822_05_add_last_accepted_totp_period_to_administrators`: `ALTER TABLE administrators ADD COLUMN last_accepted_totp_period bigint NOT NULL DEFAULT 0`, con el mismo chequeo `HasColumn` que la corrección de PB-062 (en una base nueva, la migración `01` ya crea la columna reflejando el struct Go actual).

## Error handling impact

Un sentinel nuevo: `ErrInvalidCredentials` (`0x401009U`, 401) — cubre email inexistente, password incorrecta, TOTP incorrecto o reutilizado en login, sin distinguir la causa (ADR-008). Se reutiliza `ErrInactiveAdministratorSession` (`0x401003U`, ya existente desde `AKR-BACKEND-FOUNDATION` pero nunca antes conectado a un flujo real) para cookie ausente/sesión no encontrada/expirada/revocada en `GET`/`DELETE /auth/session`. Se reutiliza el patrón de rate limit ya establecido (sentinel de paquete, fuera del catálogo de dominio, 429 + `Retry-After`), con un sentinel propio (`ErrLoginRateLimited`) distinto de `ErrSetupRateLimited` para no mezclar el vocabulario de dos flujos distintos.

## Test strategy

- Dominio: `ExtendIdle` (extiende, capa al absoluto, rechaza sobre sesión inactiva) y el sentinel nuevo en el test de catálogo.
- Ports: compilación de las firmas extendidas/nuevas.
- Usecases: fakes de cada out port para las cuatro operaciones nuevas, cubriendo camino feliz y cada rama de error (rate limit por IP, por cuenta, credenciales inexistentes/incorrectas/TOTP inválido/período reutilizado para login; sesión ausente/expirada/revocada para authenticate; propagación de errores de infraestructura en cada paso).
- Adapters: `PasswordHasher.Verify` contra hashes reales generados por `Hash`; `totpVerifier.Verify` con vectores generados con `totp.GenerateCode` (período actual, ±1, fuera de tolerancia, código incorrecto) verificando también el período devuelto; `SessionTokenGenerator.Hash` determinístico y coincidente con el hash que ya devuelve `Generate`.
- Repositorios (Postgres real, local, como en PB-061/062): `FindByID`/`FindByEmail` (incluyendo "no encontrado"), `UpdateLastAcceptedTOTPPeriod`, `FindByTokenHash`, `UpdateIdleExpiry`, `Revoke`.
- REST: handler de login (200 + cookie + body; 401 sin distinguir causa; 429 + `Retry-After`); middleware (401 sin cookie, con cookie inválida/expirada/revocada; 200 con sesión válida, inyectando la sesión en el contexto); `GetCurrentSession`/`Logout` con la sesión ya inyectada (200 con proyección correcta; 204 + cookie expirada).
- Validaciones finales: `go test ./...`, `go vet ./...`, `gofmt -l .`, los tres scripts del harness.
- Verificación manual end-to-end contra Postgres local (como en PB-061/062): login con TOTP real generado a partir del secreto activado en PB-062, reutilización del mismo código rechazada, `GET /auth/session` deslizando el idle TTL (comparar `idle_expires_at` entre dos llamadas), `DELETE /auth/session` seguido de un `GET` que ya no debe funcionar.

## Risks

- El middleware no valida `Origin` en `DELETE /auth/session` (mutación autenticada) porque `AKRITAS_PUBLIC_URL`/`AKRITAS_ALLOWED_ORIGINS` todavía no existen en `config.go` — gap documentado, no decisión silenciosa; candidato natural para PB-065 o una tarea de hardening dedicada.
- El idle TTL deslizante implica un `UPDATE` por cada request autenticada exitosa (incluida cada llamada a `GET /auth/session` y la que precede a `DELETE /auth/session`); es el costo esperado de que "idle" signifique inactividad real, aceptado como parte de la decisión ya confirmada.
- Cambiar la firma de `TOTPVerifier.Verify` toca código ya aprobado de PB-062 (`verify_administrator_setup.go`, sus tests, y los fakes compartidos de `internal/usecase/auth`) — cambio mecánico de una línea por call site, sin alterar comportamiento; se señala explícitamente para que la revisión lo confirme.

## Files likely to change

Ver "Archivos o zonas probablemente afectadas" en `task.md`.
