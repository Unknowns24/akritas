# Akritas — Incident Lifecycle

## 1. Observing

Akritas obtiene periódicamente nuevos logs de cada Project habilitado.

Los logs brutos no son considerados incidentes.

## 2. Candidate Detection

El Detection Engine identifica eventos potencialmente relevantes.

Un evento debe contener suficiente señal como para justificar análisis posterior.

## 3. Fingerprinting

El evento recibe un fingerprint estable derivado de sus características principales.

Objetivo:

- reconocer repeticiones del mismo error;
- evitar una Issue por cada ocurrencia;
- agrupar errores equivalentes aunque contengan valores variables.

Ejemplo conceptual:

```text
panic: nil pointer at users/service.go:81 user=123
panic: nil pointer at users/service.go:81 user=929
```

Deben poder pertenecer al mismo incidente.

## 4. Incident Creation / Update

Si no existe un incidente abierto compatible con el fingerprint, se crea uno.

Si ya existe, se actualizan:

- `last_seen_at`;
- `occurrence_count`;
- evidencia relevante adicional.

## 5. Investigation

QVAC recibe un contexto acotado del incidente y puede solicitar herramientas para investigar.

La investigación puede iterar entre:

```text
Reason
  ↓
Request tool
  ↓
Receive evidence
  ↓
Reason again
```

Hasta producir una conclusión o alcanzar límites de seguridad/recursos.

## 6. Classification

Toda investigación debe producir dos clasificaciones independientes:

### Root Cause

- identified
- suspected
- unknown

### Resolution

- fixable
- requires_human

## 7. Issue Creation

La Issue se crea siempre para toda Investigation completada.

La Issue es la salida obligatoria de la investigacion y el registro humano/auditable del incidente. Akritas persiste una `GitHubIssueReference` ligada a `Incident` e `Investigation`; si la misma Investigation ya tiene referencia, el workflow reutiliza ese resultado y no republica. Si GitHub crea la Issue pero falla la persistencia local, el workflow falla de forma explicita y el marcador HTML con el UUID de Investigation queda preparado para reconciliacion futura.

## 8. Optional Remediation

Solo si `resolution_status = fixable`. En H4, un incidente fixable queda esperando en `publishing_issue` despues de crear la Issue; H5 reutiliza esa fase para iniciar remediation.

Flujo:

```text
Issue created
   ↓
Create branch
   ↓
Generate changes
   ↓
Run validation
   ↓
Validation successful?
   ├── No → record failure / update issue
   └── Yes
         ↓
      Create PR
```

## 9. Human Boundary

El MVP se detiene luego de crear la Pull Request.

Un humano debe decidir:

- revisar cambios;
- aprobar;
- mergear;
- desplegar.

## 10. Repeated Incidents

Si un error previamente resuelto reaparece, Akritas debe poder registrar una nueva ocurrencia o reabrir/asociar el incidente según política futura.

La estrategia exacta queda fuera del MVP y puede definirse posteriormente.
