# TDD Test Plan

## Scope

Definir mediante tests el contrato de `GET /auth/setup-status` y `POST /auth/setup`: dominio (`PendingEnrollment`, errores nuevos), los ocho out ports y dos in ports, los dos usecases de `internal/usecase/auth`, los dos repositorios Postgres, los cinco adapters de seguridad, y los dos handlers REST con su wiring en el router. No se prueba verificación TOTP, login, recovery ni persistencia real del `Administrator` porque no forman parte de esta tarea.

## Tests to add/update

### 1. Dominio — `PendingEnrollment`

- `NewPendingEnrollment` acepta id/email/display_name/password_hash/encrypted_totp_secret/created_at/expires_at válidos y devuelve una entidad cuyos campos coinciden exactamente con la entrada.
- Rechaza `id` cero (`uuid.Nil`), `email` vacío o sin `@`, `display_name` vacío, `password_hash` vacío, `encrypted_totp_secret` vacío, `created_at` zero-value, y `expires_at` que no sea estrictamente posterior a `created_at` — todos con `ErrInvalidPendingEnrollment`.
- `IsExpired(now)` es `false` estrictamente antes de `expires_at`, `true` en `expires_at` y después (borde inclusivo hacia expirado, simétrico a `AdministratorSession.IsActive`).
- El package `domain` sigue sin tags `json`/`gorm` y sin importar HTTP, Chi, GORM ni adapters concretos.

### 2. Dominio — catálogo de errores

- `ErrInvalidBootstrapToken` = `0x401005V`, `ErrInvalidPendingEnrollment` = `0x401006V`, `ErrAdministratorAlreadyExists` = `0x401007C`; los tres cumplen el pattern `^[0-9A-F]x[0-9A-F]{6}[VUFNCI]$` y están registrados en `DomainErrors()`.
- Los tres códigos son únicos en el catálogo completo (ningún choque con los 40 sentinels existentes).
- `errors.Is` reconoce cada sentinel por código; `Wrap` conserva `errors.Is`/`errors.As` vía `Unwrap`; `Error()` no expone la causa envuelta.

### 3. Ports — compilación de contratos

- Cada out port (`AdministratorRepository`, `PendingEnrollmentRepository`, `CredentialStore`, `TOTPSecretGenerator`, `PasswordHasher`, `BootstrapTokenVerifier`, `RateLimiter`, `Clock`) es una interfaz en `internal/core/ports/out`, un archivo por port, sin importar GORM/Chi/SDKs.
- Cada in port (`GetSetupStatusUseCase`, `StartAdministratorSetupUseCase`) es una interfaz en `internal/core/ports/in` que no depende de `internal/adapter/**`.
- No hay test de comportamiento aquí (son interfaces); se verifica indirectamente por compilación al implementar fakes en los tests de usecase.

### 4. Usecase — `GetSetupStatus`

- Con `AdministratorRepository.ExistsActive` devolviendo `false`: `Execute` devuelve `{SetupRequired: true, RegistrationOpen: true}`.
- Con `ExistsActive` devolviendo `true`: `Execute` devuelve `{SetupRequired: false, RegistrationOpen: false}`.
- Si `ExistsActive` devuelve error, `Execute` propaga el error sin envolverlo en un sentinel de dominio (es un fallo de infraestructura, no una regla de negocio).

### 5. Usecase — `StartAdministratorSetup`

Con fakes/stubs de los 8 out ports:

