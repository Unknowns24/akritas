# Implementation Brief

## Task

`AKR-AUTH-INTEGRATION`: integrar limpiamente `origin/feat/authentication` sobre
la rama milestone actual sin crear el merge commit.

## Estado actual

- `feat/backend-milestone-1` pasa `go test ./...` antes del merge.
- La configuración vigente vive en `config/config.go`, usa una instancia Viper
  local, `AKRITAS_DATABASE_URL`, environment sobre `app.env` y secretos
  derivados validados.
- El adapter PostgreSQL actual persiste directamente las entidades de dominio
  autorizadas por ADR-012 y ya posee Credential Store AES-256-GCM con AAD.
- El router de integraciones usa `net/http` y falla cerrado sin middleware de
  administrador; `cmd/main.go` permanece vacío hasta que auth provea ese
  boundary.
- `origin/feat/authentication` implementa setup/status/verify/login/session y
  logout, pero trae configuración manual, modelos GORM duplicados, un segundo
  cipher denominado Credential Store, router Chi y helpers REST paralelos.
- El merge sintético detecta conflictos materiales en configuración,
  dependencias, migraciones, router, catálogo de errores y Credential Store.

## Fuentes autoritativas

- Harness, profile `backend_api`, workflow y todas sus policies requeridas.
- ADR-005, ADR-008 y ADR-010..013.
- `../docs/backend-architecture.md`, `configuration.md` y memoria del proyecto.
- `docs/openapi.yaml` como contrato canónico.
- La decisión humana registrada en `task.md` sobre merge, migraciones,
  transactor, credential table y error `R`.

## Estrategia de integración

1. Tras aprobar el plan TDD, ejecutar
   `git merge --no-commit origin/feat/authentication` y conservar `AGENTS.md`
   sin alteraciones adicionales.
2. Resolver los conflictos usando como base la arquitectura de la rama actual;
   no aceptar archivos completos por estrategia ours/theirs cuando haya
   conflicto semántico.
3. Conservar los usecases y comportamiento observable de auth, reescribiendo
   sus adapters, persistence y REST para cumplir las convenciones vigentes.
4. Mantener los artefactos históricos de las tres tareas auth de la rama y
   registrar en este task qué decisiones quedaron reemplazadas.
5. Dejar el índice resuelto y el working tree dentro del merge, validado y sin
   commit.

## Dominio y persistencia

- `Administrator`, `PendingEnrollment` y `AdministratorSession` serán las
  estructuras activas de GORM mediante tags pasivos. Ninguna importará GORM ni
  contendrá password hash, token hash, ciphertext, nonce o seed TOTP.
- El período TOTP consumido es estado de autenticación no secreto y podrá formar
  parte de `Administrator` para hacer explícita la regla anti-replay.
- Password hashes y token hashes continuarán en columnas de sus tablas, pero los
  repositorios los manejarán mediante proyecciones privadas mínimas o updates
  explícitos, sin modelos completos duplicados.
- `PendingEnrollment` persistirá identidad, email, display name y expiración;
  su password hash será material privado del repositorio y el seed TOTP vivirá
  cifrado en Credential Store.

Migraciones finales, sobre base descartable:

1. Mantener GitHub accounts y Dokploy servers como `01` y `02`.
2. Reemplazar `03_add_integration_credentials` por
   `03_add_credentials`, con tabla/constraints genéricos.
3. Mantener GitHub App registrations/bindings como `04` y `05`.
4. Agregar `06_add_administrators`.
5. Agregar `07_add_pending_enrollments`.
6. Agregar `08_add_administrator_sessions`.

Las migraciones usarán SQL estable, registro explícito y rollback. No se
conservarán los IDs `01..05` experimentales de la rama auth ni `AutoMigrate`
histórico dependiente del struct Go actual.

## Credential Store

- Conservar el cipher y repository actuales; eliminar
  `adapter/security/credential_store.go` de la rama.
- Renombrar usos de `integration_credentials` a `credentials`.
- Agregar `CredentialOwnerPendingEnrollment`,
  `CredentialOwnerAdministrator` y `SecretKindAdministratorTOTP`.
- El setup guarda el seed bajo el enrollment; verify lo mueve al administrator,
  provocando re-encryption con el nuevo AAD.
