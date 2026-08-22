# Akritas — Demo Story

## Objetivo de la demo

Mostrar que Akritas no es un chatbot de logs sino un ingeniero autónomo de incident response que trabaja localmente.

## Escenario sugerido

Aplicación simple desplegada en Dokploy con un bug reproducible introducido recientemente.

Ejemplo:

- endpoint `GET /users/:id`;
- para un usuario inexistente el código dereferencia un valor nil/null;
- la aplicación devuelve 500 y produce stack trace;
- el bug fue introducido por un commit reciente.

## Secuencia

### 0. Acceso seguro

En una instalación vacía, mostrar una sola vez:

```text
bootstrap autorizado
      ↓
QR TOTP
      ↓
código confirmado
      ↓
Administrator activo
```

En ejecuciones posteriores, comenzar con login mediante email, password y TOTP.
La demostración no debe mostrar el bootstrap token, seed TOTP ni session cookie.

### 1. Estado inicial

Mostrar Akritas con:

- Project activo;
- GitHub conectado;
- Dokploy conectado;
- estado `Monitoring`.

### 2. Provocar incidente

Realizar requests que disparen el bug.

Mostrar que la aplicación falla sin interactuar manualmente con Akritas.

### 3. Detection

Akritas detecta múltiples errores equivalentes y crea un único Incident.

Mostrar:

- first seen;
- occurrences;
- stack trace resumido;
- estado `Investigating`.

### 4. Investigación local

Mostrar timeline de tools:

```text
✓ inspected incident logs
✓ extracted stack trace
✓ searched repository
✓ read users/service.go
✓ inspected recent commits
✓ identified probable regression
```

Mensaje clave:

> None of this production data is being sent to a cloud AI provider.

### 5. Issue

Akritas crea automáticamente una Issue.

La Issue debe explicar:

- qué pasó;
- evidencia;
- causa raíz;
- confidence;
- archivo/línea relevante.

### 6. Fix

Akritas clasifica el incidente como `fixable`.

Mostrar:

```text
Generating remediation...
✓ branch created
✓ regression test added
✓ fix generated
✓ tests passed
✓ pull request created
```

### 7. Pull Request

Abrir la PR y mostrar:

- cambios pequeños y entendibles;
- test de regresión;
- explicación;
- vínculo a la Issue.

## Mensaje final sugerido

Un flujo que normalmente exige que un ingeniero abra logs, busque stack traces, entienda el repositorio, revise cambios recientes y prepare un fix fue realizado automáticamente por Akritas, con la inferencia ejecutándose localmente.
