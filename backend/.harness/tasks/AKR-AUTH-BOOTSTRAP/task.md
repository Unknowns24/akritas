# AKR-AUTH-BOOTSTRAP - Bootstrap del único Administrator

## Estado

complete

## Tipo de tarea

backend-api-feature

## Modo de proyecto

existing_project

## Contexto

`AKR-BACKEND-FOUNDATION` dejó el dominio `Administrator`/`AdministratorSession` modelado en `internal/core/domain`, pero ports, usecases, adapters, migraciones y HTTP quedaron intencionalmente vacíos (placeholders `.gitkeep`). `AKR-OPENAPI-MVP` ya fijó el contrato completo del flujo de autenticación en `docs/openapi.yaml`, incluyendo `GET /auth/setup-status` y `POST /auth/setup`, y ADR-008/ADR-005 documentan el flujo de bootstrap, TOTP y Credential Store. `go.mod` solo depende de `github.com/google/uuid` y `cmd/main.go` es un entrypoint inerte: esta tarea es la primera vez que el proyecto persiste algo.

## Objetivo

Implementar `GET /auth/setup-status` (`getAuthSetupStatus`) y `POST /auth/setup` (`startAdministratorSetup`) exactamente como especifica `docs/openapi.yaml`: exponer si el registro inicial sigue abierto, y validar el bootstrap token para crear un pending enrollment de corta duración con un secreto TOTP nuevo, sin persistir todavía al `Administrator`.

## Requerimiento funcional

- `GET /auth/setup-status` devuelve `SetupStatusResponse { data: { setup_required, registration_open } }` sin revelar la identidad del administrador.
- `POST /auth/setup` valida `SetupRequest { email, display_name, password, bootstrap_token }`.
- El bootstrap token se compara en tiempo constante contra `AKRITAS_BOOTSTRAP_TOKEN` y el endpoint está rate limited.
- Se genera un secreto TOTP nuevo e independiente del bootstrap token y de `AKRITAS_MASTER_KEY`; se cifra en reposo vía Credential Store antes de persistir el pending enrollment.
- La respuesta `201 TotpEnrollmentResponse { data: { enrollment_id, otpauth_uri, manual_entry_key, expires_at } }` expone el material de aprovisionamiento en claro exactamente una vez, con `Cache-Control: no-store`.
- Si ya existe un Administrator activo, la solicitud responde `409`.
- Si el bootstrap token es inválido, la solicitud responde `400`.
- Si se supera el límite de intentos, la solicitud responde `429` con `Retry-After`.
- No se persiste el `Administrator` en esta tarea.

## Criterios de aceptación

- `go test ./...`, `go vet ./...` y `gofmt -l .` (sin diffs) finalizan correctamente.
- Los dos endpoints devuelven exactamente los schemas y códigos de estado definidos en `docs/openapi.yaml` (`SetupStatusResponse`, `TotpEnrollmentResponse`, `ErrorResponse` para 400/409/429/500).
- El secreto TOTP nunca se persiste en claro; solo el ciphertext del Credential Store llega a la base de datos.
- La comparación del bootstrap token es constant-time y nunca se loguea ni se devuelve por API.
- `internal/core/**` no importa GORM, Chi, `net/http` ni SDKs externos.
- Los errores nuevos cumplen el patrón `0x4XXNNNT` ya usado en `errors.go` y están documentados.
- Los checks de arquitectura, OpenAPI y seguridad del harness pasan sin modificar `docs/openapi.yaml`.

## Restricciones técnicas

- Profile: `.harness/kernel/profiles/go-hexagonal-api.yaml`.
- Workflow: `.harness/kernel/workflows/backend-api-feature.yaml`.
- Motor de base de datos: PostgreSQL, vía `gorm.io/gorm` + `gorm.io/driver/postgres` + `github.com/go-gormigrate/gormigrate/v2`. Decisión tomada explícitamente con el usuario por no estar documentada en el repo.
- Nueva variable de entorno `AKRITAS_DB_DSN` (connection string), documentada en `docs/configuration.md` como parte de esta tarea.
- Router HTTP: `github.com/go-chi/chi/v5` (implícito por exclusión ya declarada en `project-structure.md`).
- Hashing de password: Argon2id (`golang.org/x/crypto/argon2`) con `m=19456,t=2,p=1` según ADR-008.
- Generación TOTP: `github.com/pquerna/otp/totp`, RFC 6238, 6 dígitos, período 30s.
- Migraciones versionadas `YYYYMMDD_NN_<descripcion>`, cada una con `Rollback` explícito; ninguna usa `AutoMigrate` global.
- No implementar código antes de la aprobación humana de `tdd-test-plan.md`.

## Archivos o zonas probablemente afectadas

- `internal/core/domain/pending_enrollment.go`, `errors.go`.
- `internal/core/ports/in/`, `internal/core/ports/out/`.
- `internal/usecase/auth/`.
- `internal/adapter/db/postgres/` (modelos, migraciones, repositorios).
- `internal/adapter/security/` (password hasher, TOTP generator, credential store, bootstrap token verifier, rate limiter, clock).
- `internal/adapter/rest/dto/auth/`, `internal/adapter/rest/handler/auth/`, `internal/adapter/rest/router/`.
- `config/`, `cmd/main.go`.
- `docs/configuration.md` (nueva variable `AKRITAS_DB_DSN`).
- `go.mod`, `go.sum`.
- `.harness/tasks/AKR-AUTH-BOOTSTRAP/` y `.harness/tasks/index.md`.

## Fuera de alcance

- Verificación del código TOTP y activación del Administrator (`POST /auth/setup/verify`, PB-062).
- Login (`POST /auth/login`, PB-063).
- Recovery (`POST /auth/recovery*`, PB-064).
- Rate limiting avanzado más allá de un limitador simple en memoria (PB-065).
- Persistir el `Administrator` (primer `INSERT` real ocurre en PB-062).
- Sesiones (`AdministratorSession` ya modelada, pero no se emite ninguna sesión en esta tarea).

## Instrucción para el harness

Primero generar `implementation-brief.md` y `tdd-test-plan.md`. No implementar código hasta aprobación humana.
