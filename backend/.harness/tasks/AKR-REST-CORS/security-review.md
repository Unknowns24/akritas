# Security Review

## Summary

CORS quedó limitado a la allowlist runtime exacta y habilita cookies únicamente
para orígenes aprobados.

## Auth / permissions

CORS no autentica ni autoriza. La sesión, middleware administrador y
`RequireAllowedOrigin` conservan sus boundaries. Un `401` sigue siendo `401`;
ahora el browser autorizado puede inspeccionarlo.

## Input validation

Preflight restringe métodos y headers a un conjunto explícito. No se habilitó
wildcard de orígenes ni headers.

## Data exposure

`Access-Control-Allow-Origin` sólo refleja un valor que coincide con
`Config.AllowedOrigins`. Orígenes no configurados no reciben autorización CORS.

## Error leakage

Los statuses y envelopes existentes no cambiaron. CORS sólo añade headers de
transporte sobre responses existentes.

## Findings

Sin hallazgos. `AllowCredentials=true` se combina con orígenes exactos, nunca
con `*`. La defensa CSRF no fue reemplazada ni relajada.

## Result

pass
