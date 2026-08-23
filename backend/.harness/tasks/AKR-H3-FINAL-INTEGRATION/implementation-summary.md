# Implementation Summary — AKR-H3-FINAL-INTEGRATION

## Resultado

H3 quedó integrado sobre H2 real. El flujo productivo es `Incident → LogEvents → Evidence → QVAC local → GitHub read-only tools → structured result → PostgreSQL`; no hay `IncidentReader`, deny-all ni runner stub en producción y H3 no ejecuta `StartIssuePublication`.

## Compilación y catálogos

- Se reemplazaron los dos `Catalog()` PostgreSQL superpuestos por una única función cerrada y alcanzable.
- El catálogo conserva `ErrIntegrationPersistence`, `ErrProjectPersistence`, `ErrInvestigationPersistence`, `ErrOperationPersistence`, `ErrEvidencePersistence`, `ErrIncidentPersistence` y `ErrMonitoringPersistence`.
- `ErrIncidentNotFound` pertenece sólo al catálogo Incident; se eliminó su segunda catalogación desde Investigation.
- El test agregado global verifica formato, código exacto y nombre únicos, y una única fila por sentinel/código en `docs/errors/aaa-map.md` para dominio/usecase, PostgreSQL, REST y GitHub externo.

## Arquitectura H2 → H3

- Se eliminó el port legado `IncidentReader` y su `Exists`. H3 usa capacidades pequeñas reales: `IncidentGetter`, `IncidentLogEventLister` e `IncidentWorkflowStore`; `IncidentStore` compone las lecturas REST H2.
- El repositorio Incident PostgreSQL implementa `Get/List/ListLogEvents/Lock/Update`. No se agregó `Exists` ni un adapter puente sin propósito.
- `StartIncidentInvestigation` usa una transacción, lock de Incident, revalidación de idempotencia/actividad, transición de Incident y creación de Investigation+Operation; sólo despacha luego del commit.
- Start/run/failure/recovery persisten pares Investigation/Operation de manera transaccional. `unknown` válido completa; fallo técnico marca Investigation, Operation e Incident como failed con mensajes públicos.
- Recovery previo a HTTP reencola `pending+queued` y falla trabajo iniciado/inconsistente conservando Evidence, bajo la decisión aprobada de una instancia MVP.
- El éxito H3 deja Incident en `investigating`; no inicia semánticamente H4.

## Evidence y QVAC

- `InvestigationContextAssembler` resuelve `Investigation → Incident → Project → GitHubAccount/repository` y los LogEvents persistidos H2.
- Persiste metadata no sensible, un `log_excerpt` JSON por evento real y `stack_trace` sólo ante la regla H2 real. Incluye message, timestamp, severity, rules, source y before/after; no fabrica file/commit/diff.
- Corpus: 17 Evidence iniciales como máximo, 25/128 KiB persistidos; la Evidence inicial enviada al prompt tiene techo 24 KiB y se reduce según `ctx_size=16384` y la reserva para sistema/tools/output.
- Evidence inicial se persiste antes de QVAC; Evidence de tools se persiste incluso si luego falla el resultado. Lecturas repetidas y persistencia se deduplican.
- QVAC recibe `InvestigationRunContext` completo y no consulta PostgreSQL. Sólo puede citar UUIDs de Evidence realmente mostrada; Investigation persiste `evidence_ids` y `evidence_count` conserva el corpus total.
- Toda DATA no confiable se redacta y viaja como JSON entre `UNTRUSTED_DATA_BEGIN/END`. Se cubren bearer/PAT/installation token/private key/access key/JWT, variables secretas, cookies/sesiones y URLs con userinfo.

## Tools y resultado

- `RepositoryInspector` expone exclusivamente `SearchCode`, `ReadFile`, `ListRecentCommits`, `ReadCommit`, `ReadDiff` con tipos application-level.
- El adapter QVAC crea exactamente las cinco tools, fijadas al account/owner/name/default branch resueltos desde Project. Argumentos de owner/repo/ref extra son rechazados; repo A/B se prueba explícitamente.
- Límites: 8 rondas, 24 calls, 8 KiB JSON válido por payload, 16 KiB acumulados. Unknown falla cerrado.
- El resultado usa schema/decoder estricto: required fields, no extras/trailing JSON, enums, confidence `[0,1]`, tamaños/arrays y Evidence UUIDs únicos/citables.
- Hallazgos se persisten como `code_location`, `commit` o `diff`; no existe ninguna tool mutativa ni invocación de Remediation.

## Persistencia y API

- `20260823_04_link_investigation_history`: FKs RESTRICT Investigation→Incident y Evidence→Investigation, más índice único parcial para `pending|running`.
- `20260823_05_add_investigation_evidence_ids`: JSONB no nulo con default `[]`.
- No se modificaron migraciones aplicadas. La referencia histórica dentro de una migración aplicada que decía que H2 estaba pendiente queda superseded por `_04`, no se usa como estado vigente.
- OpenAPI 1.6.0 agrega `created_at`, `evidence_ids` y hace `started_at` opcional para pending, sin rutas nuevas.
- Se preservaron `projectusecase.NewWithMonitoring(...)`, monitoring runner y composición REST Incident/Investigation/Operation/Evidence.
