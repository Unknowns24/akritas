# TDD Test Plan

## Scope

Agregar CORS credentialed global al router Chi usando exclusivamente la
allowlist exacta de `Config.AllowedOrigins`, sin cambiar autenticación, CSRF ni
contrato OpenAPI.

## Tests to add/update

1. `GET /api/v1/auth/session` desde `Origin: http://localhost:3000`, sin cookie,
   conserva `401` y agrega `Access-Control-Allow-Origin` con ese origen exacto,
   `Access-Control-Allow-Credentials: true` y `Vary: Origin`.
2. Un preflight `OPTIONS` permitido para una mutación JSON responde exitosamente
   antes del router, publica métodos autorizados y acepta `Content-Type`.
3. El preflight de comandos permite `Idempotency-Key`.
4. Un Origin que no pertenece a `AllowedOrigins` no recibe
   `Access-Control-Allow-Origin` ni autorización de credenciales.
5. Un request sin `Origin` conserva status/body existentes y no recibe headers
   CORS, manteniendo clientes no-browser.
6. La configuración con múltiples orígenes refleja únicamente el Origin exacto
   que coincide; nunca devuelve `*` con credenciales.
7. Una mutación autenticada con Origin permitido continúa atravesando
   `RequireAllowedOrigin`; una no permitida sigue fallando cerrada.
8. La construcción fail-closed del router ante `AllowedOrigins` vacía permanece
   sin cambios.

## Expected failing tests before implementation

- Las respuestas actuales no incluyen `Access-Control-Allow-Origin` ni
  `Access-Control-Allow-Credentials`.
- `OPTIONS` no está resuelto globalmente y cae en 404/405 según el path.
- No existe la dependencia `github.com/go-chi/cors`.

## Acceptance criteria covered

Lectura browser del `401`, preflight JSON/idempotencia, cookies credentialed,
allowlist exacta, ausencia de wildcard, compatibilidad no-browser y preservación
de CSRF/fail-closed.

## TDD sequence

1. Agregar tests CORS al router y confirmar que fallan contra el comportamiento
   actual.
2. Agregar `github.com/go-chi/cors` v1.2.2.
3. Montar/configurar el middleware top-level mínimo para ponerlos en verde.
4. Ejecutar regresión y validar manualmente con Origin
   `http://localhost:3000` sobre el backend local.

## Final validations

- `gofmt` check.
- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `.harness/kernel/scripts/check-backend-architecture.sh`
- `.harness/kernel/scripts/check-openapi.sh`
- `.harness/kernel/scripts/check-security.sh`
- `git diff --check`

## Open questions / human approval notes

Sin preguntas abiertas. El usuario aprobó explícitamente este plan en la
conversación del 2026-08-23 con la respuesta “si apruebo”.
