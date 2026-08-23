# Implementation Summary

## Implemented

- Merge real de `origin/feat/authentication` abierto con `--no-commit` y
  resolución semántica sobre la arquitectura de `feat/backend-milestone-1`.
- Config Viper única con `AKRITAS_DATABASE_URL`, bootstrap token, TTLs de
  sesión, cookie segura y origins exactos; se eliminó `AKRITAS_DB_DSN`.
- Persistencia directa de `Administrator`, `PendingEnrollment` y
  `AdministratorSession` mediante tags GORM pasivos. Hashes de password/token y
  ciphertext continúan fuera de las entidades de dominio.
- Credential Store vigente ampliado con owners de auth y secret kind TOTP. La
  tabla fundacional ahora es `credentials` y el seed se re-cifra con nuevo AAD
  al pasar de enrollment a administrator.
- Migraciones explícitas `06..08` para auth, sin `AutoMigrate`, y eliminación de
  los modelos/migraciones experimentales duplicados de la rama.
- `out.Transactor` propagado mediante contexto privado del adapter y ADR-014.
  Setup, activación y login poseen boundaries atómicos compartidos con el Store.
- Replay TOTP corregido con compare-and-set estricto y creación de sesión en la
  misma transacción; un período igual o anterior nunca se acepta.
- `adapter/security` dividido en bootstrap, password, TOTP, session token y
  rate limiter; se eliminó el segundo Credential Store y el Clock adapter.
- Router único `net/http`, middleware de sesión y Origin para mutaciones
  autenticadas, y bootstrap compartido por auth e integraciones sobre un solo
  pool/migrador/Store.
- DTOs auth con sufijo `DTO`, decoder/envelopes comunes y mappers dedicados.
- Tipo de error `R` mapeado a HTTP 429; catálogo, policy, OpenAPI y gate
  actualizados a `1.3.0`.

## Implementaciones históricas reemplazadas

Se conservan los artefactos `AKR-AUTH-BOOTSTRAP`, `AKR-AUTH-TOTP-VERIFY` y
`AKR-AUTH-LOGIN-SESSION` como historial, pero esta integración reemplaza:

- loader manual y `AKRITAS_DB_DSN` por el Config Viper vigente;
- modelos PostgreSQL duplicados por entidades de dominio persistibles;
- cipher/credential store de `adapter/security` por el Store AES-GCM compartido;
- migraciones auth experimentales `01..05` por `06..08` posteriores;
- Chi y helpers REST paralelos por `net/http`, decoder, envelopes y mappers
  actuales;
- actualización no condicional del período TOTP por compare-and-set;
- Clock adapter por `func() time.Time` inyectable;
- gap de Origin previamente documentado por validación fail-closed.

## Deliberadamente no implementado

- Recovery y rate limiting distribuido/persistente.
- Compatibilidad de datos con `integration_credentials` o migraciones auth
  experimentales, conforme a la decisión de recrear la base descartable.
- Merge commit.

## Validation result

- `go test ./...`: pass.
- `go test -race ./...`: pass.
- `go test -tags=integration ./internal/adapter/db/postgres/...`: pass.
- `go vet ./...`: pass.
- `gofmt` y `git diff --check`: pass.
- Gates backend architecture, OpenAPI (59 operaciones, 112 schemas) y security:
  pass.
