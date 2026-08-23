# TDD Test Plan

## Scope

Definir mediante tests el contrato de `POST /auth/login`, `GET /auth/session` y `DELETE /auth/session`: dominio (`ExtendIdle`, nuevo sentinel), los out ports extendidos/nuevos, los cuatro usecases nuevos, los adapters de seguridad extendidos, los repositorios nuevos/extendidos, el middleware de autenticación, y los tres handlers+router. No se prueban recovery ni rate limiting avanzado porque no forman parte de esta tarea.

## Tests to add/update

### 1. Dominio

- `AdministratorSession.ExtendIdle(now, idleTTL)`: sobre una sesión activa, `IdleExpiresAt` pasa a `now+idleTTL`; si `now+idleTTL` supera `AbsoluteExpiresAt`, queda capado en `AbsoluteExpiresAt` (no lo supera nunca); sobre una sesión ya inactiva (expirada o revocada) devuelve error y no modifica el struct.
- Catálogo: `ErrInvalidCredentials` = `0x401009U`, cumple el pattern, único, registrado en `DomainErrors()` y `docs/errors/aaa-map.md`.

### 2. Ports — compilación de contratos

- `AdministratorRepository` gana `FindByID(ctx, uuid.UUID) (*domain.Administrator, error)` (nil,nil si no existe), `FindByEmail(ctx, email string) (*out.AdministratorCredentials, error)` (nil,nil si no existe; `AdministratorCredentials{Administrator, PasswordHash, EncryptedTOTPSecret, LastAcceptedTOTPPeriod}`), `UpdateLastAcceptedTOTPPeriod(ctx, uuid.UUID, int64) error`.
- `AdministratorSessionRepository` gana `FindByTokenHash(ctx, string) (*domain.AdministratorSession, error)` (nil,nil si no existe), `UpdateIdleExpiry(ctx, uuid.UUID, time.Time) error`, `Revoke(ctx, uuid.UUID, time.Time) error`.
- `PasswordHasher` gana `Verify(password, hash string) (bool, error)`.
- `SessionTokenGenerator` gana `Hash(token string) string`.
- `TOTPVerifier.Verify` pasa a `(secret, code string, at time.Time) (valid bool, period int64, err error)`.
- Nuevos in ports: `LoginAdministratorUseCase`, `AuthenticateSessionUseCase`, `GetCurrentSessionUseCase`, `LogoutAdministratorUseCase`.

### 3. Usecase — `LoginAdministrator`

Con fakes de rate limiter, administrator repo, password hasher, credential store, totp verifier, session token generator, administrator session repo, transactor, clock:

- Camino feliz: ambos `Allow` (key `"ip:"+RateLimitKey`, key `"account:"+Email`) → `FindByEmail` → `Verify` password → `Decrypt` → `TOTPVerifier.Verify` (período distinto del `LastAcceptedTOTPPeriod` guardado) → dentro de la transacción, `UpdateLastAcceptedTOTPPeriod` con el período devuelto por `Verify` + `Save` de una sesión **nueva** (no se toca ninguna sesión previa) → `Output` con el token en claro y los tres timestamps.
- `Allow` de IP devuelve `false` → `ErrLoginRateLimited`; no se llama `Allow` de cuenta ni nada posterior.
- `Allow` de cuenta devuelve `false` (IP sí permitido) → `ErrLoginRateLimited`; no se llama `FindByEmail`.
- `FindByEmail` devuelve `(nil, nil)` → `ErrInvalidCredentials`; no se llama `PasswordHasher`/`TOTPVerifier`/nada posterior.
- `PasswordHasher.Verify` devuelve `false` → `ErrInvalidCredentials`; no se llama `Decrypt`/`TOTPVerifier`/nada posterior.
- `TOTPVerifier.Verify` devuelve `valid=false` → `ErrInvalidCredentials`; no se llama `Transactor` ni nada posterior.
- `TOTPVerifier.Verify` devuelve `valid=true` pero `period == LastAcceptedTOTPPeriod` del `AdministratorCredentials` → `ErrInvalidCredentials` (reutilización); no se llama `Transactor`.
- Errores de infraestructura en cualquier paso (`Decrypt`, `SessionTokenGenerator`, `Transactor`, `UpdateLastAcceptedTOTPPeriod`, `Save`) se propagan tal cual y detienen los pasos siguientes.

