# Implementation Summary — AKR-68/69

## Resultado

Se implementaron los dos endpoints de recovery definidos por OpenAPI 1.5.0 y
se endurecieron los entry points de autenticación y la autoridad de sesión del
servidor, preservando el único Administrator y el enrollment TOTP en dos pasos.

## Recovery

- Start valida bootstrap/email/password con error externo genérico, costo
  Argon2id acotado para identidades ausentes y buckets independientes por IP y
  cuenta normalizada.
- La password y seed nuevos permanecen pendientes y cifrados; las credenciales
  activas no cambian hasta confirmar el TOTP.
- Verify consume el enrollment exactamente una vez y, en una transacción,
  rota hash y seed, registra el período TOTP, revoca sesiones anteriores y crea
  una sesión nueva.
- La respuesta de provisioning se entrega con `Cache-Control: no-store`; nunca
  expone hash, ciphertext ni token de sesión.

## Rate limiting y sesiones

- El fixed-window en memoria tiene reloj inyectable, límite estricto de keys,
  limpieza de ventanas vencidas y fail-closed para keys nuevas al saturarse.
- Setup, setup verify, login, recovery y recovery verify usan instancias
  separadas y claves de peer directo más cuenta/enrollment según corresponde.
- La sesión se autentica con un único `UPDATE ... WHERE ... RETURNING` que
  verifica revocación e idle/absolute expiry y extiende el idle sin exceder el
  absoluto.
- Logout conserva el primer `revoked_at`; recovery usa revoke-all. El CAS TOTP
  de login queda ligado al hash observado para bloquear logins stale.

## Persistencia y configuración

No se agregaron migraciones: las tablas vigentes ya contienen hash, Credential
Store cifrado, expiraciones y `revoked_at`. Se agregaron settings Viper para
attempts, window y max keys, con defaults, límites de startup y documentación.
