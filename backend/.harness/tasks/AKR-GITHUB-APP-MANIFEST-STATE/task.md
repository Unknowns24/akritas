# Task

## ID

`AKR-GITHUB-APP-MANIFEST-STATE`

## Estado

`complete`

## Tipo y modo

- Profile: `backend_api` (`.harness/kernel/profiles/go-hexagonal-api.yaml`).
- Workflow: `backend-api-feature` (`.harness/kernel/workflows/backend-api-feature.yaml`).
- Modo: `existing_project`.
- Base Git verificada: `HEAD`, `main` y `origin/main` apuntan a `a9df482` al iniciar la tarea.

## Contexto

`StartRegistration` genera un `state` criptográficamente aleatorio, persiste
únicamente `sha256(state)` y devuelve el nonce en la respuesta, pero construye
`form_action` sin incluirlo. Como el navegador publica el Manifest contra esa
URL incompleta, GitHub no puede devolver el `state` en el callback y
`CompleteManifest` rechaza correctamente el callback.

## Objetivo

Entregar una `form_action` autosuficiente que contenga exactamente el `state`
generado, tanto para cuentas personales como organizaciones, manteniendo las
garantías de expiración, digest y consumo único del flujo Manifest.

## Criterios de aceptación

- La URL personal apunta a `https://github.com/settings/apps/new` e incluye un
  único query parameter `state` con el nonce generado.
- La URL de organización apunta al endpoint de la organización, con escaping
  correcto del slug y del `state`.
- El callback puede recuperar ese mismo `state` desde `form_action` y completar
  la conversión mediante el digest persistido.
- Callback sin `state`, expirado o repetido continúa fallando.
- El valor persistido sigue siendo únicamente `sha256(state)`.
- El campo `state` de la respuesta se conserva por compatibilidad.
- OpenAPI documenta que `form_action` ya incorpora el protocolo de correlación.

## Restricciones

- No relajar `CompleteManifest` ni cambiar el mecanismo de GitHub App Manifest.
- No persistir el nonce en claro, agregar cookies/sesiones, estado en memoria o
  infraestructura nueva.
- No introducir cambios de frontend ni funcionalidades ajenas al Hito actual.

## Fuera de alcance

- GitHub Enterprise, OAuth Apps, webhooks, cambios de permisos o instalación.
- Refactors arquitectónicos no necesarios para corregir el handoff.

## Gate humano

Primero se generan `implementation-brief.md` y `tdd-test-plan.md`. No se agregan
tests ni se modifica código productivo antes de la aprobación humana explícita.

Aprobación recibida el 2026-08-23: `TDD Aprobado`.
