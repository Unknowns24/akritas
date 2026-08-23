# Implementation Brief

## Task

`AKR-GITHUB-APP-MANIFEST-STATE`: corregir el handoff de registro mediante
GitHub App Manifest para que GitHub reciba y devuelva el nonce de correlación.

## Current project context

- `StartRegistration` genera 32 bytes con `crypto/rand` y los codifica mediante
  Base64 URL-safe sin padding.
- Persiste exclusivamente `sha256(state)` en `ConversionStateDigest`, con
  expiración de una hora.
- `CompleteManifest` exige longitudes válidas, vuelve a calcular SHA-256 y llama
  a `ConsumeConversionState`, que aplica consumo único y expiración.
- `form_action` selecciona correctamente el endpoint personal u organización,
  pero actualmente omite `state`; esa omisión explica el callback productivo
  con sólo `code` y el error `ErrManifestStateInvalid`.
- El test de flujo existente usa `started.State` directamente y por eso no
  reproduce el handoff real del navegador a través de `form_action`.
- El baseline focalizado `go test ./internal/usecase/githubapp` pasa antes de
  introducir los tests de regresión.

## Proposed approach

1. Agregar primero tests de usecase que parseen `form_action`, validen host,
   path, cardinalidad del query y equivalencia exacta de `state` para ownership
   personal y organization, incluyendo caracteres que requieren escaping.
2. Ajustar el test end-to-end con fakes para tomar el `state` de `form_action`,
   simular el callback `code + state` y completar la conversión.
3. Mantener pruebas explícitas para digest-only, replay, ausencia de `state` y
   expiración.
4. Construir la URL base existente y agregar `state` con `net/url` mediante
   `Query`, `Set`, `Encode` y `String`, sin helpers nuevos.
5. Aclarar el contrato OpenAPI y aplicar el incremento patch exigido por la
   policy SemVer, alineando el gate que valida la versión canónica.

## Architecture impact

El cambio permanece en `internal/usecase/githubapp/start_registration.go`, donde
ya vive la orquestación del handoff provider-specific. No se agregan puertos,
adapters, dependencias ni decisiones arquitectónicas.

## API/OpenAPI impact

No cambia el shape de la respuesta. Cambia el valor funcional de `form_action`
para cumplir el contrato previsto: será una URL lista para el POST y contendrá
el `state`. Se conserva el campo separado `state` por compatibilidad. La
documentación indicará que el consumidor sólo debe publicar el campo `manifest`
y no reconstruir query parameters. Al ser una corrección compatible del
contrato, se propone `info.version: 1.5.1`.

## Data/persistence impact

Ningún cambio de schema o migración. `ConversionStateDigest` conserva exactamente
32 bytes SHA-256; el nonce público no se persiste en claro.

## Error handling impact

No se agregan ni modifican errores. `CompleteManifest` conserva validación de
longitud, lookup por digest, expiración y consumo único.

## Test strategy

Tests unitarios/usecase con `registrationStoreFake` y `appGatewayFake`; no se
realizan requests contra GitHub ni se requieren secretos reales.

## Risks

- Doble encoding o concatenación frágil: mitigado parseando la URL en tests y
  usando `net/url` en producción.
- Breaking change al remover `state`: evitado conservando el campo.
- Regresión de seguridad: cubierta verificando digest-only, callback inválido,
  expiración y replay.

## Files likely to change

- `internal/usecase/githubapp/start_registration.go`
- `internal/usecase/githubapp/flow_test.go`
- `docs/openapi.yaml`
- `.harness/kernel/scripts/check-openapi.sh`
- artefactos de `.harness/tasks/AKR-GITHUB-APP-MANIFEST-STATE/`

## Human gate

No crear tests ni modificar implementación hasta recibir aprobación explícita
de `tdd-test-plan.md`.
