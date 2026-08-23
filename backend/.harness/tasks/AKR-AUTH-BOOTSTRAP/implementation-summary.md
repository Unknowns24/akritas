# Implementation Summary

## Implemented

- `PendingEnrollment` en `internal/core/domain`, con constructor + `Validate()` y `IsExpired`, siguiendo el estilo de `administrator_session.go`.
- Tres sentinels nuevos en `errors.go` (`ErrInvalidBootstrapToken` `0x401005V`, `ErrInvalidPendingEnrollment` `0x401006V`, `ErrAdministratorAlreadyExists` `0x401007C`), registrados en `DomainErrors()` y en `docs/errors/aaa-map.md`.
- Dos in ports (`GetSetupStatusUseCase`, `StartAdministratorSetupUseCase`) y ocho out ports (repos de Administrator/PendingEnrollment, CredentialStore, TOTPSecretGenerator, PasswordHasher, BootstrapTokenVerifier, RateLimiter, Clock).
- Usecases `GetSetupStatus` y `StartAdministratorSetup` en `internal/usecase/auth`, con el orden de verificación rate limit → bootstrap token → administrator existente, y `ErrSetupRateLimited` como sentinel de transporte fuera del catálogo de dominio (el pattern `ErrorCode` no tiene letra para 429).
- Primera migración Postgres del proyecto (`20260822_01_create_administrators`, `20260822_02_create_pending_enrollments`, vía GORM + gormigrate, con `Rollback`), y los repositorios `AdministratorRepository.ExistsActive` / `PendingEnrollmentRepository.Save` (upsert: reemplaza cualquier pending enrollment anterior).
- Adapters de seguridad: Argon2id (`m=19456,t=2,p=1`), generador TOTP (`pquerna/otp`, RFC 6238 defaults), Credential Store AES-256-GCM sobre `AKRITAS_MASTER_KEY` (clave derivada vía SHA-256, stdlib únicamente), verificador de bootstrap token constant-time, rate limiter en memoria (5 intentos / 15 min por IP, confirmado por el usuario), reloj UTC.
- DTOs REST espejando `SetupStatus(Response)`, `SetupRequest`, `TotpEnrollment(Response)` del OpenAPI; envelope de error compartido (`dto.Error`/`ErrorResponse`) con dos códigos de capa REST (`0x100001I` interno, `0x100002V` request malformado, `0x100003C` rate limit — documentados en `aaa-map.md`, fuera del catálogo de dominio).
- Handlers `GetSetupStatus`/`StartAdministratorSetup` y router chi bajo `/api/v1/auth/{setup-status,setup}`, wiring completo en `cmd/main.go` (config → DB+migraciones → repos → adapters → usecases → router).
- `AKRITAS_DB_DSN` documentada en `docs/configuration.md`.

## Deliberately not implemented

- Verificación del código TOTP y activación del Administrator (`POST /auth/setup/verify`) — PB-062.
- Login, sesión, logout — PB-063.
- Recovery — PB-064.
- Rate limiting avanzado (persistente/distribuido) — PB-065.
- Cualquier `INSERT` en `administrators`: esta tarea sólo lee (`ExistsActive`).

## Validation result

- `go test ./...`: pass.
- `go test -race ./...`: pass.
- Coverage: domain 82.2%, usecase/auth 91.2%, adapter/security 87.2%, adapter/db/postgres/migrations 100%, repository/administrator 80.0%, repository/pending_enrollment 85.7%, adapter/rest/handler/auth 97.7%, adapter/rest/response 86.7%, adapter/rest/router 100%, config 92.9%.
- `go vet ./...`: pass.
- `gofmt -l .`: sin diferencias.
- `check-backend-architecture.sh`, `check-openapi.sh` (59 operaciones, 112 schemas, sin cambios al spec), `check-security.sh`: pass.
- Verificación manual end-to-end contra Postgres local: `setup-status` (200), `setup` exitoso (201 + `Cache-Control: no-store`), token incorrecto (400), administrator existente (409), rate limit tras 5 intentos (429 + `Retry-After: 60`); logs y respuestas sin bootstrap token, password ni secreto TOTP en texto plano.
