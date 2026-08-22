# Security Review

## Summary

El contrato aplica autenticación por defecto, limita explícitamente las rutas públicas y mantiene los secretos fuera de los recursos persistentes y responses normales.

## Auth / permissions

- `cookieAuth` es el requisito global; health/readiness, enrollment/login/recovery y callbacks GitHub son las únicas excepciones.
- Setup y recovery requieren el bootstrap token como campo write-only; TOTP utiliza enrollment de una sola visualización.
- La sesión se documenta opaca y server-side, con cookie HttpOnly, Secure, SameSite=Lax y revocación en recovery.
- Los callbacks GitHub requieren `state` de un solo uso y responden mediante redirect 303.

## Input validation

- Password 12–128, TOTP de seis dígitos, UUID, URI, rangos de contexto, durations y enums se expresan en schemas.
- Automation y QVAC incluyen invariantes condicionales en JSON Schema 2020-12.
- Comandos manuales con efectos externos exigen `Idempotency-Key`.
- QVAC restringe endpoints a loopback/redes privadas y prohíbe redirects hacia hosts públicos.

## Data exposure

- Passwords, bootstrap token, PAT y credenciales Dokploy/QVAC son write-only.
- PAT, private key, webhook secret, token de instalación, seed TOTP y token de sesión no forman parte de los DTOs de recursos.
- Logs, evidence, patches y resultados se documentan sanitizados/redactados.

## Error leakage

- Las operaciones reutilizan un error normalizado con código de dominio, mensaje para usuario y `request_id`.
- Los tests de conexión y fallos de proveedores no devuelven payloads o errores crudos.
- Login y recovery documentan errores genéricos y rate limiting.

## Findings

No se encontraron secretos hardcodeados ni hallazgos bloqueantes. La implementación futura deberá verificar criptografía, Argon2id, anti-replay TOTP, expiraciones, rate limiting, validación de Origin/state y redacción en runtime.

## Result

pass
