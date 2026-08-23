# AKR-REST-CORS - CORS credentialed en el router Chi

## Estado

in_progress

## Tipo de tarea

backend-api-feature

## Modo de proyecto

existing_project

## Contexto

El frontend de desarrollo corre en `http://localhost:3000` y consume el backend
en `http://localhost:8080`. La configuración ya contiene una allowlist exacta
en `AKRITAS_ALLOWED_ORIGINS`, pero actualmente sólo se usa para proteger
mutaciones contra CSRF. El router no emite headers CORS y el navegador oculta
incluso respuestas esperadas como el `401` de una sesión ausente.

## Objetivo

Configurar CORS credentialed como middleware global del router Chi usando
`github.com/go-chi/cors` y la allowlist runtime ya validada.

## Requerimiento funcional

- Permitir requests browser desde orígenes configurados exactamente.
- Permitir cookies mediante `Access-Control-Allow-Credentials: true`.
- Resolver preflight `OPTIONS` antes del matching de endpoints.
- Agregar headers CORS también a respuestas de error como `401`.
- No permitir wildcard ni reflejar orígenes fuera de la allowlist.
- Conservar `RequireAllowedOrigin` como protección CSRF de mutaciones.

## Criterios de aceptación

- `http://localhost:3000` puede leer el `401` de `GET /api/v1/auth/session`.
- Un preflight permitido responde con los métodos y headers aprobados.
- Orígenes no configurados no reciben autorización CORS.
- Requests sin `Origin`, como curl y comunicación server-to-server, conservan
  su comportamiento actual.
- Tests, race, vet y gates del profile `backend_api` pasan.

## Restricciones técnicas

- Middleware oficial `github.com/go-chi/cors` v1.2.2.
- Montaje top-level mediante `root.Use`, requerido por el middleware para que
  `OPTIONS` funcione correctamente.
- `AllowCredentials=true` exige orígenes exactos; se prohíbe `*`.
- Sin cambios de DTOs, casos de uso, persistencia o auth.

## Archivos o zonas probablemente afectadas

- `internal/adapter/rest/router/router.go`
- `internal/adapter/rest/router/router_test.go`
- `go.mod` y `go.sum`
- `docs/configuration.md`, si requiere aclaración operativa
- `.harness/tasks/AKR-REST-CORS/`

## Fuera de alcance

- Desactivar cookies `Secure` o relajar producción a HTTP.
- Permitir subredes, patrones arbitrarios o wildcard origins.
- Cambiar la semántica del `401` de sesión ausente.
- Modificar OpenAPI, endpoints o payloads.

## Instrucción para el harness

No implementar código hasta aprobación humana explícita de
`tdd-test-plan.md`.

## Aprobación humana

El usuario aprobó explícitamente el plan TDD en la conversación del 2026-08-23
con la respuesta “si apruebo”.