- Password hashes permanecen irreversibles mediante Argon2id y no se cifran ni
  se almacenan en Credential Store.

## Transacciones

Se incorporará ADR-014 para `out.Transactor.WithinTransaction(ctx, fn)`:

- Sólo coordina operaciones cortas sobre el mismo PostgreSQL.
- El handle GORM activo se propaga mediante una key privada del adapter.
- Repositorios PostgreSQL y Credential Store resuelven consistentemente el
  handle activo antes de operar.
- No se permiten HTTP, providers, filesystem, generación Argon2/TOTP ni otro
  trabajo lento dentro del callback.
- Un repositorio puede conservar una transacción local para una única operación
  cohesiva, pero debe reutilizar la transacción exterior cuando exista.

Boundaries atómicos:

- reemplazar enrollment previo, su hash y seed por el nuevo enrollment;
- crear administrator y sesión, mover seed y consumir enrollment;
- actualizar condicionalmente el último período TOTP y crear la sesión.

El login usará compare-and-set `candidate_period > last_accepted_period`. Dos
requests concurrentes con el mismo período no podrán crear dos sesiones y un
período tolerado anterior nunca hará retroceder el valor persistido.

## Seguridad y configuración

- Ampliar el `Config` Viper vigente con bootstrap token, TTL idle/absolute,
  cookie secure y allowlist exacta de origins; no introducir
  `AKRITAS_DB_DSN`.
- Mantener la master key Base64 de 32 bytes y limpiar representaciones crudas
  antes de devolver `Config`.
- Separar adapters por responsabilidad: bootstrap verifier, password Argon2id,
  TOTP, session token y fixed-window rate limiter.
- Inyectar `func() time.Time` en usecases en lugar de un Clock adapter ubicado
  dentro de security.
- Validar y acotar versión/parámetros Argon2 antes de derivar memoria.
- Mantener budgets distintos para setup, verify y login; login usa keys
  namespaceadas por IP y email normalizado.
- El limiter elimina buckets expirados para evitar crecimiento ilimitado.
- Los tokens de sesión son 32 bytes aleatorios Base64URL y sólo su SHA-256 se
  persiste.
- Las mutaciones autenticadas validan Origin contra el origen público y la
  allowlist configurada antes de ejecutar el handler.

## REST, errores y contrato

- Conservar `net/http`, request decoder, response envelope y request ID
  vigentes; eliminar router Chi, root DTO error y response helpers duplicados.
- DTOs de auth tendrán sufijo `DTO`, una estructura por archivo bajo
  `rest/dto/auth`; conversiones vivirán bajo `rest/mapper`.
- Los endpoints auth públicos se montan junto a callbacks GitHub públicos.
  Session y todas las integraciones usan el mismo middleware admin.
- No se agregan endpoints ni campos. Recovery continúa sin implementación.
- Agregar `R` al formato `DxAAABBBT`, mapearlo a 429 y actualizar policy,
  catálogo, OpenAPI y `info.version` a `1.3.0`.

## Wiring

El bootstrap final seguirá:

```text
config
→ PostgreSQL + migraciones
→ cipher + Credential Store compartido
→ repositorios auth/integraciones + transactor
→ security/external adapters
→ usecases auth/integraciones
→ auth middleware + router único
→ HTTP server
```

No se abrirán dos pools ni se ejecutarán migraciones desde módulos separados.

## Riesgos y mitigaciones

- Replay TOTP concurrente: compare-and-set dentro de la transacción de login.
- Setup parcialmente consumido: activación, secret move y delete son atómicos.
- Secret leakage: DTOs write-only, store actual con AAD, errores genéricos y
  tests negativos.
- CSRF: Origin obligatorio para mutaciones con cookie.
- Crecimiento del limiter: expiración y limpieza de buckets.
- Regresión de integraciones: router, bootstrap y tests verifican que todos los
  endpoints existentes permanezcan montados y protegidos.
- Historia engañosa: los summaries históricos se conservan, pero este task
  documentará las decisiones sustituidas.

## Human gate

No ejecutar el merge, crear tests ni implementar código hasta la aprobación
humana explícita de `tdd-test-plan.md`.
