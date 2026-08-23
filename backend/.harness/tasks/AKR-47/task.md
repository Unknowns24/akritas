# AKR-47 - Redaccion segura y contenido auditable

## Objetivo

Fortalecer `internal/service/evidencesafety` y `internal/service/issuecontent` para que ningun secreto frecuente pueda publicarse en una GitHub Issue, manteniendo contenido auditable, deterministico y separado entre Evidence observada y conclusiones QVAC.

## Alcance

- Redaccion defensiva de secretos en logs, Evidence, titulo y body final de Issue.
- Tests table-driven para formatos de secreto requeridos.
- Tests del builder con secretos en todos los campos solicitados.
- Integridad PostgreSQL para impedir referencias inconsistentes entre Incident, Investigation e IssueReference.
- Documentacion/harness summaries ajustados a cobertura realmente implementada y ejecutada.

## Fuera de alcance

- Cambios en contrato HTTP publico u OpenAPI, salvo incompatibilidad real.
- Cambios en publicacion remota de GitHub fuera del contenido que se envia.
- Remediation, Pull Requests o reconciliacion remota post-fallo.

