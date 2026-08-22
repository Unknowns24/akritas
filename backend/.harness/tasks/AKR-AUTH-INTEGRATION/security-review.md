# Security Review

## Authentication and secrets

- Bootstrap token se compara en tiempo constante y nunca se persiste.
- Passwords usan Argon2id versionado; parámetros embebidos se validan y acotan
  antes de reservar memoria o derivar la clave.
- TOTP acepta RFC 6238 con tolerancia ±1 y el período consumido sólo avanza por
  compare-and-set dentro de la transacción que crea la sesión.
- Seeds TOTP se cifran con AES-256-GCM y AAD del owner; mover el owner vuelve a
  cifrar. Password hashes, token hashes, ciphertext y seeds no viven en dominio.
- Tokens de sesión son opacos, 32 bytes aleatorios y sólo SHA-256 llega a DB.

## HTTP boundary

- Fallos de credenciales y verificación son genéricos.
- Respuestas sensibles usan `Cache-Control: no-store`; cookies son HttpOnly,
  Secure, SameSite=Lax y tienen expiración absoluta.
- Session e integraciones exigen middleware de sesión. Toda mutación autenticada
  exige un Origin exacto allowlisted; missing, wildcard y mismatch devuelven 403.
- Rate limits son independientes para setup, verify y login; login separa IP y
  cuenta normalizada. Buckets expirados se limpian.
- El tipo `R` devuelve 429 y `Retry-After` sin revelar qué bucket se agotó.

## Transaction safety

- Argon2, generación de secretos y TOTP se ejecutan antes de abrir transacciones.
- Setup, verify y login revierten metadata, sesión y Credential Store juntos.
- No se realizan llamadas externas dentro de los nuevos boundaries.

## Findings

- El rate limiter no sobrevive reinicios ni coordina múltiples réplicas; queda
  fuera de alcance por decisión explícita.
- Recovery continúa sin implementar y no debe considerarse una vía disponible.

## Result

pass