### 4. Usecase — `AuthenticateSession`

- Token vacío → `ErrInactiveAdministratorSession`; no se llama `Hash` ni `FindByTokenHash`.
- `FindByTokenHash` devuelve `(nil, nil)` → `ErrInactiveAdministratorSession`.
- Sesión encontrada pero `IsActive(now)` es `false` (expirada o revocada) → `ErrInactiveAdministratorSession`; no se llama `ExtendIdle`/`UpdateIdleExpiry`.
- Sesión activa → se llama `ExtendIdle` y `UpdateIdleExpiry` con el nuevo `IdleExpiresAt`; el valor devuelto por el usecase refleja el idle ya extendido.
- Errores de infraestructura (`FindByTokenHash`, `UpdateIdleExpiry`) se propagan tal cual.

### 5. Usecase — `GetCurrentSession` y `LogoutAdministrator`

- `GetCurrentSession`: dada una sesión ya resuelta, `AdministratorRepository.FindByID` arma el `Administrator`; el output expone los timestamps de la sesión recibida (ya deslizados por `AuthenticateSession`, no se recalculan acá). `FindByID` devolviendo error se propaga.
- `LogoutAdministrator`: dada una sesión ya resuelta, llama `session.Revoke(now)` y persiste con `AdministratorSessionRepository.Revoke`; un segundo `Revoke` sobre una sesión ya revocada es idempotente (mismo comportamiento que el método de dominio ya probado en PB-061). Error de `Revoke` del repositorio se propaga.

### 6. Adapters de seguridad

- `PasswordHasher.Verify`: un hash generado por `Hash("password-real")` verifica `true` contra la misma password y `false` contra otra; un hash malformado devuelve error, no panic.
- `totpVerifier.Verify`: acepta el código del período actual (devolviendo su contador), del período anterior y del siguiente (devolviendo el contador correspondiente, no siempre el mismo); rechaza dos períodos de distancia y un código incorrecto (`valid=false`, `period=0`).
- `SessionTokenGenerator.Hash`: determinístico (mismo input → mismo output); coincide exactamente con el `hash` que devuelve `Generate()` para el `token` que generó.

### 7. Adapters de persistencia (Postgres, real, local)

- `AdministratorRepository.FindByID`/`FindByEmail`: devuelven los datos guardados por `Create` (incluye `password_hash`, `encrypted_totp_secret`, `last_accepted_totp_period`); devuelven `(nil, nil)` para un id/email inexistente.
- `AdministratorRepository.UpdateLastAcceptedTOTPPeriod`: persiste el nuevo valor; una lectura posterior con `FindByEmail` lo refleja.
- `AdministratorSessionRepository.FindByTokenHash`: devuelve la sesión guardada por `Save`; `(nil, nil)` para un hash inexistente.
- `AdministratorSessionRepository.UpdateIdleExpiry`: persiste el nuevo `idle_expires_at`.
- `AdministratorSessionRepository.Revoke`: persiste `revoked_at`; una sesión revocada deja de aparecer como activa vía `FindByTokenHash` + `IsActive`.
- Migración `20260822_05`: la columna `last_accepted_totp_period` existe en `administrators` tras migrar, con default `0`; no toca ni redefine `20260822_01`-`04`; test de idempotencia (`Migrate` dos veces no falla), mismo patrón que `TestMigration03IsIdempotentWhenColumnAlreadyExists` de PB-062.

### 8. Middleware — `RequireSession`

- Sin cookie → `401`, no se llama al handler siguiente.
- Cookie con valor que no matchea ninguna sesión (hash no encontrado) → `401`.
- Cookie de una sesión expirada o revocada → `401`.
- Cookie de una sesión activa → se llama al handler siguiente con `domain.AdministratorSession` disponible vía `SessionFromContext`.
- Ningún caso de error expone el valor de la cookie ni detalles internos en la respuesta.

### 9. REST — DTOs, handlers y router

