# Security Review - AKR-H4-GITHUB-ISSUE

## Veredicto

PASS.

## Credenciales

- El publisher GitHub obtiene tokens exclusivamente via `accountToken` y Credential Store.
- No se agregaron variables ad hoc, `gh`, tokens directos ni logs de credenciales.
- El adapter nunca retorna payload completo de GitHub ni bodies de error del proveedor.

## Contenido publicado

- El builder separa Evidence observada de conclusiones QVAC.
- El body usa Evidence sanitizada y limites deterministas.
- La redaccion defensiva cubre PAT/App tokens, private keys, JWT/session, cookies, DSN, variables secretas y URLs TOTP.
- El marcador HTML contiene solo el UUID de Investigation, sin secretos.

## REST/timeline

- Incident detail expone IssueReference publica: number, url, repository, created_at.
- Timeline deriva eventos de registros persistidos y no copia Evidence sensible.
- Errores publicos de publicacion no incluyen causas internas ni respuesta GitHub.

## Riesgos residuales

Si GitHub acepta la Issue y PostgreSQL falla antes de persistir IssueReference, H4 falla explicitamente y no inventa exito. La reconciliacion por marcador queda para H6.
