# ADR-003 — Separar detección de logs de investigación con IA

## Estado

Accepted

## Contexto

Enviar cada línea de logs a un modelo local sería ineficiente, lento y propenso a ruido.
Los modelos pequeños son más confiables cuando reciben contexto ya filtrado y estructurado.

## Decisión

Akritas tendrá una etapa determinística de detección y agrupación antes de invocar QVAC.
QVAC sólo será invocado cuando exista un Incident Candidate con suficiente señal.

## Consecuencias

- menor costo computacional;
- menor latencia;
- contexto más limpio;
- mejor funcionamiento con modelos pequeños;
- necesidad de diseñar fingerprints y reglas de detección.
