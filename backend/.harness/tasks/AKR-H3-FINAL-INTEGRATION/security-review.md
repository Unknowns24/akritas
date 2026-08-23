# Security Review — AKR-H3-FINAL-INTEGRATION

## Veredicto

PASS

## Prompt injection

- Logs, stacks, código, comentarios, commits, diffs y tool results son DATA no confiable.
- Prompt inicial y outputs usan JSON dentro de `UNTRUSTED_DATA_BEGIN/END`.
- El system prompt ordena ignorar instrucciones/exfiltración embebidas en DATA.
- QVAC sólo puede citar Evidence efectivamente mostrada.

## Secretos

- El assembler usa únicamente metadata Project no sensible y LogEvents H2 ya sanitizados, con redacción defensiva adicional.
- La redacción cubre Authorization bearer, GitHub tokens, access keys, private keys, JWT/session tokens, nombres de variables secretas, cookies y URLs con credenciales.
- Credential Store, configuración de proceso y entorno host nunca forman parte de `InvestigationRunContext`.
- No se registran errores crudos en Investigation/Operation; se persisten mensaje público y código estable cuando existe.

## GitHub

- Sólo cinco tools read-only.
- Schemas rechazan campos desconocidos; owner/repository/default branch no son argumentos del modelo.
- Paths/SHA se validan y los adapters mantienen path safety.
- Unknown tool falla inmediatamente; no hay shell ni operaciones write/issue/branch/commit/push/PR.

## Límites

- 24 KiB máximo adaptativo para Evidence inicial.
- 8 KiB por envelope JSON tool y 16 KiB acumulados.
- 8 rondas, 24 calls, 25 Evidence/128 KiB persistidos.
- Overflow o salida malformada falla de forma predecible.

## Gate

`.harness/kernel/scripts/check-security.sh`: PASS.
