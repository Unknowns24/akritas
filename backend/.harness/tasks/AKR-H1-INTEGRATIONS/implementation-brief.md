# Implementation Brief

## Task

`AKR-H1-INTEGRATIONS`: resolver AKR-5..12 y AKR-21 como el incremento de
integraciones del Hito 1.

## Estado actual

- El módulo Go 1.26 contiene el dominio de `GitHubAccount`,
  `GitHubRepository`, `DokployServer`, `DokployApplication`, `Project`, estados
  de integración y errores enriquecidos.
- `cmd/main.go` está vacío. No existen implementaciones de puertos, usecases,
  repositorios, migraciones, adapters externos ni REST runtime.
- `docs/openapi.yaml` es el contrato canónico y ya declara los endpoints y DTOs
  requeridos por este incremento.
- El repositorio no había elegido motor de persistencia. Para este incremento se
  fija PostgreSQL con GORM y gormigrate; el Credential Store utilizará el mismo
  PostgreSQL, no un SQLite separado.
- La línea base pasa `go test ./...`, `go vet ./...` y los gates de arquitectura,
  OpenAPI y seguridad.

Fuentes autoritativas revisadas:

- `AGENTS.md`, harness, profile `backend_api`, policies requeridas y workflow.
- `../docs/backend-architecture.md`, `domain.md`, `integrations.md`, `mvp.md`,
  `product-backlog.md`, `spec.md` y `configuration.md`.
- ADR-005, ADR-008 y ADR-009, además del resto de ADRs aceptados como límites de
  hitos posteriores.
- `docs/openapi.yaml`, `docs/errors/aaa-map.md` y memoria del harness.
- Documentación oficial vigente de GitHub App Manifest/GitHub Apps y Dokploy API.

## Estrategia y orden de implementación

1. Crear la fundación PostgreSQL, migraciones, transacciones y Credential Store
   compartido para PB-002/PB-006.
2. Implementar GitHubAccount mediante PAT y validación de identidad para
   PB-001/PB-003.
3. Implementar el flujo GitHub App Manifest de PB-066 reutilizando la cuenta y
   el Credential Store.
4. Implementar discovery de repositorios para PAT y App en PB-004.
5. Implementar DokployServer, credencial y connection test para
   PB-005/PB-006/PB-007.
6. Implementar discovery de aplicaciones para PB-008.
7. Agregar REST, paginación firmada, wiring disponible y alineación de
   OpenAPI/documentación.

PB-004 será consumido por PB-010 y PB-008 por PB-011. Este incremento no
implementará persistencia ni endpoints de Project.

## Arquitectura y módulos

### Reutilización

- Entidades, constructores y estados de `internal/core/domain`.
- `domain.Error` y catálogo `docs/errors/aaa-map.md`.
- Paths, schemas, envelopes y seguridad de `docs/openapi.yaml`.
- UUID, UTC/RFC 3339 y convención de cursor firmado.
- `AKRITAS_MASTER_KEY`, `AKRITAS_PUBLIC_URL` y
  `AKRITAS_PAGINATION_SECRET` definidos en configuración.

### Puertos de entrada

Crear contratos cohesionados para:

- administrar GitHubAccount;
- conectar mediante PAT o GitHub App Manifest;
- probar la conexión y listar repositorios GitHub;
- administrar y probar DokployServer;
- listar aplicaciones Dokploy.

Los comandos podrán transportar secretos write-only durante la llamada, pero
ningún secreto formará parte de una entidad o resultado público.

### Puertos de salida

- Repositorios de GitHubAccount y DokployServer.
- `CredentialStore` y unidad transaccional de integraciones.
- Repositorio de intentos GitHub App y bindings App/installation.
- Capabilities GitHub y Dokploy orientadas al negocio, sin copiar APIs completas
  de proveedores.
- `IntegrationUsageReader` para que Project pueda bloquear borrados cuando
  PB-010/PB-011 exista.

### Implementaciones

- Usecases en `internal/usecase/githubaccount/` y
  `internal/usecase/dokployserver/`, con una operación pública por archivo.
- PostgreSQL en `internal/adapter/db/postgres/`, con modelos y mapping exclusivos
  del adapter.
- Credential Store/cipher como adapter de infraestructura.
- GitHub y Dokploy en `internal/adapter/external/` usando `net/http`.
- DTOs y handlers feature-oriented en `internal/adapter/rest/`.
- Router construido con una dependencia obligatoria de middleware de
  administrador. La implementación concreta de PB-061..063 queda fuera.
- Cursor HMAC-SHA256 que firme integración, posición, filtros y límite.
- Wiring en el orden config → DB/migrations → repositories → adapters →
  usecases → router, sin exponer rutas sin autenticación.

## Credential Store y persistencia

Migraciones previstas:

1. `20260822_01_add_github_accounts`.
2. `20260822_02_add_dokploy_servers`.
3. `20260822_03_add_integration_credentials`.
4. `20260822_04_add_github_app_registrations`.
5. `20260822_05_add_github_app_bindings`.

Cada migración tendrá ID inmutable, registro explícito y rollback. No se usará
`AutoMigrate` global ni se agregarán data migrations.

`integration_credentials` guardará únicamente UUID, owner type/ID, secret kind,
ciphertext, nonce, versión y timestamps. Se aplicará unicidad por owner/kind.

Reglas criptográficas:

- `AKRITAS_MASTER_KEY` será Base64 de exactamente 32 bytes y se validará en
  startup.