- Camino feliz: `RateLimiter.Allow` = true, `BootstrapTokenVerifier.Verify` = true, `AdministratorRepository.ExistsActive` = false → se llama `PasswordHasher.Hash` con la password recibida, `TOTPSecretGenerator.Generate` con issuer `"Akritas"` y el email, `CredentialStore.Encrypt` con el `Base32Key` generado, `PendingEnrollmentRepository.Save` con un `PendingEnrollment` cuyo `ExpiresAt` es `Clock.Now() + 10min`; el `Output` devuelto expone `OtpauthURI`/`ManualEntryKey` en claro (los del generador, no el ciphertext) y el mismo `ExpiresAt` persistido.
- `RateLimiter.Allow` = false → `Execute` devuelve un error de rate limit distinguible (para que el handler responda 429 + `Retry-After`); no se llama a ningún otro port.
- `BootstrapTokenVerifier.Verify` = false → `Execute` devuelve `ErrInvalidBootstrapToken`; no se llama `PasswordHasher`, `TOTPSecretGenerator`, `CredentialStore` ni `PendingEnrollmentRepository.Save`.
- `AdministratorRepository.ExistsActive` = true → `Execute` devuelve `ErrAdministratorAlreadyExists`; no se llama `PasswordHasher`, `TOTPSecretGenerator`, `CredentialStore` ni `Save` (se verifica antes de generar cualquier secreto).
- Orden de verificación es exactamente rate limit → bootstrap token → administrator existente (se prueba con un fake que registra la secuencia de llamadas).
- Si `PasswordHasher.Hash`, `TOTPSecretGenerator.Generate`, `CredentialStore.Encrypt` o `PendingEnrollmentRepository.Save` devuelven error, `Execute` propaga el error y no persiste un `PendingEnrollment` parcial.
- El secreto TOTP en claro (`Base32Key`/`OtpauthURI`) nunca se pasa a `PendingEnrollmentRepository.Save` — solo el resultado de `CredentialStore.Encrypt` llega al `PendingEnrollment`.

### 6. Adapters de seguridad

- `PasswordHasher.Hash`: dos hashes de la misma password son distintos (salt aleatorio); el string codificado contiene el prefijo `$argon2id$` y los parámetros `m=19456,t=2,p=1`.
- `TOTPSecretGenerator.Generate("Akritas", email)`: `Base32Key` matchea `^[A-Z2-7]{16,128}$` (mismo pattern que `TotpEnrollment.manual_entry_key` en el OpenAPI); `OtpauthURI` empieza con `otpauth://totp/` e incluye el email y el secreto codificado.
- `CredentialStore.Encrypt`: el mismo plaintext produce ciphertexts distintos en llamadas sucesivas (nonce aleatorio); el ciphertext nunca contiene el plaintext como substring.
- `BootstrapTokenVerifier.Verify`: `true` solo cuando el candidato coincide byte a byte con el token configurado; se documenta (no se puede testear en tiempo real de forma determinística) que la comparación usa `crypto/subtle.ConstantTimeCompare`.
- `RateLimiter.Allow`: permite hasta el límite configurado para una key dentro de la ventana y lo deniega al superarlo; una key distinta no se ve afectada por el límite de otra.
- `Clock.Now()`: devuelve un `time.Time` en UTC no-zero.

### 7. Adapters de persistencia (Postgres)

- Migraciones: `20260822_01_create_administrators` y `20260822_02_create_pending_enrollments` se registran en `migrate.go` en ese orden, cada una con `ID`, `Migrate` y `Rollback` no vacíos.
- `AdministratorRepository.ExistsActive`: devuelve `false` contra una tabla vacía y `true` tras insertar una fila directamente (vía el modelo GORM, no vía este repo, ya que esta tarea no persiste Administrator).
- `PendingEnrollmentRepository.Save`: persiste un `PendingEnrollment` nuevo; una segunda llamada reemplaza (upsert) el enrollment pendiente anterior en vez de acumular filas — se verifica que solo queda una fila vigente tras dos `Save` sucesivos.
- Estos tests requieren una base Postgres real (o `sqlmock` si se decide evitar la dependencia de infraestructura en CI — a confirmar en la aprobación de este plan); no se usa `AutoMigrate` global en ningún test.

### 8. REST — DTOs, handlers y router

