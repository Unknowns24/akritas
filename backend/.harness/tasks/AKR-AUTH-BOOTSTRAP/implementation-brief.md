# Implementation Brief

## Task

`AKR-AUTH-BOOTSTRAP`: implementar `GET /auth/setup-status` y `POST /auth/setup`, el arranque del flujo de bootstrap del único Administrator, sin persistir todavía la cuenta.

## Current project context

`internal/core/domain` ya define `Administrator` (id, email, display_name, created_at, updated_at — sin password ni TOTP) y `AdministratorSession`, con sus errores `0x401001V`–`0x401004C` en `errors.go`. Todo lo demás bajo `internal/` es placeholder: no hay ports, usecases, service layer, adapters de DB, adapters externos, DTOs, handlers, middleware ni router. `cmd/main.go` es `package main; func main() {}`. `go.mod` solo tiene `github.com/google/uuid`; no hay ORM, router HTTP ni librerías de hashing/TOTP. No existe ningún archivo de migración en el repo — esta tarea agrega la primera.

Fuentes autoritativas revisadas:

- `docs/openapi.yaml`: paths `/auth/setup-status` (`getAuthSetupStatus`) y `/auth/setup` (`startAdministratorSetup`); schemas `SetupStatus`, `SetupStatusResponse`, `SetupRequest`, `TotpEnrollment`, `TotpEnrollmentResponse`, `Error`/`ErrorResponse`, y las `responses` reutilizables `BadRequest`/`Conflict`/`TooManyRequests`/`InternalError`.
- ADR-008 (single-admin, bootstrap token, TOTP, Argon2id) y ADR-005 (Credential Store cifrado con `AKRITAS_MASTER_KEY`).
- `docs/configuration.md`: contrato de `AKRITAS_BOOTSTRAP_TOKEN` y `AKRITAS_MASTER_KEY` (alta entropía, nunca persistidos/logueados, comparación constant-time, rate limit en el borde HTTP).
- Policies del profile `backend_api`: `project-structure.md`, `wiring.md`, `modularity-srp.md`, `domain-errors.md`, `migrations.md`, `external-adapters.md`, `uker.md`, `openapi.md`, `testing.md`, `security.md`, `architecture-decisions.md`.
- `internal/core/domain/administrator.go`, `administrator_session.go`, `errors.go`, `validation.go` como referencia de estilo (constructor + `Validate()`, helpers `nonBlank`/`validTime`, sentinels `newDomainError`).

Decisiones tomadas explícitamente con el usuario, por no estar documentadas en ningún lado del repo (regla de `architecture-decisions.md`: no decidir en silencio):

- Motor de base de datos: **PostgreSQL**.
- Variable de entorno del connection string: **`AKRITAS_DB_DSN`**.

## Proposed approach

1. Dominio: agregar `PendingEnrollment` (id, email, display_name, password_hash, encrypted_totp_secret, created_at, expires_at) con constructor + `Validate()` siguiendo el estilo de `administrator_session.go`, y tres sentinels nuevos en `errors.go` continuando la secuencia del grupo `01` (auth): `ErrInvalidBootstrapToken` (`0x401005V`), `ErrInvalidPendingEnrollment` (`0x401006V`), `ErrAdministratorAlreadyExists` (`0x401007C`).
2. Ports: `internal/core/ports/in/` con `GetSetupStatusUseCase` y `StartAdministratorSetupUseCase`; `internal/core/ports/out/` con ocho ports de un archivo cada uno: `AdministratorRepository` (solo `ExistsActive`), `PendingEnrollmentRepository` (`Save`), `CredentialStore` (`Encrypt`), `TOTPSecretGenerator`, `PasswordHasher` (`Hash`), `BootstrapTokenVerifier`, `RateLimiter`, `Clock`.
3. Usecases en `internal/usecase/auth/`: `GetSetupStatus` deriva `setup_required`/`registration_open` de `!AdministratorRepository.ExistsActive`. `StartAdministratorSetup` encadena: rate limit → verificación constant-time del bootstrap token → chequeo de Administrator existente (409) → hash Argon2id de la password → generación del secreto TOTP → cifrado del secreto vía Credential Store → construcción y persistencia del `PendingEnrollment` → respuesta con el material en claro (única vez).
4. Adapters de persistencia en `internal/adapter/db/postgres/`: modelos GORM, dos migraciones gormigrate (`administrators`, `pending_enrollments`) con `Rollback`, y los dos repositorios (`ExistsActive`, `Save`).
5. Adapters de seguridad en `internal/adapter/security/` (directorio nuevo, mismo nivel que `adapter/db`/`adapter/external`): Argon2id hasher, generador TOTP (`pquerna/otp/totp`), Credential Store AES-256-GCM sobre `AKRITAS_MASTER_KEY` (stdlib, sin dependencia nueva), verificador de bootstrap token (`crypto/subtle`), rate limiter en memoria, reloj.
6. REST: DTOs en `internal/adapter/rest/dto/auth/` espejando los schemas del OpenAPI; handlers en `internal/adapter/rest/handler/auth/` (uno por operación); registro de rutas bajo `/api/v1/auth/` en `internal/adapter/rest/router/` con `go-chi/chi/v5`.
7. Wiring en `cmd/main.go` y `config/` siguiendo el orden de `wiring.md`: config → DB + migraciones → repositorios → adapters → usecases → router.
8. Documentar `AKRITAS_DB_DSN` en `docs/configuration.md`.

