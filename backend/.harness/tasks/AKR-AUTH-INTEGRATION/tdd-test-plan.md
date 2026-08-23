# TDD Test Plan

## Scope

Definir mediante tests la integración de setup/status/verify/login/session y
logout con la arquitectura vigente, Credential Store compartido, transacciones
PostgreSQL y protección del router de integraciones.

Recovery, RBAC, rate limiting distribuido y compatibilidad con bases no
descartables quedan fuera de alcance.

Los tests deben describir comportamiento e invariantes; no fijarán nombres de
helpers privados ni un orden de llamadas sin efecto observable.

## 1. Configuración Viper

- Environment prevalece sobre `app.env` también para valores auth.
- Defaults: idle TTL 12h, absolute TTL 168h y cookie secure habilitada.
- Bootstrap token ausente, TTLs no positivos, idle mayor que absolute, origin
  inválido o cookie insegura bajo public URL HTTPS fallan cerrado.
- Origins son exactos, sin wildcard, userinfo, query ni fragment.
- `AKRITAS_DATABASE_URL` continúa siendo el único DSN soportado;
  `AKRITAS_DB_DSN` no afecta la configuración.
- Master key, pagination secret y bootstrap token no aparecen en errores y sus
  representaciones raw se limpian del resultado.

## 2. Dominio y tags persistibles

- `PendingEnrollment` normaliza email/display name, valida UUID/timestamps y
  considera expirado `now >= expires_at`.
- `AdministratorSession.ExtendIdle` desliza el vencimiento sin superar el
  absoluto y rechaza transiciones inválidas.
- El período TOTP consumido nunca disminuye mediante comportamiento de dominio o
  repositorio.
- Administrator, enrollment y session contienen tags estructurales esperados,
  sin importar GORM.
- Ninguna de esas entidades contiene password hash, token hash, plaintext TOTP,
  ciphertext o nonce.
- No existe un package de modelos PostgreSQL que duplique esas entidades.

## 3. Migraciones PostgreSQL

- El registry contiene exactamente `01..08` en orden y cada migración declara
  migrate/rollback.
- La tercera migración crea `credentials`, no
  `integration_credentials`, con unicidad owner/type/kind y checks de nonce y
  ciphertext.
- Las migraciones `06..08` crean administrators, pending enrollments y sessions
  con sus constraints, índices y foreign keys.
- Una base nueva ejecuta y revierte el registry completo.
- No existe `AutoMigrate` global ni migración auth basada en un modelo Go vivo.
- PostgreSQL persiste/lee directamente los campos de las entidades de dominio;
  los campos técnicos extra se proyectan sin structs completos duplicados.

## 4. Credential Store y transactor

- GitHub/Dokploy conservan round-trip, rotación, delete y AAD sobre la tabla
  renombrada.
- Un seed TOTP se guarda bajo pending enrollment, nunca como plaintext.
- MoveOwner re-encripta el seed con AAD de administrator y elimina el owner
  anterior.
- Get con owner, kind, versión o key incorrectos falla sin revelar material.
- Repositorios y Store usan la transacción activa propagada por contexto.
- Commit conserva todas las escrituras participantes; un error intermedio
  revierte metadata, sesión y credentials.
- Un callback nil/configuración inválida falla de forma estable.
- Operaciones existentes con transacciones internas reutilizan correctamente
  una transacción externa o savepoint sin escapar del rollback.

## 5. Setup y setup status

- Sin administrator, setup status devuelve `setup_required=true` y
  `registration_open=true`; con administrator ambos son false.
- Rate limit se evalúa antes de verificar bootstrap y usar DB.
- Bootstrap token usa comparación constante y un valor incorrecto produce el
  error público genérico esperado.
- Administrator existente impide generar/persistir otro enrollment.
- Password 12..128 se hashea con Argon2id y parámetros ADR-008; nunca se
  persiste plaintext.
- Setup genera TOTP RFC 6238 y devuelve provisioning material una sola vez con
  expiración de 10 minutos.
- Reemplazar enrollment elimina atómicamente metadata, hash y seed anteriores.
- Fallar al guardar hash o seed revierte el nuevo enrollment y conserva un
  estado coherente.
- Setup y verify tienen budgets independientes.

## 6. Verificación TOTP y activación

- Enrollment ID inválido, inexistente, reemplazado, expirado o código incorrecto
  producen el mismo error seguro.
- Se acepta el período actual y la tolerancia RFC 6238 de uno anterior/posterior.
- Verify exitoso crea exactamente un Administrator y una sesión, mueve el seed
  al owner definitivo y consume el enrollment.
- La respuesta contiene administrator/session seguros y el token sólo viaja en
  cookie.
- Fallo en create administrator, move secret, save session o delete enrollment
  revierte toda la activación.
- Dos verifies concurrentes no crean dos administradores.
- Repetir el enrollment consumido no crea otra sesión.

## 7. Login y replay TOTP

- Email inexistente, password incorrecta, hash inválido, seed inaccesible, TOTP
  incorrecto o replay devuelven el mismo 401 público.
- Argon2id valida versión y límites de memory/iterations/parallelism antes de
  derivar; un hash manipulado no puede forzar recursos no acotados.
- Login aplica budgets independientes por IP y cuenta normalizada.
- Un período TOTP sólo es válido si es estrictamente posterior al persistido.
- Un período anterior tolerado nunca hace retroceder el estado.
- Dos logins concurrentes con el mismo período producen una sola sesión.
- Si guardar sesión falla, el compare-and-set del período se revierte.
- Login exitoso no revoca sesiones anteriores y crea un token opaco nuevo.

## 8. Sesiones, middleware y Origin

