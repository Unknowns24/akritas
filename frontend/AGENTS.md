# AGENTS.md — Harness Engineering

Este proyecto usa Harness Engineering para guiar el trabajo de agentes de IA.

El harness vive en:

```txt
.harness/kernel/
```

Este archivo es la entrada raíz para agentes. No debe duplicar todas las reglas técnicas del sistema. La fuente de verdad para reglas específicas está en los `profiles`, `policies` y `workflows` del kernel.

---

## Principio de normalización

Los agentes no deben hardcodear listas de políticas.

El flujo correcto es:

```txt
Agente
→ declara o infiere profile activo
→ lee .harness/kernel/harness.yaml
→ resuelve profiles.<profile>
→ carga required_policies del profile
→ aplica workflow correspondiente
```

Si una policy cambia, se actualiza el profile correspondiente, no cada agente.

---

## Resolución obligatoria del harness

Antes de modificar código, cualquier agente debe:

1. Leer `.harness/kernel/harness.yaml`.
2. Identificar el profile activo según la tarea.
3. Abrir el archivo indicado por `profiles.<profile>`.
4. Cargar todas las policies listadas en `required_policies`.
5. Cargar fuentes declaradas por el profile, incluyendo, cuando estén presentes:

- `openapi_locations`
- `design_locations`
- `architecture_locations`

Las fuentes de arquitectura son contexto autoritativo del proyecto y deben
revisarse antes de realizar cambios arquitectónicos o de dominio.

6. Aplicar el workflow correspondiente.
7. Respetar los gates humanos definidos por el harness.

Profiles esperados:

```txt
backend_api     → .harness/kernel/profiles/go-hexagonal-api.yaml
backend_service → .harness/kernel/profiles/go-hexagonal-service.yaml
frontend        → .harness/kernel/profiles/next-feature-based.yaml
mobile_flutter  → .harness/kernel/profiles/flutter-hexagonal.yaml
```

---

## Selección de profile

### Backend Go API

Usar profile:

```txt
backend_api
```

Aplicar el workflow:

```txt
.harness/kernel/workflows/backend-api-feature.yaml
```

Usar este profile para tareas sobre servicios backend Go que exponen una API REST, incluyendo arquitectura hexagonal, handlers, casos de uso, puertos, adapters, repositorios, OpenAPI, seguridad backend y tests.

### Backend Go Service

Usar profile:

```txt
backend_service
```

Aplicar el workflow:

```txt
.harness/kernel/workflows/backend-service-feature.yaml
```

Usar este profile para tareas sobre servicios Go que no necesariamente exponen una API REST, incluyendo arquitectura hexagonal, casos de uso, puertos, adapters, repositorios, procesamiento interno, seguridad y tests.

No asumir la existencia de HTTP, handlers, routers u OpenAPI salvo que la especificación del proyecto o la tarea lo requiera.

---

### Frontend Next.js

Usar profile:

```txt
frontend
```

Aplicar el workflow:

```txt
.harness/kernel/workflows/frontend-feature.yaml
```

Usar este profile para tareas sobre dashboard/frontend Next.js, App Router, TypeScript, CSS Modules, integración con OpenAPI, UI web, runtime envs y tests frontend.

Si el profile declara `design_locations`, el agente debe buscar el contrato visual en esas rutas y respetarlo.

---

### Mobile Flutter

Usar profile:

```txt
mobile_flutter
```

Aplicar el workflow:

```txt
.harness/kernel/workflows/mobile-flutter-feature.yaml
```

Usar este profile para tareas sobre apps Flutter, arquitectura hexagonal feature-based, Riverpod, GoRouter, Dio, integración mobile con backend, permisos de dispositivo, UI mobile y tests Flutter.

Si el profile declara `design_locations`, el agente debe buscar el contrato visual en esas rutas y respetarlo.

---

### Fullstack web

Aplicar el workflow:

```txt
.harness/kernel/workflows/fullstack-web-feature.yaml
```

Este workflow coordina backend + frontend web.

Orden esperado:

```txt
backend/OpenAPI primero
→ frontend consume contrato actualizado
```

---

### Fullstack mobile

Aplicar el workflow:

```txt
.harness/kernel/workflows/fullstack-mobile-feature.yaml
```

Este workflow coordina backend + Flutter.

Orden esperado:

```txt
backend/OpenAPI primero
→ Flutter consume contrato actualizado
```

---

## Regla TDD obligatoria

No se debe implementar código antes de generar y aprobar el plan de tests TDD cuando la tarea sea manejada por el harness.

