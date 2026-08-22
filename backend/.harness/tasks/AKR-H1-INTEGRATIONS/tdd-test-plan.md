# TDD Test Plan

## Scope

Definir mediante tests el contrato de AKR-5..12 y AKR-21: almacenamiento seguro,
administración y validación de GitHub/Dokploy, GitHub App Manifest y discovery de
repositorios/aplicaciones.

No se probarán ni implementarán Project/Monitoring, autenticación single-admin,
logs, detección, Incidents, QVAC, Issues o Remediation.

## Corrección arquitectónica — tests antes del refactor

- Uker parsea primeras páginas y cursores firmados, rechaza tampering,
  expiración, overrides y scopes de otra integración.
- PostgreSQL aplica `pagination.Apply` para datos y `ApplyFilters` para el total,
  con orden total `created_at,id` y navegación siguiente/anterior.
- Discovery traduce boundaries firmados Uker a página GitHub u offset Dokploy
  sin codec/serialización propios.
- Viper respeta environment sobre `app.env`, aplica defaults y falla cerrado
  ante DSN, URL, master key o pagination secret inválidos.
- GORM persiste entidades de dominio contra las cinco migraciones existentes;
  ciphertext/nonce sólo aparecen en el record privado del Credential Store.
- Los errores por capa conservan códigos únicos y no aparecen en catálogos de
  otra capa.
- Los contratos REST usan DTOs tipados seguros, agrupados por feature/common y
  mappers separados sin cambiar los payloads OpenAPI; el gate rechaza DTOs
  sueltos en la raíz `rest/dto`.
- Bootstrap rechaza middleware admin ausente antes de abrir DB o ejecutar
  migraciones.

Estos tests priorizan comportamiento e invariantes; no fijan nombres de helpers
privados ni orden interno de llamadas sin efecto observable.

## 1. Dominio, contratos y errores

- `GitHubAccountType` y `GitHubAuthenticationMethod` permanecen independientes.
- Las proyecciones GitHub/Dokploy contienen sólo metadata segura y estados
  normalizados.
- Ninguna entidad, input port de lectura o DTO de salida contiene PAT, API key,
  private key, webhook secret, ciphertext o nonce.
- Repositorios y aplicaciones provider se mapean a constructores válidos del
  dominio; status desconocido de Dokploy se transforma en `unknown`.
- Los errores nuevos REST/DB/external/usecase cumplen `DxAAABBBT`, son únicos y
  coinciden con `docs/errors/aaa-map.md`.
- Las causas envueltas no aparecen en mensajes REST.

## 2. Cipher y configuración

- Una master key Base64 de 32 bytes crea un cipher válido.
- Base64 inválido, longitud distinta o key vacía fallan startup sin imprimir el
  valor.
- Cifrar y descifrar con AES-256-GCM conserva el plaintext correcto.
- Dos cifrados del mismo secreto producen nonce/ciphertext diferentes.
- AAD con owner, kind o versión diferente impide descifrar.
- Ciphertext, errores y logs nunca contienen el plaintext.
- Nonces se generan con `crypto/rand`; no se aceptan nonces suministrados por
  callers productivos.

## 3. PostgreSQL, Credential Store y migraciones

- Las cinco migraciones se registran en orden con los IDs aprobados.
- Una base vacía termina con tablas, constraints, índices y FKs esperados.
- Cada migración declara rollback; el rollback reversible elimina sólo su schema.
- No existe `AutoMigrate` fuera de una migración versionada.
- `Put` persiste ciphertext/nonce/version y nunca plaintext.
- `Get` devuelve el secreto sólo al caller interno autorizado y falla de forma
  segura con key/AAD incorrectos.
- Owner/kind es único; rotación reemplaza la fila sin cambiar la integración.
- Crear metadata+secreto y eliminar metadata+binding+secretos son atómicos.
- Un fallo inyectado a mitad de transacción no deja metadata ni credenciales
  huérfanas.
- Tests reales PostgreSQL se ejecutan con Testcontainers bajo
  `go test -tags=integration ./...`; SQLite no se usa como sustituto.

## 4. GitHub PAT usecases

- Create valida el PAT antes de persistir la cuenta o credencial.
- Auth inválida, target personal incorrecto u organización inaccesible no dejan
  filas ni secretos.
