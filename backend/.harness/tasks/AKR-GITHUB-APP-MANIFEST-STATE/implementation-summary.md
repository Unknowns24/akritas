# Implementation Summary

## Task

`AKR-GITHUB-APP-MANIFEST-STATE`

## Cause confirmed

`StartRegistration` devolvía el nonce en `state` y persistía correctamente su
SHA-256, pero construía `form_action` sin query parameters. El frontend hacía el
POST del Manifest a esa URL, por lo que GitHub nunca recibía el nonce que debía
devolver al callback. `CompleteManifest` rechazaba correctamente el callback
sin `state` mediante `ErrManifestStateInvalid`.

## TDD execution

Después de la aprobación humana se agregaron primero tests de personal,
organization, escaping y handoff completo. El estado RED se observó porque las
URLs existentes no contenían `state` y el flujo no podía extraerlo desde
`form_action`.

La implementación mínima construye la URL base existente, usa `net/url` para
establecer `state` como query parameter y devuelve `formURL.String()`.

## Files changed

- `internal/usecase/githubapp/start_registration.go`: incorpora el nonce a la
  URL oficial personal u organización mediante `net/url`.
- `internal/usecase/githubapp/start_registration_test.go`: cubre endpoints,
  escaping, query exacto, digest SHA-256, ausencia y expiración.
- `internal/usecase/githubapp/flow_test.go`: toma el callback state desde
  `form_action`, verifica conversión y replay, y conserva invariantes del
  Manifest.
- `docs/openapi.yaml`: documenta la URL autosuficiente y conserva `state` por
  compatibilidad; versión patch `1.5.1`.
- `.harness/kernel/scripts/check-openapi.sh`: alinea el gate con `1.5.1`.
- artefactos de `.harness/tasks/AKR-GITHUB-APP-MANIFEST-STATE/` e índice.

## Security invariants preserved

- `randomState` continúa usando 32 bytes de `crypto/rand` con Base64 URL-safe.
- Persistencia continúa recibiendo solamente `sha256(state)`.
- La expiración continúa siendo de una hora.
- `CompleteManifest` conserva validación de longitud y lookup por digest.
- `ConsumeConversionState` conserva expiración y consumo único.
- Callback sin state, vencido o repetido continúa rechazado.
- El Manifest no incorpora secretos ni modifica permisos/webhooks.

## Contract behavior

Antes, el consumidor recibía una `form_action` incompleta y debía inferir que
tenía que agregar el nonce. Ahora puede publicar únicamente el hidden input
`manifest` contra la URL recibida. El campo separado `state` permanece para no
introducir un breaking change, pero no debe usarse para reconstruir la URL.
