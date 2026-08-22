# ADR-001 — QVAC como único motor de inferencia

## Estado

Accepted

## Contexto

Akritas procesa información extremadamente sensible: logs de producción, código fuente, stack traces, arquitectura interna y potencialmente identificadores o secretos accidentales.

Además, el proyecto compite en el QVAC Track, que exige inferencia local.

## Decisión

Toda inferencia de IA del MVP se ejecutará localmente mediante QVAC.

No se utilizarán APIs cloud de IA para análisis, clasificación, generación de Issues ni generación de correcciones.

## Consecuencias

### Positivas

- privacidad;
- menor exposición de información sensible;
- operación sin dependencia de un proveedor externo de IA;
- alineación directa con el track.

### Negativas

- modelos más pequeños;
- límites de contexto y memoria;
- mayor necesidad de tool calling y outputs estructurados;
- necesidad de controlar estrictamente el contexto enviado al modelo.
