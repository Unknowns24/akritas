# ADR-002 — Todo incidente investigado crea una GitHub Issue

## Estado

Accepted

## Contexto

Una investigación puede terminar con una solución automática o puede requerir intervención humana.

Si la creación de una Issue dependiera de que Akritas pudiera solucionar el problema, los incidentes no solucionables quedarían fuera del workflow normal del equipo.

## Decisión

Todo Incident que alcance la etapa de investigación debe producir una GitHub Issue.
La Issue es el registro canónico y auditable del incidente.
La capacidad de remediación se evalúa después.

## Consecuencias

- todos los incidentes quedan visibles para humanos;
- las PRs pueden referenciar siempre una Issue;
- Akritas aporta valor aunque no pueda corregir un problema;
- las decisiones del agente quedan auditables.
