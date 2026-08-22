# Security Review

## Summary

El endpoint es público por diseño (no hay sesión previa en el bootstrap). Protege el secreto TOTP (descifrado sólo en el punto de uso), el código enviado (nunca se distingue la causa del fallo) y el token de sesión (sólo su hash se persiste). Dos bugs reales de esta capa de seguridad se encontraron y corrigieron durante la propia verificación manual de esta revisión (ver Findings).

## Auth / permissions

- `POST /auth/setup/verify` no requiere sesión (`security: []` en el OpenAPI) — correcto, porque en el momento del bootstrap todavía no existe ninguna.
- El secreto TOTP se descifra (`CredentialStore.Decrypt`) únicamente dentro del usecase, en el momento de verificar — nunca antes, nunca se cachea (ADR-005).
- Rate limit de 5 intentos / 15 minutos por IP, ahora con **budget independiente** del de `/auth/setup` (bug encontrado y corregido en esta revisión — ver Findings), aplicado antes de tocar la DB o descifrar nada.
- La violación del índice único de `email` en `AdministratorRepository.Create` se traduce al mismo `409 ErrAdministratorAlreadyExists` que el chequeo temprano `ExistsActive`, cerrando la ventana de carrera entre dos verificaciones casi simultáneas sin degradar a un 500 genérico (ajuste explícitamente pedido por el usuario tras la primera versión del plan).
- `Administrator.Create` + `AdministratorSession.Save` corren en una única transacción GORM — no puede quedar un Administrator activado sin su sesión, ni viceversa.

## Input validation

- `enrollment_id` debe ser un UUID válido; `totp_code` debe matchear `^[0-9]{6}$` — validado en el handler antes de invocar el usecase.
- El usecase revalida en profundidad: enrollment inexistente, expirado (`IsExpired`), o código incorrecto (`TOTPVerifier.Verify`, RFC 6238 con tolerancia ±1 período) — todos colapsan al mismo `400 ErrInvalidTotpEnrollmentVerification`, sin distinguir la causa en la respuesta (ADR-008).

## Data exposure

- El token de sesión en claro sólo viaja en el `Set-Cookie` (`HttpOnly`, `Secure` según config, `SameSite=Lax`); nunca aparece en el body JSON — verificado por test (`TestVerifyAdministratorSetupHappyPath`) y manualmente.
- Sólo el SHA-256 hex del token se persiste en `administrator_sessions.token_hash` — confirmado manualmente comparando el valor de la cookie real contra la fila en la base (no coinciden, como se espera de un hash).
- El ciphertext del secreto TOTP (`encrypted_totp_secret`, ya cifrado desde `/auth/setup`) se copia tal cual a `administrators` — nunca se vuelve a exponer en claro; verificado manualmente inspeccionando la fila (60 bytes de ciphertext, no el secreto Base32 de 32 caracteres).
- El logger de GORM sigue en modo `Silent` (heredado de PB-061) — sin SQL con valores en logs.

## Error leakage

- `writeVerifyAdministratorSetupError` reutiliza el mismo patrón de PB-061: errores de dominio mapean por sufijo de código, el sentinel de rate limit mapea a 429 con `Retry-After`, y cualquier otro error cae al genérico `500` sin exponer la causa.
- Ninguna respuesta de error incluye `totp_code` del request — verificado por test y manualmente.

## Findings

Dos bugs de seguridad/usabilidad reales, encontrados y corregidos durante esta revisión (no quedan pendientes):

1. **Migración no idempotente** (`20260822_03`): no es una falla de confidencialidad/integridad de datos, pero un despliegue fresco fallaba directamente al arrancar. Corregido con un chequeo `HasColumn`.
2. **Rate limiter compartido entre `/auth/setup` y `/auth/setup/verify`**: un atacante que agota deliberadamente el budget de `verify` con códigos incorrectos también bloqueaba `setup` para el mismo IP (y viceversa) — una superficie de denegación de servicio más amplia de la prevista por ADR-008 ("límites independientes"). Corregido con instancias de `RateLimiter` separadas por endpoint en `cmd/main.go`.

## Result

pass