## Architecture impact

Primeras implementaciones concretas de las capas hexagonales que hasta ahora eran placeholders:

```text
internal/core/ports/{in,out}/          (8 out + 2 in)
internal/usecase/auth/
internal/adapter/db/postgres/{migrations,model,repository}/
internal/adapter/security/             (nuevo, no listado en project-structure.md — se agrupa junto a db/external por ser infraestructura, no core)
internal/adapter/rest/{dto,handler,router}/auth/
config/
```

`internal/core/**` sigue sin depender de GORM, Chi, `net/http` ni SDKs externos — solo dominio + ports. `internal/usecase/**` depende únicamente de `internal/core/ports/out`, nunca de adapters concretos (`wiring.md`, `external-adapters.md`).

## API/OpenAPI impact

No se modifica `docs/openapi.yaml`: el contrato de `/auth/setup-status` y `/auth/setup` ya está completo desde `AKR-OPENAPI-MVP`. Se replica tal cual, incluyendo que `startAdministratorSetup` declara su propio header `Cache-Control: no-store` inline en vez del `$ref: "#/components/headers/NoStore"` que usan los endpoints hermanos — es una inconsistencia real y preexistente del spec, no se "corrige" sin que se pida explícitamente.

## Data/persistence impact

Primera migración del proyecto. Dos tablas nuevas vía gormigrate, IDs `20260822_01_create_administrators` y `20260822_02_create_pending_enrollments`, cada una con `Rollback` (drop table). `administrators` se crea pero esta tarea solo la lee (`ExistsActive`); el primer `INSERT` real es de PB-062. `pending_enrollments` guarda `password_hash` (ya hasheado) y `encrypted_totp_secret` (ciphertext) — nunca texto plano. Motor: PostgreSQL vía `AKRITAS_DB_DSN` (decisión explícita del usuario).

## Error handling impact

Tres sentinels nuevos en el grupo `01` (auth) de `errors.go`, mismo patrón `0x4XXNNNT` que los cuatro existentes:

- `ErrInvalidBootstrapToken` (`0x401005V`) → 400.
- `ErrInvalidPendingEnrollment` (`0x401006V`) → 400.
- `ErrAdministratorAlreadyExists` (`0x401007C`) → 409.

El rate limit (429, con `Retry-After`) no pasa por el catálogo de errores de dominio: el patrón `ErrorCode` del OpenAPI (`^[0-9A-F]x[0-9A-F]{6}[VUFNCI]$`) no tiene letra para 429, así que se resuelve directamente en usecase/handler.

## Test strategy

- Dominio: `pending_enrollment_test.go` (estilo `administrator_session_test.go`, `testing` estándar sin testify, `t.Parallel()`); sumar los 3 códigos nuevos al test de catálogo de `errors_test.go`.
- Usecases: tests de `GetSetupStatus` y `StartAdministratorSetup` con fakes/stubs de los 8 out ports, sin base de datos real — cubren el flujo feliz y cada rama de error (rate limited, token inválido, administrator existente).
- Handlers: tests de mapeo transporte↔dominio, parsing de body y códigos de estado, con las dependencias `in` mockeadas.
- Repositorios: solo si el comportamiento de persistencia no es trivial (upsert de `PendingEnrollment`) — se prueba contra una base real o vía `sqlmock`, a definir en `tdd-test-plan.md`.
- Validaciones finales: `go test ./...`, `go vet ./...`, `gofmt -l .`, y los tres scripts del harness.

## Risks

- Elegir PostgreSQL agrega infraestructura externa a un entorno que hoy no la requiere (Go compila sin DB); se documenta como decisión explícita y no silenciosa.
- Sumar `pquerna/otp` como dependencia nueva solo para *generar* el secreto (no verificarlo en esta tarea) podría parecer prematuro, pero evita reimplementar a mano el formato del otpauth URI y ya hará falta en PB-062 para verificar.
- El tiempo de expiración del pending enrollment ("corta duración" según ADR-008, sin número fijo) queda como constante propuesta (10 minutos) en el usecase — ajustable en revisión si no es lo esperado.
- Un rate limiter en memoria no sobrevive reinicios ni escala a múltiples instancias; aceptable porque el alcance excluye explícitamente rate limiting avanzado (PB-065).
- El upsert de `PendingEnrollment` (reemplazar cualquier pendiente anterior) es una interpretación razonable pero no está escrita literalmente en ADR-008; se deja explícita para que la revisión la objete si no es la esperada.

## Files likely to change

- `internal/core/domain/pending_enrollment.go`, `pending_enrollment_test.go`, `errors.go`, `errors_test.go`.
- `internal/core/ports/in/get_setup_status.go`, `start_administrator_setup.go`.
- `internal/core/ports/out/{administrator_repository,pending_enrollment_repository,credential_store,totp_secret_generator,password_hasher,bootstrap_token_verifier,rate_limiter,clock}.go`.
- `internal/usecase/auth/*.go` y sus tests.
- `internal/adapter/db/postgres/**`.
- `internal/adapter/security/**`.
- `internal/adapter/rest/{dto,handler,router}/auth/**`.
- `config/config.go`, `cmd/main.go`.
- `go.mod`, `go.sum`.
- `docs/configuration.md`.
- Artefactos de `.harness/tasks/AKR-AUTH-BOOTSTRAP/`.
