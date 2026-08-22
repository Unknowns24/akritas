# Task

## ID

`AKR-AUTH-INTEGRATION`

## Objetivo

Integrar mediante un merge real sin commit la funcionalidad existente de
`origin/feat/authentication` sobre `feat/backend-milestone-1`, conservando
setup, verificación TOTP, login, consulta/revocación de sesión y protección de
las integraciones, pero reemplazando las decisiones de infraestructura que
contradicen los ADRs y la estructura vigente.

## Profile y workflow

- Profile: `backend_api` (`.harness/kernel/profiles/go-hexagonal-api.yaml`).
- Workflow: `.harness/kernel/workflows/backend-api-feature.yaml`.
- Estado: `complete` (plan TDD aprobado y validaciones completadas el 2026-08-22).

## Alcance

- Merge `origin/feat/authentication` con `--no-commit` y resolución semántica de
  conflictos.
- `GET /auth/setup-status`.
- `POST /auth/setup` y `POST /auth/setup/verify`.
- `POST /auth/login`.
- `GET /auth/session` y `DELETE /auth/session`.
- Middleware de sesión y validación de Origin para mutaciones autenticadas.
- Configuración centralizada con Viper.
- Persistencia PostgreSQL de autenticación mediante entidades de dominio con
  tags GORM pasivos.
- Reutilización del Credential Store PostgreSQL y cipher vigentes.
- Transacciones application-level documentadas mediante ADR.
- Tipo estable de error `R` para HTTP 429.
- Wiring compartido de auth e integraciones.

## Fuera de alcance

- Recovery (`POST /auth/recovery*`).
- RBAC, múltiples administradores, SSO, passkeys o recovery codes.
- Rate limiting distribuido o persistente.
- Compatibilidad con bases no descartables que hayan aplicado las migraciones
  experimentales de `feat/authentication`.
- Creación del merge commit.

## Decisiones humanas confirmadas

- La integración debe quedar como merge real sin commit.
- Se conserva únicamente la funcionalidad ya implementada en la rama auth;
  recovery queda para una tarea separada.
- Se conserva `out.Transactor` y se documenta en un ADR.
- La tabla del Credential Store pasa de `integration_credentials` a
  `credentials`; la base actual puede recrearse.
- Las migraciones de auth sólo se aplicaron en entornos locales/descartables.
- El formato de errores incorpora `R` para rate limiting y HTTP 429.

## Artefactos del workflow

- `implementation-brief.md`.
- `tdd-test-plan.md`.
- `implementation-summary.md`.
- `architecture-review.md`.
- `security-review.md`.
- `final-summary.md`.

## Human gate

No ejecutar el merge, crear tests ni modificar código productivo hasta recibir
aprobación humana explícita de `tdd-test-plan.md`.