- `LoginRequest` deserializa `email`/`password`/`totp_code`; validación de shape igual que `SetupRequest` (email formato+254, password 12-128, totp_code `^[0-9]{6}$`) → `400` sin invocar el usecase si falla.
- Login exitoso → `200`, `Set-Cookie` con los mismos atributos que PB-062, `Cache-Control: no-store`, body `SessionResponse`.
- Login con `ErrInvalidCredentials` → `401` (a diferencia de `/auth/setup*`, que usa `400`/`409` — este endpoint usa el mapeo estándar de sufijo `U`→401 que ya provee `response.WriteDomainError`, sin código nuevo en la capa REST).
- Login rate limited → `429` + `Retry-After`.
- `GET /auth/session` sin sesión válida → `401` (vía middleware, antes de llegar al handler). Con sesión válida → `200` con la proyección correcta, y un segundo `GET` inmediato muestra `idle_expires_at` estrictamente mayor (test de integración con reloj real o fake inyectado en el usecase, según convenga).
- `DELETE /auth/session` sin sesión válida → `401`. Con sesión válida → `204`, header `Set-Cookie` que expira la cookie (valor vacío, `Max-Age`/`Expires` en el pasado); un `GET /auth/session` posterior con la misma cookie original → `401`.
- Router: las tres rutas responden a través del router real; `GET`/`DELETE /auth/session` pasan por el middleware.

### 10. Validaciones finales

- `go test ./...`, `go test -race ./...`, `go vet ./...`, `gofmt -l .` sin diferencias.
- `.harness/kernel/scripts/check-backend-architecture.sh`, `check-openapi.sh` (sin modificar `docs/openapi.yaml`), `check-security.sh`.
- Verificación manual end-to-end contra Postgres local: activar un Administrator (flujo PB-061→PB-062), generar el código TOTP real vigente, `POST /auth/login` exitoso, reintentar el mismo código → rechazado por reutilización, `GET /auth/session` dos veces confirmando que `idle_expires_at` avanza, `DELETE /auth/session`, `GET /auth/session` posterior con la cookie ya revocada → `401`. Inspección de logs y respuestas sin password/TOTP/token en claro.

## Expected failing tests before implementation

- No existe `ExtendIdle` ni `ErrInvalidCredentials`.
- `AdministratorRepository`/`AdministratorSessionRepository`/`PasswordHasher`/`SessionTokenGenerator` no tienen los métodos nuevos; `TOTPVerifier.Verify` no devuelve `period`.
- No existen los in ports ni los usecases nuevos.
- No existe `model.Administrator.LastAcceptedTOTPPeriod` ni la migración `20260822_05`.
- No existen los repositorios/adapters nuevos ni el middleware.
- No existen el DTO `LoginRequest`, los handlers `Login`/`GetCurrentSession`/`Logout`, ni las rutas correspondientes.

Los tests se escribirán después de la aprobación de este plan y deberán fallar (o no compilar) antes de incorporar la implementación correspondiente.

## Acceptance criteria covered

- Los tres endpoints devuelven exactamente los schemas y códigos de estado del OpenAPI.
- El idle TTL se desliza en cada request autenticada, capado al absoluto.
- Se rechaza la reutilización del último período TOTP aceptado.
- El secreto TOTP y el password hash nunca se exponen.
- `internal/core/**` no importa GORM, Chi, `net/http` ni SDKs externos.
- Los tres scripts del harness pasan sin modificar `docs/openapi.yaml`.
- Las migraciones `20260822_01`-`04` quedan intactas.

## Open questions / human approval notes — resueltas

- Idle TTL deslizante vs. fijo: **desliza en cada request autenticada exitosa**, confirmado explícitamente por el usuario antes de este plan.
- Cambio de firma de `TOTPVerifier.Verify` (agrega `period`): necesario para cumplir el requisito de esta tarea; toca el call site de PB-062 de forma puramente mecánica (ignora el valor nuevo), documentado en `implementation-brief.md`.
- Validación de `Origin` en `DELETE /auth/session`: fuera de alcance por falta de configuración previa (`AKRITAS_PUBLIC_URL`), documentado como gap explícito en `implementation-brief.md`, no decidido en silencio.

No quedan decisiones abiertas. Se requiere aprobación humana explícita de este archivo antes de crear tests o implementar código.
