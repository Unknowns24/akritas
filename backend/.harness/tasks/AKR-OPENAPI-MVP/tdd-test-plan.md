# TDD Test Plan

## Scope

Validar primero el contrato OpenAPI y la documentacion de comportamiento. No se testean handlers ni adapters porque no forman parte de esta tarea.

## Tests to add/update

### 1. Estructura OpenAPI

- El YAML parsea como objeto y contiene `openapi: 3.1.0`, `info.version: 1.0.0`, `servers`, `paths` y `components.schemas`.
- Cada operacion tiene `operationId` no vacio y unico.
- Cada parametro de path declarado en una ruta existe en la operacion y tiene `required: true`.
- Todo `$ref` local apunta a un nodo existente.
- Todos los status codes y response objects tienen contenido o descripcion valida.

### 2. Seguridad y exposicion publica

- El security default es `cookieAuth`.
- Solo son publicos: health, readiness, setup status/setup/verify, login, recovery/verify y callbacks GitHub.
- Logout y consulta de sesion requieren cookie.
- Mutaciones autenticadas documentan 401; recursos restringidos documentan 403 cuando aplica.
- Setup, login y recovery documentan 429 y errores genericos.

### 3. Secret isolation

- Request schemas marcan como `writeOnly` password, bootstrap token, PAT y credenciales Dokploy/QVAC.
- `totp_code`, `otpauth_uri` y manual entry secret solo aparecen en schemas especificos de enrollment/verification y nunca en recursos persistentes.
- Ningun response schema reutilizable de Administrator, GitHubAccount, DokployServer, QVACConfiguration, Project, Incident, Evidence o Session contiene nombres sensibles prohibidos.
- Responses de integraciones solo exponen `credential_configured` y estados normalizados.
- Evidence, log context, diffs y validation output se documentan como redactados/sanitizados.

### 4. Contrato de autenticacion

- Setup cerrado devuelve 409; setup valido devuelve enrollment de un solo uso con `Cache-Control: no-store`.
- Verify activa al unico Administrator y crea cookie; login requiere email/password/TOTP.
- Logout devuelve 204; session devuelve proyeccion segura del Administrator.
- Recovery requiere bootstrap token, password nuevo y nuevo enrollment TOTP; verify revoca sesiones previas.
- Schemas fijan password 12..128, TOTP de seis digitos y expiraciones documentadas.

### 5. GitHub, Dokploy y QVAC

- GitHub soporta PAT y App Manifest sin mezclar `account_type` con `authentication_method`.
- Manifest registration devuelve action/manifest/state/expiry; callbacks requieren code/state o installation/state y responden 303.
- Repositorios GitHub y aplicaciones Dokploy son seleccionables mediante endpoints documentados.
- Tests de conexion tienen resultados normalizados y nunca retornan secretos.
- QVAC permite none/bearer/basic, declara restriccion a hosts locales/privados y prohibicion de redirects publicos.

### 6. Projects, Monitoring y Automation

- Projects referencian una cuenta/repository GitHub y un server/application Dokploy sin credenciales.
- MonitoringConfiguration contiene exactamente los seis campos aprobados.
- Defaults: disabled, arrays vacios, `PT30M`, context before/after 20.
- Regex invalidas y duraciones/contextos invalidos documentan 400.
- Automation defaults true; remediation implica investigation y automatic pull request implica remediation.
- La Issue obligatoria no aparece como toggle.
- Reglas built-in son read-only y no existe detector basado en IA.

### 7. Incident lifecycle y operaciones asincronas

- Enums coinciden con root cause, resolution, investigation, remediation, phase y terminal outcome aprobados.
- Incidents exponen identificador humano, project, severity, fingerprint, occurrence count, first/last seen y phase.
- Detalle incluye latest investigation, issue reference, remediation/validation y PR opcional.
- LogEvents, Evidence, timeline, validation results, Pull Requests y Activity tienen colecciones paginadas cuando corresponde.
- Comandos manuales de investigation, remediation y PR requieren `Idempotency-Key` y devuelven 202 `Operation`.
- Operation soporta queued/running/succeeded/failed y referencia segura al recurso resultante.
- `requires_human`, remediation failed y PR created son outcomes terminales sin merge/deploy.

### 8. Paginacion, errores y cobertura funcional

- Toda coleccion operacional reutiliza `Paging` y documenta que un cursor posterior no se combina con filtros.
- Se documentan filtros necesarios para Projects, Incidents, Activity, repositories, applications y Pull Requests.
- Error envelope respeta el patron `DxAAABBBT` y no filtra errores de proveedores.
- Una matriz en la documentacion de la tarea vincula PB-001..PB-055 con operaciones o schemas concretos.
- Overview cubre proyectos monitoreados, incidentes activos, workflows completados, PR creadas e investigaciones activas.
- Cada pantalla compatible del ZIP tiene datos contractuales; Team, New Incident y metricas de “resolved” quedan excluidos.

## Expected failing tests before implementation

- El gate actual omite validacion porque `docs/openapi.yaml` no existe.
- No existen paths, schemas, security scheme, operation IDs ni matriz de cobertura.
- No existen ADR/documentacion para auth TOTP ni GitHub App Manifest.

Tras la aprobacion, se reforzara primero el gate para que falle por ausencia/incompletitud del contrato; luego se agregaran OpenAPI y documentacion hasta satisfacerlo.

## Acceptance criteria covered

- Contrato completo y consumible por frontend para H1-H6.
- Autenticacion single-admin con bootstrap/password/TOTP/session/recovery.
- Integraciones seguras GitHub PAT/App, Dokploy y QVAC.
- Control plane, monitoring, detection projections e incident lifecycle completos.
- Issue obligatoria, validacion previa a PR y frontera humana preservadas.
- Ausencia estructural de secretos en responses.
- Compatibilidad con la interfaz prevista sin incorporar conceptos fuera de scope.

## Open questions / human approval notes

- No quedan decisiones funcionales abiertas; se aplican los defaults aprobados en el plan.
- Se requiere aprobacion humana explicita de este archivo antes de modificar `docs/openapi.yaml`, el gate o la documentacion de producto.
- La escritura sobre `../docs/` puede requerir ampliar el permiso del workspace cuando comience la implementacion.