- Token vacío/desconocido, sesión revocada, idle expirada o absolute expirada
  produce unauthorized sin reflejar la cookie.
- Sesión activa desliza idle TTL y lo limita al absolute TTL.
- Get session devuelve administrator y tiempos seguros, nunca token hash.
- Logout revoca sólo la sesión actual y expira la cookie.
- Mutaciones autenticadas sin Origin, con Origin distinto o wildcard son
  rechazadas con 403 antes del handler.
- Origin público y cada origin exacto allowlisted son aceptados.
- Lecturas protegidas requieren sesión pero no validación CSRF de mutación.

## 9. Security adapters

- Bootstrap verifier rechaza token configurado/candidato vacío y acepta sólo
  coincidencia exacta.
- Argon2id usa salt aleatorio, formato versionado y comparación constante.
- TOTP generator/verifier comparten seis dígitos, SHA-1 y 30 segundos.
- Session token usa 32 bytes de `crypto/rand`, Base64URL y SHA-256 persistible.
- Rate limiter mantiene namespaces y budgets aislados, reinicia ventanas
  expiradas y limpia buckets antiguos.
- Ningún error/log contiene password, bootstrap token, seed TOTP, session token,
  master key, ciphertext o nonce.

## 10. REST, DTOs y errores

- DTOs auth coinciden campo por campo con OpenAPI, usan sufijo `DTO`, una
  estructura por archivo y mappers separados.
- Decoder rechaza JSON vacío, desconocido, múltiple o mayor al límite.
- Setup/verify/login respetan status, envelopes, `Cache-Control: no-store`,
  `Set-Cookie`, HttpOnly, Secure, SameSite=Lax y Path=/.
- Logout devuelve 204 y cookie expirada.
- Error `R` conserva código estable, envelope común, `Retry-After` y HTTP 429.
- Causas DB/crypto no aparecen en respuestas internas.
- Request IDs siguen la implementación común vigente.

## 11. Router, bootstrap y regresión de integraciones

- Auth público monta setup-status, setup, verify y login.
- Session GET/DELETE requiere el middleware concreto de auth.
- Callbacks GitHub continúan públicos y protegidos por state.
- Todas las rutas GitHub/Dokploy existentes continúan montadas detrás del mismo
  middleware de sesión; sus mutaciones también validan Origin.
- Bootstrap abre un solo pool PostgreSQL, ejecuta migraciones una vez y comparte
  cipher/Credential Store entre auth e integraciones.
- Cualquier dependencia auth/DB/config ausente hace fallar el build del router
  antes de exponer integraciones.
- Los tests actuales de integraciones y paginación continúan pasando.

## 12. OpenAPI, policy y documentación

- `ErrorCode` admite `R` y el mapper/policy lo asocian exclusivamente a 429.
- `info.version` es `1.3.0`; paths, requests y response payloads auth no cambian.
- El catálogo registra errores auth/REST/DB en su capa y no duplica códigos.
- ADR-014 explica usos, límites, propagación, nested transactions y tests del
  Transactor.
- ADR-005/008/010..013 no son contradichos por dominio, schema o wiring.
- Artefactos históricos de auth permanecen; el summary de integración enumera
  implementaciones descartadas y reemplazos.

## Expected failing tests before implementation

- La rama milestone no contiene usecases/repositorios/handlers concretos de
  auth ni middleware admin.
- El Credential Store no reconoce owners TOTP ni una transacción propagada.
- La tabla aún se llama `integration_credentials`.
- No existen migraciones auth `06..08`, Transactor ni ADR-014.
- Config no carga bootstrap/session/origins.
- ErrorCode no admite `R` y OpenAPI continúa en `1.2.0`.
- El router falla cerrado porque auth todavía no provee el middleware.

## Secuencia TDD

1. Tras la aprobación, ejecutar el merge sin commit y resolver únicamente lo
   necesario para que el árbol compile como punto de partida.
2. Escribir/ajustar primero tests de config, dominio, migrations/transactor,
   Credential Store y security adapters; confirmar fallos relevantes.
3. Implementar infraestructura mínima para esos tests.
4. Escribir tests de usecases setup/verify/login/session y confirmar fallos.
5. Implementar usecases y puertos hasta verde.
6. Escribir tests REST/middleware/router/bootstrap y confirmar fallos.
7. Implementar el boundary HTTP y wiring.
8. Ejecutar tests PostgreSQL reales, race detector y gates; corregir mediante el
   fix-plan requerido si reviews encuentran fallos no triviales.

## Validaciones finales

- `go test ./...`.
- `go test -race ./...`.
- `go test -tags=integration ./...` con PostgreSQL efímero.
- `go vet ./...`.
- `gofmt` y comprobación del merge/worktree.
- `.harness/kernel/scripts/check-backend-architecture.sh`.
- `.harness/kernel/scripts/check-openapi.sh`.
- `.harness/kernel/scripts/check-security.sh`.
- Architecture review y security review del workflow.

## Acceptance criteria

- Setup, verify, login, current session y logout cumplen OpenAPI y ADR-008.
- El seed TOTP usa exclusivamente el Credential Store PostgreSQL actual.
- Dominio no contiene material secreto ni modelos duplicados.
- Todas las invariantes multi-repository son atómicas y replay-safe.
- Integraciones existentes quedan accesibles sólo con sesión válida.
- `R` representa 429 de forma consistente en code, policy y OpenAPI.
- Todas las validaciones pasan y el merge queda resuelto sin commit.

## Human approval required

Este archivo requiere aprobación humana explícita. No ejecutar el merge, crear
tests ni implementar código antes de recibirla.