- AES-256-GCM, nonce aleatorio por escritura y AAD compuesto por owner, secret
  kind y versión.
- Clases separadas para PAT, App private key, App webhook secret y Dokploy API
  key.
- El GitHub client secret se descartará porque H1 no implementa OAuth de usuario.
- Los installation tokens permanecerán sólo en memoria hasta su expiración menos
  un margen de seguridad.
- Metadata y secretos se crearán, rotarán y borrarán dentro de una misma
  transacción PostgreSQL, pero nunca se mantendrá una transacción abierta durante
  una llamada externa.

## GitHub

- Fijar el host productivo a `https://api.github.com`, versión explícita de API,
  timeouts y redirects seguros.
- PAT personal: validar `/user` y coincidencia case-insensitive del identifier.
- PAT de organización: validar token y acceso a la organización configurada.
- PAT classic: comprobar scopes declarados cuando estén disponibles.
- PAT fine-grained: validar identidad y discovery; el acceso al repositorio se
  comprobará otra vez al seleccionarlo porque GitHub no permite introspección
  completa de permisos write sin efectos externos.
- Manifest usará dos states diferentes: conversión e instalación. Cada uno será
  aleatorio, persistido sólo como SHA-256, expirable en una hora y consumible una
  vez.
- El callback de conversión reclamará el state antes de llamar a GitHub; un fallo
  posterior marcará el intento como fallido y requerirá reiniciar el flujo.
- Private key y webhook secret se cifrarán antes del redirect de instalación.
- El `installation_id` se verificará con JWT de la App antes de crear el
  GitHubAccount conectado.
- La App será privada, sin webhooks y con metadata read, contents read/write,
  issues write y pull requests write.

## Dokploy

- `base_url` será un origen normalizado sin userinfo, path de API duplicado,
  query ni fragment.
- HTTPS podrá apuntar a destinos públicos o privados. HTTP sólo podrá apuntar a
  loopback/RFC1918/private IPv6.
- Se bloquearán IPs link-local, metadata, multicast y unspecified; el destino se
  revalidará al conectar y los redirects prohibidos se rechazarán.
- La API key se enviará exclusivamente en `x-api-key`.
- Connection test: `GET /api/settings.health`.
- Discovery: `GET /api/application.search`.
- `application_identifier` mapeará `applicationId`; `instance_identifier`, el
  `appName`/identificador runtime disponible. Estados desconocidos mapearán a
  `unknown`.
- `server_identifier` será el SHA-256 hexadecimal del origen normalizado porque
  Dokploy no documenta un identificador estable de su instancia.

## REST y OpenAPI

Endpoints GitHub:

- `GET|POST /integrations/github/accounts`.
- `GET|PATCH|DELETE /integrations/github/accounts/{account_id}`.
- `POST /integrations/github/accounts/{account_id}/connection-test`.
- `GET /integrations/github/accounts/{account_id}/repositories`.
- `POST /integrations/github/app-manifest/registrations`.
- `GET /integrations/github/app-manifest/callback`.
- `GET /integrations/github/app-installations/callback`.

Endpoints Dokploy:

- `GET|POST /integrations/dokploy/servers`.
- `GET|PATCH|DELETE /integrations/dokploy/servers/{server_id}`.
- `POST /integrations/dokploy/servers/{server_id}/connection-test`.
- `GET /integrations/dokploy/servers/{server_id}/applications`.

No se implementarán ni modificarán `POST /projects` o
`PATCH /projects/{project_id}`.

Los callbacks GitHub serán públicos, con `state` y `Cache-Control: no-store`.
Las demás operaciones exigirán sesión administrativa; las mutaciones también
validarán Origin.

Al agregar respuestas `403` compatibles y aclarar el significado de
`server_identifier`, subir `info.version` a `1.1.0` y actualizar el gate/memoria
que hoy fijan `1.0.0`.

## Inconsistencias detectadas

- `mvp.md` clasifica connection tests/discovery como Nice to Have mientras el
  backlog y OpenAPI los definen P0 de H1.
- `spec.md` puede leerse como si Project almacenara credenciales, contradiciendo
  dominio y ADR-005.
- `domain.md` omite proyecciones seguras ya presentes en código/OpenAPI:
  `credential_configured`, conteos, last checked/synced, URLs y status.
- `server_identifier` es obligatorio en response pero no tiene input ni origen
  estable documentado por Dokploy.
- Las mutaciones no documentan uniformemente el `403` por Origin inválido.
- `check-openapi.sh` fija `1.0.0`, en tensión con la policy SemVer.
- El guard de borrado contra Projects está contratado, pero la persistencia
  Project aún no existe. Se definirá y probará el puerto ahora; PB-010/PB-011
  completará su adapter concreto antes de habilitar Projects.

## Riesgos y mitigaciones

- Secret leakage: tipos separados, DTOs write-only, redacción, límites de body y
  tests negativos.
- SSRF Dokploy: validación DNS/IP por request y redirects cerrados.
- Replay de callbacks: states distintos, digest, CAS, expiración y consumo único.
- Fallo después de convertir un Manifest: persistencia transaccional y estado
  failed; el usuario deberá eliminar la App huérfana o reiniciar el flujo.
- Permisos fine-grained no introspectables: validar identidad/discovery y repetir
  la autorización sobre el repositorio en el caso de uso que lo seleccione.
- Dependencia auth aún ausente: constructor fail-closed y rutas no habilitadas sin
  middleware concreto.

## Human gate

No crear tests ni implementación hasta la aprobación humana explícita de
`tdd-test-plan.md`.
