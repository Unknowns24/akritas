# CLAUDE.md — Harness Engineering Bridge

Usá `AGENTS.md` como entrada principal.

La fuente canónica del harness está en:

```text
.harness/kernel/
```

Para tareas con harness:

1. Leer `.harness/tasks/<TASK-ID>/task.md`.
2. Leer el workflow correspondiente en `.harness/kernel/workflows/`.
3. Generar `implementation-brief.md` y `tdd-test-plan.md`.
4. No implementar código hasta aprobación humana explícita.
5. Luego implementar, validar, revisar arquitectura, revisar seguridad y actualizar memoria.

No inventes endpoints, campos, permisos, errores, estructura de carpetas ni patrones visuales si el proyecto ya tiene una convención existente.

## Profiles

No dupliques listas de políticas en agentes o instrucciones.

Para cualquier tarea:

1. Leer `.harness/kernel/harness.yaml`.
2. Seleccionar el profile activo: `backend`, `frontend` o `mobile_flutter`.
3. Abrir el archivo indicado por `profiles.<profile>`.
4. Aplicar todas las políticas listadas en `required_policies`.
5. Usar el workflow correspondiente.

Para tareas mobile Flutter, usar profile `mobile_flutter` y workflow:

```text
.harness/kernel/workflows/mobile-flutter-feature.yaml
```

No implementes código Flutter antes de aprobación explícita del plan TDD.
No rediseñes pantallas existentes salvo pedido explícito.
No llames ApiClient/Dio desde presentation.
