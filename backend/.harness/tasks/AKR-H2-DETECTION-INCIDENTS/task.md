# AKR-H2-DETECTION-INCIDENTS

Implementar AKR-23 a AKR-34 como un único incremento backend H2 sobre la
implementación H1 vigente.

## Alcance aprobado

- ingesta incremental de logs desde el adapter Dokploy existente;
- checkpoint durable, protección contra replay y continuidad fail-closed;
- reconstrucción multilinea, redacción, reglas determinísticas y fingerprint;
- contexto acotado, LogEvent e Incident agrupado transaccionalmente;
- ejecución background con shutdown ordenado;
- lista, detalle y ocurrencias de Incidents según OpenAPI;
- migraciones explícitas y pruebas unitarias, REST e integración PostgreSQL.

## Fuera de alcance

QVAC, Investigation, análisis de causa raíz, GitHub Issues, Remediation,
modificación de código y Pull Requests no participan en H2.
