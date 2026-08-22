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

La Issue se crea siempre.

La Issue es la salida obligatoria de la investigación y el registro humano/auditable del incidente.

## 8. Optional Remediation

Sólo si `resolution_status = fixable`.

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