Flujo mínimo:

```txt
implementation-brief.md
→ tdd-test-plan.md
→ aprobación humana
→ implementación
→ validaciones
→ reviews
→ final-summary.md
```

Si el usuario pide explícitamente implementar sin pasar por harness/TDD, el agente debe aclarar que eso rompe el workflow estándar del harness y pedir confirmación.

---

## Contrato backend

Si la tarea toca endpoints, DTOs, payloads, errores o integración con backend:

1. Buscar OpenAPI en las rutas declaradas por el profile activo.
2. No inventar endpoints.
3. No inventar campos.
4. No consumir endpoints fuera del scope de la app.
5. Actualizar OpenAPI cuando el contrato cambia.
6. Validar OpenAPI con el script correspondiente si existe.

La policy común está en:

```txt
.harness/kernel/policies/openapi.md
```

---

## Diseño y UX

Si la tarea toca UI web o mobile:

1. Buscar diseño en `design_locations` del profile activo.
2. Respetar el contrato visual existente.
3. No rediseñar pantallas existentes salvo pedido explícito.
4. No hardcodear colores, spacing o componentes si existen tokens/componentes compartidos.
5. Mantener loading, empty, error y retry en pantallas conectadas al backend.

Las reglas específicas deben venir desde las policies del profile activo.

---

## Seguridad

Toda tarea debe respetar:

```txt
.harness/kernel/policies/security.md
```

Además, si el profile activo incluye una policy de seguridad específica, también debe aplicarse.

Reglas generales:

- No loguear tokens, cookies ni contraseñas.
- No exponer secretos.
- No duplicar clientes HTTP con credenciales propias.
- No enviar passwords en texto plano salvo que el contrato lo exija explícitamente.
- No relajar validaciones de auth/RBAC sin aprobación explícita.

---

## Validaciones

Antes de cerrar una tarea, ejecutar los scripts que correspondan al profile activo y al workflow usado.

Scripts habituales:

```txt
.harness/kernel/scripts/check-backend-architecture.sh
.harness/kernel/scripts/check-frontend-architecture.sh
.harness/kernel/scripts/check-flutter-architecture.sh
.harness/kernel/scripts/check-openapi.sh
.harness/kernel/scripts/check-security.sh
```

Si un script no existe o no puede ejecutarse en el entorno actual, dejarlo informado en el summary final.

---

## Artefactos esperados por tarea

Cuando la tarea use el workflow del harness, generar o actualizar dentro de `.harness/tasks/<TASK-ID>/`:

```txt
task.md
implementation-brief.md
tdd-test-plan.md
implementation-summary.md
architecture-review.md
security-review.md
final-summary.md
```

Si la tarea es pequeña y el usuario no pidió usar formalmente el harness, responder con criterio práctico, pero sin violar las policies del profile detectado.

---

## Memoria del proyecto

Usar:

```txt
.harness/memory/project-summary.md
.harness/memory/decisions.md
```

para recuperar decisiones estables del proyecto.

No usar memory para logs temporales o información generada de una ejecución puntual. Para eso usar:

```txt
.harness/runtime/
```

---

## No hacer

- No duplicar listas de policies dentro de agentes.
- No implementar antes de aprobación humana del plan TDD cuando el workflow lo requiere.
- No inventar endpoints, DTOs ni campos.
- No ignorar OpenAPI si existe.
- No ignorar design files si el profile declara `design_locations`.
- No crear clientes HTTP paralelos.
- No filtrar DTOs a UI si la arquitectura lo prohíbe.
- No romper dirección de dependencias definida por el profile.
- No meter lógica de negocio en capas de presentación.
- No consumir endpoints administrativos desde apps públicas o de usuario final.
- No dejar mocks como fuente productiva.
- No cerrar una tarea sin summary de validaciones ejecutadas o pendientes.

---

## Regla final

La fuente de verdad del harness es:

```txt
.harness/kernel/harness.yaml
→ profiles
→ required_policies
→ workflows
→ scripts
```

Este archivo solo define cómo resolver y aplicar el harness. Las reglas técnicas detalladas deben vivir en policies.

<!-- BEGIN:nextjs-agent-rules -->

# This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` (resolved from this file's directory; in monorepos the `next` package may not be visible from the repo root) before writing any code. Heed deprecation notices.

This block is written and re-added by `next dev` — verify at `node_modules/next/dist/server/lib/generate-agent-files.js`. Removing it from a diff only re-creates the uncommitted change; committing it with your work keeps the tree clean.

<!-- END:nextjs-agent-rules -->
