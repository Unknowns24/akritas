# Implementation Summary

## Task

AKR-REST-CORS — CORS credentialed global sobre el router Chi.

## Implemented changes

- Se agregó `github.com/go-chi/cors` v1.2.2 como dependencia directa.
- `router.New` monta `cors.Handler` como middleware top-level, antes del
  registro de `/api/v1`.
- La lista de orígenes proviene exclusivamente de `Config.AllowedOrigins`.
- Se habilitaron credenciales y los métodos `GET`, `HEAD`, `POST`, `PUT`,
  `PATCH`, `DELETE` y `OPTIONS`.
- Se habilitaron `Accept`, `Content-Type` e `Idempotency-Key` para preflight.
- El preflight se cachea durante 300 segundos.
- `RequireAllowedOrigin` permanece activo para mutaciones autenticadas.

## Contract and persistence

No hubo cambios en OpenAPI, DTOs, endpoints, casos de uso, persistencia ni
migraciones.

## Tests added

- Headers CORS credentialed sobre el `401` de sesión ausente.
- Preflight global para JSON e `Idempotency-Key`.
- Rechazo de Origin no configurado.
- Compatibilidad con requests sin Origin.
- Reflexión exacta entre múltiples orígenes, sin wildcard.
- Preservación de la protección CSRF en mutaciones.

## TDD result

La fase roja confirmó ausencia de headers CORS y preflight `405`. Después de
montar el middleware oficial, la suite target y la regresión completa quedaron
verdes.
