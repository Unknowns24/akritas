# Security Review

## Summary

Ambos endpoints son públicos por diseño (`security: []`, acorde a ADR-008: el bootstrap no puede requerir autenticación previa). El flujo protege el secreto TOTP y el bootstrap token en cada punto: comparación constant-time, cifrado en reposo, logging silenciado, y respuestas de error genéricas verificadas manualmente contra el servidor real.

## Auth / permissions

- `POST /auth/setup` valida el bootstrap token con `crypto/subtle.ConstantTimeCompare` antes de cualquier otra operación.
- `config.Load()` exige `AKRITAS_BOOTSTRAP_TOKEN` no vacío al arrancar (falla rápido). Como defensa en profundidad adicional, `Verify` también rechaza explícitamente un token configurado vacío o un candidato vacío, para no depender únicamente del wiring de `main.go`.
- Rate limit de 5 intentos / 15 minutos por IP (confirmado con el usuario), aplicado antes de tocar bootstrap token, DB o generación de secretos — verificado manualmente (6to intento → 429 + `Retry-After: 60`).
- El rate limiter usa `r.RemoteAddr` (no `X-Forwarded-For`), evitando bypass por spoofing de headers; la contrapartida documentada es que, detrás de un proxy compartido sin configuración de proxies confiables, múltiples clientes comparten el mismo bucket — aceptable para el MVP (`PB-065` es rate limiting avanzado, fuera de alcance).
- Origin validation (ADR-008) no aplica a este endpoint: no hay sesión ni cookie todavía en el bootstrap; queda pendiente para los endpoints autenticados de PB-062 en adelante.

## Input validation

- Validación de shape en el handler (email, `display_name` 1-100, password 12-128, bootstrap_token 32-512) antes de invocar el usecase.
- `domain.PendingEnrollment.Validate()` es una segunda línea de defensa a nivel de dominio (UUID no-nil, email parseable, campos no vacíos, `expires_at` posterior a `created_at`).
- Ninguna validación distingue en el mensaje de error cuál de email/password/bootstrap_token falló específicamente (ADR-008: respuestas genéricas).

## Data exposure

- El secreto TOTP en claro (`otpauth_uri`, `manual_entry_key`) sólo existe en la respuesta HTTP de un único `POST /auth/setup`; lo persistido es exclusivamente el resultado de `CredentialStore.Encrypt` (AES-256-GCM, nonce aleatorio por operación).
- `password_hash` es Argon2id (`m=19456,t=2,p=1`, salt aleatorio de 16 bytes); nunca se persiste ni se loguea la password en claro.
- El logger de GORM se configura en modo `Silent` explícitamente — evita que el driver vuelque SQL con valores (incluido `password_hash`/`encrypted_totp_secret`) a stdout/stderr.
- Verificado manualmente: los logs del proceso no contienen el bootstrap token, la password ni el master key tras ejercer los cuatro caminos (éxito, token inválido, administrator existente, rate limit).

## Error leakage

- `response.WriteInternalError` devuelve siempre el mismo mensaje genérico (`0x100001I`); nunca incluye `err.Error()` de la causa real (fallos de DB, hashing, etc.), verificado por `TestStartAdministratorSetupUsecaseErrors/unexpected_error`.
- `domain.Error.Error()` ya excluye la causa envuelta (comportamiento heredado de `AKR-BACKEND-FOUNDATION`, cubierto por `TestDomainErrorWrapPreservesIdentityAndCause`).
- Las respuestas de error nunca incluyen `bootstrap_token` ni `password` del request, verificado por test y manualmente.

## Findings

Ningún hallazgo bloqueante. Se aplicó una mejora de defensa en profundidad durante esta revisión: `BootstrapTokenVerifier.Verify` ahora rechaza explícitamente un token configurado vacío o un candidato vacío, en vez de depender únicamente de que `config.Load()` nunca deje pasar un token vacío al arrancar (cubierto por `TestBootstrapTokenVerifierRejectsEmptyConfiguredToken`).

## Result

pass
