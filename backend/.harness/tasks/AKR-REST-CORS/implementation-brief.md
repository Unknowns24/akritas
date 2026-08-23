# Implementation Brief

## Task

AKR-REST-CORS — CORS credentialed global sobre el router Chi.

## Current project context

`config.Config.AllowedOrigins` ya contiene orígenes exactos y validados. El
router utiliza esos valores en `RequireAllowedOrigin` únicamente para rechazar
mutaciones cross-origin no confiables; no existe un middleware que genere
`Access-Control-Allow-*`. Por eso un browser desde `http://localhost:3000` no
puede inspeccionar la respuesta `401` de `/api/v1/auth/session`.

## Proposed approach

- Incorporar `github.com/go-chi/cors` v1.2.2.
- Montar `cors.Handler` en el router Chi raíz antes de registrar `/api/v1`.
- Usar `Config.AllowedOrigins` sin patrones adicionales.
- Permitir credenciales, métodos `GET`, `HEAD`, `POST`, `PUT`, `PATCH`,
  `DELETE`, `OPTIONS` y headers browser necesarios (`Accept`, `Content-Type`,
  `Idempotency-Key`).
- Mantener `RequireAllowedOrigin` en sus scopes actuales como segunda defensa
  para mutaciones autenticadas.
- Configurar un `MaxAge` corto para desarrollo y operación predecible.

## Architecture impact

El cambio permanece en el adapter REST. No introduce dependencias en core,
usecases ni adapters de persistencia. La allowlist continúa resuelta por
`config` y entregada al router desde bootstrap.

## API/OpenAPI impact

No cambian paths, métodos, payloads, estados ni schemas. CORS es comportamiento
de transporte browser y no requiere incremento de versión del contrato.

## Data/persistence impact

Ninguno. No hay migraciones.

## Error handling impact

Los errores existentes conservan su status y envelope. Para un Origin permitido,
el middleware añadirá headers CORS incluso sobre `400`, `401`, `403`, `404`,
`409`, `429` y `500`, permitiendo que el browser los inspeccione.

## Test strategy

Tests de router con `httptest` para requests simples, respuestas `401`,
preflight permitido, Origin rechazado y requests sin Origin; luego regresión
completa y gates del harness.

## Risks

- Combinar wildcard con cookies filtraría responses a orígenes no confiables.
- Montar CORS dentro de un subrouter rompería preflight `OPTIONS`.
- Reemplazar `RequireAllowedOrigin` reduciría la protección CSRF.
- Una lista insuficiente de headers permitidos bloquearía requests JSON o
  comandos idempotentes legítimos.

## Files likely to change

`internal/adapter/rest/router/router.go`, sus tests, `go.mod`, `go.sum`,
documentación runtime y artefactos de la tarea.
