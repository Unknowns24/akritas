# TDD Test Plan

## Scope

Definir mediante tests el contrato de `POST /auth/setup/verify`: dominio (nuevo sentinel), las extensiones/nuevos out ports, el usecase `VerifyAdministratorSetup`, los adapters de persistencia y seguridad nuevos/extendidos, y el handler+router. No se prueban login, lectura/logout de sesión ni recovery porque no forman parte de esta tarea.

## Tests to add/update

### 1. Dominio — catálogo de errores

- `ErrInvalidTotpEnrollmentVerification` = `0x401008V`, cumple el pattern `^[0-9A-F]x[0-9A-F]{6}[VUFNCI]$`, único en el catálogo, registrado en `DomainErrors()` y en `docs/errors/aaa-map.md`.
- No se agregan ni modifican tests de `PendingEnrollment`/`AdministratorSession` — se reutilizan tal cual.

### 2. Ports — compilación de contratos

- `CredentialStore` gana `Decrypt(ctx, ciphertext []byte) ([]byte, error)`.
- `AdministratorRepository` gana `Create(ctx, *domain.Administrator, passwordHash string, encryptedTOTPSecret []byte) error`.
- `PendingEnrollmentRepository` gana `FindByID(ctx, uuid.UUID) (*domain.PendingEnrollment, error)` (nil, nil = no encontrado) y `Delete(ctx, uuid.UUID) error` (idempotente: borrar un id inexistente no es error).
- Nuevos: `TOTPVerifier{ Verify(secret, code string, at time.Time) (bool, error) }`, `AdministratorSessionRepository{ Save(ctx, *domain.AdministratorSession, tokenHash string) error }`, `SessionTokenGenerator{ Generate() (token, hash string, err error) }`.
- `VerifyAdministratorSetupUseCase` en `ports/in`, con `Input{EnrollmentID, TOTPCode, RateLimitKey string}` y `Output{Administrator domain.Administrator, SessionToken string, AuthenticatedAt, IdleExpiresAt, AbsoluteExpiresAt time.Time}`.

### 3. Usecase — `VerifyAdministratorSetup`

Con fakes de los 7 out ports (rate limiter, pending enrollment repo, credential store, totp verifier, administrator repo, session token generator, clock, administrator session repo — 8 en total):

- Camino feliz: rate limit ok → `FindByID` devuelve un enrollment válido y no expirado → `Decrypt` se llama con `enrollment.EncryptedTOTPSecret` → `TOTPVerifier.Verify` se llama con el secreto descifrado y el código recibido → `ExistsActive` false → se construye `domain.NewAdministrator` con el email/display_name del enrollment → `AdministratorRepository.Create` se llama con `enrollment.PasswordHash` y **el mismo `enrollment.EncryptedTOTPSecret`** (no un nuevo cifrado) → `SessionTokenGenerator.Generate` → `NewAdministratorSession` con `IdleTTL`/`AbsoluteTTL` inyectados al usecase → `AdministratorSessionRepository.Save` se llama con el `tokenHash` devuelto por el generador → `PendingEnrollmentRepository.Delete` se llama con el id del enrollment consumido. El `Output` expone el token en claro (nunca el hash) y los tres timestamps de sesión.
- Rate limited → error de rate limit; ningún otro port se llama.
- `EnrollmentID` no parseable como UUID → `ErrInvalidTotpEnrollmentVerification`; `FindByID` no se llama.
- `FindByID` devuelve `(nil, nil)` (no encontrado) → `ErrInvalidTotpEnrollmentVerification`; no se llama `Decrypt` ni `TOTPVerifier`.
- Enrollment encontrado pero `IsExpired(now)` → `ErrInvalidTotpEnrollmentVerification`; no se llama `Decrypt` ni `TOTPVerifier`.
- `TOTPVerifier.Verify` devuelve `false` → `ErrInvalidTotpEnrollmentVerification`; no se llama `ExistsActive` ni ningún port posterior.
- `ExistsActive` devuelve `true` (carrera) → `ErrAdministratorAlreadyExists`; no se llama `Create`, `SessionTokenGenerator`, `Save` ni `Delete`.
- Errores de infraestructura en cualquier port (`Decrypt`, `Create`, `Save`, `Delete`, etc.) se propagan tal cual, sin envolver en un sentinel de dominio, y detienen los pasos siguientes.

### 4. Adapters de seguridad

- `CredentialStore.Decrypt` es inverso de `Encrypt`: cifrar y luego descifrar devuelve el plaintext original, para varios plaintexts distintos. Un ciphertext corrupto/demasiado corto devuelve error (no panic).
- `TOTPVerifier.Verify`: acepta el código válido para el período actual; acepta el código del período anterior y del siguiente (`Skew: 1`); rechaza un código de dos períodos de distancia; rechaza un código de 6 dígitos incorrecto para el secreto/momento dados. Se usa el mismo generador (`pquerna/otp/totp.GenerateCode`) para construir los códigos esperados en cada caso, evitando vectores hardcodeados frágiles.
- `SessionTokenGenerator.Generate`: dos llamadas devuelven tokens distintos; `hash` es el SHA-256 hex del `token` devuelto en la misma llamada; el `token` tiene suficiente entropía (longitud decodificada ≥ 32 bytes).

### 5. Adapters de persistencia (Postgres, real, local)