- PAT personal compara identifier sin sensibilidad a mayúsculas.
- PAT de organización exige que el token pueda consultar el target configurado.
- Classic PAT valida scopes informados; ausencia de scopes informados sigue la
  política fine-grained y no inventa permisos.
- Create exitoso devuelve `credential_configured=true`, status connected y nunca
  devuelve el PAT.
- Get/list devuelven proyecciones seguras y paginadas.
- Patch de display name no lee ni rota el secreto.
- Rotación valida primero; fallo conserva el PAT y metadata anteriores.
- Rotación exitosa reemplaza el secreto sin cambiar el account ID.
- Delete consulta `IntegrationUsageReader`; referenced produce conflict y no
  borra nada, unreferenced borra cuenta y secreto.

## 5. GitHub connection test y repositories

- Cuenta inexistente produce not found.
- Credencial válida produce `connected`, timestamp y latencia no negativa.
- Provider 401/403 produce resultado `authentication_failed`, no un 401 de la
  sesión Akritas.
- Timeout/DNS/5xx produce `unavailable` con mensaje normalizado.
- Discovery PAT devuelve sólo repositorios dentro del target configurado y
  accesibles por el token.
- Discovery App devuelve sólo repositorios de la instalación verificada.
- Cada repositorio contiene identifier, owner, name, full name, default branch,
  privacy y URL seguros.
- Name filter, limit y cursor se respetan; cursor firmado no permite cambiar
  account, filtros o límite.
- Discovery exitoso actualiza repository count/last checked; fallo no destruye el
  último conteo válido.

## 6. GitHub App Manifest

- Manifest personal y organization construyen el form action oficial correcto.
- Organization exige nombre; personal no lo acepta como target alternativo.
- Manifest contiene homepage/redirect/setup URLs derivadas exclusivamente de
  `AKRITAS_PUBLIC_URL`.
- App es privada, webhooks están inactivos y los únicos permisos son metadata
  read, contents read/write, issues write y pull requests write.
- Registration devuelve ID, form action, manifest, state y expiración con
  `Cache-Control: no-store`.
- El state se guarda sólo como digest y expira como máximo en una hora.
- State incorrecto, vencido, ya consumido o de otra etapa es rechazado.
- El callback reclama el state mediante compare-and-swap antes del exchange.
- Exchange exitoso cifra private key y webhook secret; client secret no se
  persiste.
- Fallo de persistencia después del exchange marca el intento failed y nunca
  expone los secretos.
- El redirect de instalación contiene un segundo state impredecible.
- Installation callback valida state y después verifica el installation ID con
  JWT de la App.
- Installation de otra App, suspendida o inexistente no crea GitHubAccount.
- Éxito crea binding y GitHubAccount connected en una transacción, promueve los
  secretos al owner definitivo y consume el state.
- Repetir cualquiera de los callbacks no duplica cuentas ni bindings.
- Installation token es efímero, no se persiste y se renueva antes de expirar.

## 7. GitHub external adapter

- Usa host GitHub y versión API fijados en producción.
- Tests con `httptest` validan headers Accept, Authorization y versión sin usar
  credenciales reales.
- Redirects a hosts inesperados se rechazan.
- Responses mayores al límite, JSON inválido, 401/403/404/422/429/5xx y timeouts
  se normalizan sin conservar body sensible.
- JWT contiene issuer/app ID, issued-at con skew controlado y expiración corta.
- Provider DTOs nunca salen del package adapter.

## 8. DokployServer usecases

- Create normaliza base URL, calcula fingerprint SHA-256 y valida health antes de
  persistir.
- API key inválida produce validation/authentication failure; DNS/timeout/5xx
  produce unavailable.
- Create exitoso cifra API key, devuelve safe projection y status connected.
- Get/list son paginados y no contienen api key/ciphertext.
- Patch de name no toca credenciales.
- Cambio de URL o rotación de API key valida la combinación candidata antes de
  reemplazar el estado anterior.
- Fallo de rotación conserva URL, key y status anteriores.
- Delete aplica el mismo guard de referencias Project y borra metadata/secreto de
  forma atómica.

