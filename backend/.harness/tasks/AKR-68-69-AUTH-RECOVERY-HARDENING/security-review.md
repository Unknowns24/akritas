# Security Review — AKR-68/69

## Veredicto

Aprobado sin hallazgos bloqueantes.

## Credenciales y enumeración

- Passwords sólo llegan al Argon2id vigente y no se persisten en plaintext.
- Seeds TOTP pendientes y activos usan el Credential Store AES-GCM existente;
  la rotación reemplaza el único owner activo.
- Login y recovery realizan trabajo dummy para identidades ausentes y usan el
  mismo error externo para causas sensibles valid-shape.
- Bootstrap inválido de setup ya no serializa texto específico; se conserva la
  excepción contractual de setup-state/409.
- No se agregaron logs de payloads, bootstrap, TOTP, cookies o tokens.

## Sesiones, replay y carreras

- El token opaco permanece únicamente en cookie `HttpOnly`, `Secure`,
  `SameSite=Lax`, `Path=/`, con `Expires` y `Max-Age` centralizados.
- PostgreSQL rechaza sesiones random, revocadas o vencidas sin depender del
  vencimiento del browser.
- Confirmación consume el período TOTP, el enrollment se consume una vez y la
  rotación invalida hash, seed y todas las sesiones previas atómicamente.
- Refresh/revoke y login/recovery tienen compare-and-set/locks explícitos; no
  dependen de timing de aplicación.

## Rate limiting

El limiter fixed-window es determinístico, acotado a una cantidad configurable
de buckets y falla cerrado al saturarse. Usa sólo `RemoteAddr`, no headers
forwarded spoofeables, y es compatible con la topología MVP de un proceso.

`check-security.sh` y la inspección final de patrones sensibles pasan.
