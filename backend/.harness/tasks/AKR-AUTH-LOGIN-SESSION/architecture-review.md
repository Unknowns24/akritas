# Architecture Review

## Summary

Cierra el ciclo de autenticación sobre la base ya establecida en PB-061/062. Introduce el primer middleware real del proyecto y el primer par lectura+escritura por request autenticada (resolución de sesión + idle TTL deslizante), sin romper la dirección de dependencias hexagonal ni el patrón de wiring ya validado.

## Layering

- `internal/core/domain` y `internal/core/ports/{in,out}` siguen sin importar `internal/adapter`, GORM, Chi, `net/http` ni `os/exec` (`check-backend-architecture.sh` pasa).
- `internal/usecase/auth` sólo depende de `internal/core/domain` y `internal/core/ports/{in,out}`.
- `internal/adapter/rest/middleware` depende únicamente de `ports/in` (`AuthenticateSessionUseCase`) y `internal/adapter/rest/response` — mismo patrón de dependencia que los handlers, no una excepción nueva.
- `response.SessionCookieName` se extrajo de `handler/auth` (antes privada ahí) para que middleware y handlers compartan el nombre sin que uno dependa del otro — evita el acoplamiento lateral que hubiera significado que middleware importe el paquete handler.

## Modularity / SRP

- Un archivo por operación pública se mantiene en usecases, repositorios y handlers nuevos.
- `writeSessionError` (mapeo de error compartido entre `GetCurrentSession` y `Logout`) vive junto al primero de los dos handlers, mismo patrón ya usado en PB-061/062 para `writeStartAdministratorSetupError`/`writeVerifyAdministratorSetupError`/`writeLoginError` — consistente, no una excepción nueva.
- `TOTPVerifier.Verify` cambió de firma (agrega `period int64`) — no es una violación de SRP (sigue siendo una única responsabilidad: verificar un código), es una extensión de contrato explícitamente aprobada, con su único call site anterior (PB-062) actualizado mecánicamente.

## OpenAPI consistency

- No se modificó `docs/openapi.yaml`; `check-openapi.sh` confirma 59 operaciones y 112 schemas sin cambios.
- `LoginRequest`/`SessionResponse`/`Session`/`Administrator` se usan/reusan sin diverger del contrato. La cookie de logout expira con `Max-Age: 0` en vez de `-1`; ambos son formas válidas de expirar una cookie ya emitida y el schema del header sólo pide "Expired HttpOnly session cookie" sin fijar el mecanismo exacto — verificado manualmente que el navegador (curl) la trata como expirada.

## Findings

- (Heredado de PB-061/062, mismo patrón reutilizado aquí para `ErrLoginRateLimited`) los handlers REST importan paquetes concretos de `internal/usecase/auth` en vez de sólo `ports/in`, para reconocer sentinels de rate limit que no tienen representación en los in ports. No bloqueante, permitido por `project-structure.md`, ya señalado dos veces — si PB-065 formaliza rate limiting, es el momento natural de resolverlo de una vez para los tres casos.
- `GET /auth/session` hace tres operaciones de DB por request (`FindByTokenHash` + `UpdateIdleExpiry` en el middleware, `FindByID` en el usecase del handler) en vez de una sola consulta combinada. Aceptable para el volumen de tráfico de un sistema single-admin; se señala como posible optimización futura, no como defecto.

## Result

pass