- `SetupRequest` deserializa exactamente los campos del OpenAPI (`email`, `display_name`, `password`, `bootstrap_token`) y rechaza JSON con campos desconocidos si el contrato lo exige (`additionalProperties: false`).
- `getSetupStatus` handler: responde `200` con el body `SetupStatusResponse` mapeado 1:1 desde el usecase.
- `startAdministratorSetup` handler:
  - Body válido + usecase feliz → `201`, header `Cache-Control: no-store`, body `TotpEnrollmentResponse` con los cuatro campos.
  - Body malformado (campo faltante, email inválido, password fuera de 12–128, bootstrap_token fuera de 32–512) → `400` `ErrorResponse` sin llamar al usecase.
  - Usecase devuelve `ErrInvalidBootstrapToken` → `400` `ErrorResponse`.
  - Usecase devuelve `ErrAdministratorAlreadyExists` → `409` `ErrorResponse`.
  - Usecase devuelve error de rate limit → `429` `ErrorResponse` con header `Retry-After`.
  - Ningún error de estos tres devuelve detalle que permita distinguir cuál de email/password/bootstrap_token fue el problema (mensaje genérico), consistente con ADR-008.
  - La respuesta nunca incluye `bootstrap_token` ni la password recibida, ni en el body ni en logs (se verifica con un logger de test que capture y falle si aparecen).
- Router: `GET /api/v1/auth/setup-status` y `POST /api/v1/auth/setup` responden a través del router real (test de integración liviano con `httptest`), sin autenticación previa (`security: []` en el OpenAPI).

### 9. Wiring y validaciones finales

- `cmd/main.go` arranca sin panics con las tres variables de entorno requeridas (`AKRITAS_BOOTSTRAP_TOKEN`, `AKRITAS_MASTER_KEY`, `AKRITAS_DB_DSN`) seteadas a valores de prueba, y falla rápido (log + exit) si falta alguna, sin imprimir su valor.
- Ejecutar `go test ./...`.
- Ejecutar `go vet ./...`.
- Ejecutar `gofmt` sobre los archivos Go y comprobar que no queden diferencias.
- Ejecutar `.harness/kernel/scripts/check-backend-architecture.sh`.
- Ejecutar `.harness/kernel/scripts/check-openapi.sh` sin modificar `docs/openapi.yaml`.
- Ejecutar `.harness/kernel/scripts/check-security.sh`.

## Expected failing tests before implementation

- No existe `internal/core/domain/pending_enrollment.go` ni sus errores en `errors.go`.
- No existen los ports `in`/`out` de esta tarea.
- No existe `internal/usecase/auth`.
- No existen los adapters de `internal/adapter/db/postgres` ni `internal/adapter/security`.
- No existen DTOs, handlers ni rutas de `/auth/setup-status` y `/auth/setup` bajo `internal/adapter/rest`.
- No existe `config/config.go` ni wiring en `cmd/main.go`.
- No existe la variable `AKRITAS_DB_DSN` en `docs/configuration.md`.

Los tests se escribirán después de la aprobación de este plan y deberán fallar (o no compilar) antes de incorporar la implementación correspondiente.

## Acceptance criteria covered

- Los dos endpoints devuelven exactamente los schemas y códigos de estado del OpenAPI.
- El secreto TOTP nunca se persiste en claro.
- La comparación del bootstrap token es constant-time y nunca se loguea ni se devuelve por API.
- `internal/core/**` no importa GORM, Chi, `net/http` ni SDKs externos.
- Los errores nuevos cumplen `0x4XXNNNT` y están documentados.
- Los tres scripts del harness pasan sin modificar `docs/openapi.yaml`.

## Open questions / human approval notes — resueltas

Aprobado por el usuario el 2026-08-22 con los siguientes ajustes:

1. Tests de repositorio: **Postgres real (local)**, no `sqlmock`.
2. Rate limiter: **5 intentos cada 15 minutos por IP** para `/auth/setup` (reemplaza el placeholder sin número de la sección 6).
3. Upsert del pending enrollment: **confirmado** tal como se propuso — un nuevo `POST /auth/setup` reemplaza cualquier enrollment pendiente anterior.

No quedan decisiones abiertas. Se procede a tests + implementación.