- `AdministratorRepository.Create`: persiste un `Administrator` con `password_hash` y `encrypted_totp_secret`; un segundo `Create` con el mismo email falla (constraint único), sin dejar una fila parcial.
- `PendingEnrollmentRepository.FindByID`: devuelve el enrollment guardado por `Save`; devuelve `(nil, nil)` para un id inexistente.
- `PendingEnrollmentRepository.Delete`: borra la fila; llamarlo de nuevo sobre el mismo id (ya borrado) no devuelve error.
- `AdministratorSessionRepository.Save`: persiste la sesión con su `token_hash`; el `administrator_id` referencia al `Administrator` recién creado.
- Migraciones `20260822_03`/`20260822_04`: la columna `encrypted_totp_secret` existe en `administrators` tras migrar; la tabla `administrator_sessions` existe con sus columnas; ambas migraciones no tocan ni redefinen `20260822_01`/`20260822_02`.

### 6. REST — DTOs, handler y router

- `TotpEnrollmentVerificationRequest` deserializa `enrollment_id`/`totp_code`.
- Body malformado o `totp_code` que no matchea `^[0-9]{6}$` o `enrollment_id` no-UUID → `400`, sin invocar el usecase.
- Usecase feliz → `200`, header `Set-Cookie` presente con `HttpOnly`, `SameSite=Lax`, `Path=/` (y `Secure` según config), header `Cache-Control: no-store`, body `SessionResponse` con `data.administrator.{id,email,display_name,created_at,updated_at}` y `data.{authenticated_at,idle_expires_at,absolute_expires_at}` — sin `revoked_at` ni `administrator_id` (esos campos son internos de `domain.AdministratorSession`, no de `Session` del OpenAPI).
- Usecase devuelve `ErrInvalidTotpEnrollmentVerification` → `400`.
- Usecase devuelve `ErrAdministratorAlreadyExists` → `409`.
- Usecase devuelve el error de rate limit → `429` + `Retry-After`.
- Usecase devuelve un error no reconocido → `500` genérico, sin filtrar la causa.
- Ninguna respuesta (éxito o error) incluye `totp_code` del request, el secreto TOTP, ni el valor crudo de la cookie en el body JSON.
- Router: `POST /api/v1/auth/setup/verify` responde a través del router real (test de integración liviano, como los de `/auth/setup`).

### 7. Config

- `AKRITAS_SESSION_IDLE_TTL`/`AKRITAS_SESSION_ABSOLUTE_TTL`/`AKRITAS_SESSION_COOKIE_SECURE` ausentes → `Load()` usa los defaults (`12h`/`168h`/`true`) sin fallar (a diferencia de las tres variables ya requeridas).
- Valores presentes y válidos → se parsean correctamente (`time.ParseDuration`/`strconv.ParseBool`).
- Valor inválido (duración o bool no parseable) → `Load()` falla con un mensaje que nombra la variable, sin imprimir su valor.

### 8. Validaciones finales

- `go test ./...`.
- `go vet ./...`.
- `gofmt -l .` sin diferencias.
- `.harness/kernel/scripts/check-backend-architecture.sh`.
- `.harness/kernel/scripts/check-openapi.sh` sin modificar `docs/openapi.yaml`.
- `.harness/kernel/scripts/check-security.sh`.
- Verificación manual end-to-end contra Postgres local (como en PB-061): `setup` → `setup/verify` con el código real generado a partir del `manual_entry_key` devuelto → `200` + cookie; código incorrecto → `400`; `enrollment_id` inexistente → `400`; segundo intento tras éxito (enrollment ya consumido) → `400`; inspección de logs y respuestas sin secretos.

## Expected failing tests before implementation

- No existe `ErrInvalidTotpEnrollmentVerification`.
- `CredentialStore` no tiene `Decrypt`; `AdministratorRepository` no tiene `Create`; `PendingEnrollmentRepository` no tiene `FindByID`/`Delete`.
- No existen `TOTPVerifier`, `AdministratorSessionRepository`, `SessionTokenGenerator` ni `VerifyAdministratorSetupUseCase`.
- No existe `internal/usecase/auth/verify_administrator_setup.go`.
- No existen `model.AdministratorSession`, las migraciones `20260822_03`/`04`, ni sus repositorios.
- No existen los adapters `totp_verifier.go`/`session_token_generator.go`, ni `Decrypt` en `credential_store.go`.
- No existen los DTOs `TotpEnrollmentVerificationRequest`/`Administrator`/`Session`/`SessionResponse`, el handler ni la ruta `/auth/setup/verify`.
- `config.Config` no tiene los tres campos de sesión.

Los tests se escribirán después de la aprobación de este plan y deberán fallar (o no compilar) antes de incorporar la implementación correspondiente.

## Acceptance criteria covered

- El endpoint devuelve exactamente los schemas y códigos de estado del OpenAPI.
- El secreto TOTP nunca se descifra fuera del punto de uso ni se devuelve por API.
- La cookie de sesión nunca se persiste en claro (sólo `token_hash`).
- `internal/core/**` no importa GORM, Chi, `net/http` ni SDKs externos.
- El sentinel nuevo cumple `0x4XXNNNT`; los reutilizados no se duplican.
- Los tres scripts del harness pasan sin modificar `docs/openapi.yaml`.
- Las migraciones `20260822_01`/`02` quedan intactas.

## Open questions / human approval notes — resueltas

- Representación del token de sesión: **token aleatorio + hash SHA-256 en DB**, confirmado explícitamente por el usuario antes de este plan.
- Ausencia de "último período TOTP aceptado" persistido: decisión de esta tarea (documentada en `implementation-brief.md`, sección Risks) — el consumo del pending enrollment ya protege contra reutilización del mismo `enrollment_id`; si hace falta contra logins repetidos, se agrega en PB-063.
- Ausencia de transacción cruzando `Create`+`Save`+`Delete`: riesgo aceptado y documentado (huérfano inofensivo en el peor caso), no se introduce una capa de Unit of Work nueva para esta tarea.

No quedan decisiones abiertas. Se requiere aprobación humana explícita de este archivo antes de crear tests o implementar código.