## 9. Dokploy SSRF y external adapter

- Rechazar URL con userinfo, query, fragment o path que duplique `/api`.
- Permitir HTTPS público/privado y HTTP sólo loopback/private.
- Rechazar metadata/link-local, multicast, unspecified y resolución DNS mixta que
  incluya un destino prohibido.
- Revalidar resolución en cada conexión para mitigar DNS rebinding.
- Rechazar redirects a otro origen o a IP prohibida.
- Enviar API key sólo mediante `x-api-key` y nunca en URL/log/error.
- Connection test invoca exactamente `/api/settings.health`.
- Discovery invoca `/api/application.search` con parámetros limitados.
- 401/403, 404, rate limit, 5xx, JSON inválido, body grande y timeout se
  normalizan.

## 10. Dokploy applications

- Cuenta inexistente o no conectada produce error estable sin llamar al provider.
- Discovery mapea `applicationId`, `appName`, display name, environment y status.
- Campos faltantes obligatorios descartan el elemento inválido sin filtrar el
  payload crudo; si toda la página es inválida, devolver error normalizado.
- Status desconocido o vacío mapea `unknown`.
- Name/environment filter y cursor quedan firmados; una página posterior no
  puede cambiar filtros/limit/server.
- Discovery exitoso actualiza application count/last synced.

## 11. REST, DTOs y seguridad

- BodyParser/decoder rechaza body vacío, JSON inválido, propiedades desconocidas,
  body grande y trailing data.
- Path UUID inválido y query/cursor inválidos producen 400 estable.
- CRUD devuelve exactamente los envelopes/schemas publicados.
- PAT y API key sólo aparecen en request DTOs write-only y nunca se serializan en
  responses, errores o request logs.
- Connection tests devuelven 200 con resultado normalizado para fallos del
  provider.
- Callbacks usan `Cache-Control: no-store` y Location elegido por backend.
- Rutas privadas requieren middleware administrador; mutaciones rechazan Origin
  inválido con 403.
- Callbacks no requieren cookie pero sí state válido.
- Delete referenced devuelve 409.
- Provider failures nunca se confunden con auth de Akritas.

## 12. OpenAPI y documentación

- Los endpoints implementados coinciden con operation IDs, requests, responses y
  schemas canónicos.
- Todas las mutaciones privadas documentan 401 y 403.
- `server_identifier` documenta el fingerprint derivado.
- `info.version` y el gate coinciden en `1.2.0`; el límite por defecto es 25.
- Ningún response schema adquiere propiedades sensibles.
- `mvp.md`, `spec.md` y `domain.md` quedan alineados con backlog, ADRs y OpenAPI.

## Expected failing tests before implementation

- No existen puertos, usecases ni adapters de integración.
- No existe configuración o conexión PostgreSQL.
- No existen migraciones ni Credential Store.
- No existen cipher, states Manifest, bindings o clients externos.
- No existen handlers/router para los endpoints.
- El OpenAPI y la documentación conservan las inconsistencias detectadas.

Tras la aprobación, los tests se crearán primero y deberán fallar por estas
ausencias antes de escribir implementación productiva.

## Validaciones finales

- `go test ./...`.
- `go test -race ./...`.
- `go test -tags=integration ./...` con PostgreSQL efímero.
- `go vet ./...`.
- `gofmt` y comprobación de worktree.
- `.harness/kernel/scripts/check-backend-architecture.sh`.
- `.harness/kernel/scripts/check-openapi.sh`.
- `.harness/kernel/scripts/check-security.sh`.
- Architecture review y security review del workflow.

## Acceptance criteria

- GitHub PAT y App Manifest producen GitHubAccount reutilizables y seguros.
- GitHub y Dokploy sólo persisten secretos cifrados en Credential Store.
- Credenciales inválidas nunca se convierten en integraciones utilizables.
- Repositorios y aplicaciones accesibles pueden descubrirse mediante el contrato
  OpenAPI y alimentar posteriormente a Project.
- Ningún secreto se devuelve, registra o almacena en dominio/Project.
- No se incorpora funcionalidad de H2 o posterior.

## Human approval required

Este archivo requiere aprobación humana explícita. No crear tests ni implementar
código antes de recibirla.
