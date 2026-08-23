# TDD Test Plan

## Scope

Definir el contrato corregido de `StartRegistration.form_action` y comprobar el
handoff completo hasta `CompleteManifest`, sin alterar generación, persistencia,
expiración ni consumo del `state`.

## Tests to add/update

### 1. Form action personal

Agregar un test focalizado de `StartRegistration` para `owner_type=personal`:

- inyectar un `state` determinístico de longitud válida que contenga caracteres
  con significado en query strings;
- parsear `FormAction` con `net/url`;
- comprobar `scheme=https`, `host=github.com` y path
  `/settings/apps/new`;
- comprobar que existe exactamente un parámetro, llamado `state`, cuyo valor
  decodificado coincide byte por byte con el generado;
- comprobar que la serialización aplica escaping y no introduce parámetros
  adicionales.

### 2. Form action de organización

Agregar un caso para `owner_type=organization`:

- verificar el ejemplo `organization=Unknowns24` contra el endpoint
  `/organizations/Unknowns24/settings/apps/new`;
- agregar cobertura con un valor que requiera path escaping para impedir que el
  slug altere segmentos de la URL;
- reutilizar un `state` que requiera query escaping y comprobar su round-trip
  exacto;
- exigir un único query parameter `state`.

### 3. Flujo completo sin GitHub real

Modificar `TestManifestFlowUsesSeparateOneTimeStatesAndVerifiedInstallation`:

- ejecutar `StartRegistration`;
- extraer el `state` únicamente desde `started.FormAction`, como lo recibiría
  GitHub;
- verificar que coincide con `started.State` (campo conservado por
  compatibilidad);
- invocar `CompleteManifest` con `code + state` extraído;
- comprobar que `ConsumeConversionState` encuentra la registration por el
  digest y que la conversión termina en estado `converted`;
- repetir el callback y exigir `ErrManifestStateConflict`.

### 4. Invariantes de seguridad

Agregar/ajustar assertions que comprueben:

- `ConversionStateDigest == sha256(started.State)` y longitud 32;
- el valor crudo no es el valor persistido;
- `CompleteManifest` sin `state` devuelve `ErrManifestStateInvalid` sin llamar
  al exchange de GitHub;
- un state correcto pero vencido devuelve `ErrManifestStateConflict` y tampoco
  llama al exchange;
- el manifest mantiene permisos existentes, App privada y webhooks inactivos,
  sin campos secretos nuevos.

El fake de gateway incorporará, sólo si hace falta, un contador de exchanges
para probar que los callbacks inválidos no alcanzan el provider.

### 5. Contrato OpenAPI

Actualizar la descripción de la operación y de `form_action` para dejar
explícito que:

- la URL ya contiene el `state` necesario para el handoff;
- el consumidor hace POST incluyendo sólo el input hidden `manifest`;
- no debe agregar ni reconstruir parámetros del protocolo;
- el campo separado `state` se mantiene por compatibilidad.

Aplicar el incremento patch `1.5.0 -> 1.5.1` y alinear el gate OpenAPI.

## Expected failing tests before implementation

- Los tests personal y organization fallarán porque `FormAction` no tiene query
  parameter `state`.
- El test de flujo fallará al extraer un state vacío de `FormAction`, y
  `CompleteManifest` devolverá `ErrManifestStateInvalid`.
- Las pruebas existentes de digest, replay y manifest deben seguir pasando; las
  pruebas nuevas de ausencia y expiración deben pasar sin relajar código.

Después de observar el fallo esperado se modificará la implementación mínima.

## Acceptance criteria covered

- `form_action` personal y organization contiene exactamente el nonce generado.
- Path y query mantienen escaping correcto.
- El callback válido consume la registration por SHA-256.
- Ausencia, expiración y replay siguen rechazados.
- No se persiste el nonce en claro ni se agregan secretos al manifest.
- `state` permanece en la respuesta sin trasladar responsabilidad al frontend.

## Final validations

- `go test ./...`
- `go vet ./...`
- `gofmt` sobre archivos Go modificados
- `.harness/kernel/scripts/check-backend-architecture.sh`
- `.harness/kernel/scripts/check-openapi.sh`
- `.harness/kernel/scripts/check-security.sh`
- review arquitectónica y de seguridad según el workflow

## Open questions / human approval notes

No hay decisiones de producto abiertas. Plan aprobado explícitamente por el
usuario el 2026-08-23 mediante `TDD Aprobado`.
