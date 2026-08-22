# TDD Test Plan

## Scope

Definir mediante tests el contrato de la fundación Go y del package `internal/core/domain`. No se prueban HTTP, DTOs, puertos, usecases, adapters, repositorios ni migraciones porque no forman parte de esta tarea.

## Tests to add/update

### 1. Módulo y límites arquitectónicos

- El módulo declara `github.com/Unknowns24/akritas/backend`, Go 1.26 y compila con el entrypoint inerte.
- El package domain no contiene tags `json`/`gorm` ni importa HTTP, Chi, GORM, adapters, SDKs, filesystem u `os/exec`.
- Las entidades usan `uuid.UUID` para identidades internas, `time.Time` para instantes y `time.Duration` para la ventana de agrupación.

### 2. Enums

- Cada enum acepta exactamente los valores publicados por OpenAPI para integraciones, monitoring, severidad, Incident, Investigation, Evidence, Remediation, Validation y cambios de código.
- `Validate` rechaza strings vacíos o desconocidos con el sentinel del componente correspondiente.
- Los valores terminales de Incident se reconocen de manera consistente.

### 3. Constructores y value objects

- Cada constructor rechaza UUID cero, strings obligatorios vacíos, timestamps inconsistentes, cantidades negativas y referencias incompletas.
- Confidence acepta inclusivamente 0 y 1 y rechaza valores exteriores.
- Los constructores copian slices; modificar la entrada después de construir no altera la entidad.
- `GitHubAccount` mantiene separados `account_type` y `authentication_method`.
- Project requiere referencias GitHub/Dokploy y exactamente una MonitoringConfiguration válida.
- Ningún tipo admite o expone credenciales, tokens, password, TOTP o secretos.

### 4. MonitoringConfiguration y AutomationPolicy

- `DefaultMonitoringConfiguration` devuelve `enabled=false`, listas vacías, `30m`, contexto anterior 20 y posterior 20.
- Se aceptan regex válidas; se rechazan regex inválidas, patrones vacíos, más de 100 patrones o patrones mayores a 500 caracteres.
- `grouping_window` debe ser mayor a cero; ambos contextos deben estar en 0..1000.
- Los ignored patterns permanecen separados y disponibles para que detección les dé precedencia absoluta.
- Automation acepta los cuatro estados coherentes: todo deshabilitado, solo investigation, investigation+remediation y los tres habilitados.
- Se rechaza remediation sin investigation y PR sin investigation/remediation.

### 5. AdministratorSession

- Una sesión es activa antes de sus expiraciones y mientras no esté revocada.
- Es inactiva exactamente al alcanzar la expiración idle o absoluta, y después de revocarla.
- `Revoke` es idempotente y rechaza un instante anterior a la autenticación.
- La expiración idle no puede superar la absoluta.

### 6. LogEvent e Incident

- LogEvent requiere al menos una regla de detección, fingerprint, contenido sanitizado y contexto marcado como redactado.
- Un Incident nuevo comienza en `detected`, con una ocurrencia y first/last seen iguales.
- `CanGroup` exige mismo Project, mismo fingerprint, Incident no terminal y diferencia respecto de `last_seen_at` menor o igual a la ventana.
- El borde exacto de la ventana agrupa; un instante posterior no agrupa.
- `RecordOccurrence` incrementa el contador y actualiza `last_seen_at` sin modificar `first_seen_at`.
- Se rechazan ocurrencias anteriores a `last_seen_at`, de otro Project/fingerprint o sobre Incident terminal.
- Flujo válido: detected → investigating → publishing_issue → completed/requires_human, siempre con Issue.
- Flujo fixable válido: publishing_issue con Issue → remediating → completed/pull_request_created, siempre con PR.
- Remediation fallida termina como completed/remediation_failed sin PR.
- Fallos de investigation o publicación terminan en phase failed con su outcome correspondiente; retry vuelve a investigating y limpia el outcome.
- Se rechazan saltos de fase, remediation sin `fixable`, cierre sin Issue y cierre exitoso sin PR.

### 7. Investigation y Evidence

- Investigation sigue exclusivamente `pending → running → completed|failed`.
- Completar requiere root cause status, resolution status y confidence válidos; fallar conserva un mensaje público seguro.
- `finished_at` nunca antecede `started_at`.
- Evidence acepta exclusivamente tipos documentados, exige resumen y conserva siempre contenido sanitizado/redactado.
- Slices de hipótesis, archivos, commits y acciones recomendadas son defensivos.

### 8. Remediation, ValidationResult y referencias

- Remediation sigue `planned → in_progress → validated|failed`.
- Solo una remediación `validated` puede adjuntar PullRequestReference y pasar a `pull_request_created`.
- Una validación fallida impide estado validated y creación de PR.
- ValidationResult sigue `pending → running → passed|failed` y mantiene timestamps coherentes.
- CodeChange valida path, change type, patch sanitizado y redacted=true.
- GitHubIssueReference y PullRequestReference requieren número positivo, URL, repositorio, timestamp y, para PR, branch.

### 9. Errores enriquecidos y catálogo

- Todos los sentinels cumplen `DxAAABBBT`, usan scope `0`, capa `4`, el componente reservado y un tipo permitido.
- Los códigos son únicos y coinciden con `docs/errors/aaa-map.md`.
- `errors.Is` reconoce un sentinel por código; `Wrap` conserva `errors.Is/errors.As` y la causa mediante `Unwrap`.
- `Error()` no expone la causa interna ni datos sensibles y cada error contiene `message` y `user_message` seguros.

### 10. Validaciones finales

- Ejecutar `go test ./...`.
- Ejecutar `go test -race ./...`.
- Ejecutar `go vet ./...`.
- Ejecutar `gofmt` sobre los archivos Go y comprobar que no queden diferencias.
- Ejecutar `.harness/kernel/scripts/check-backend-architecture.sh`.
- Ejecutar `.harness/kernel/scripts/check-openapi.sh` sin modificar `docs/openapi.yaml`.
- Ejecutar `.harness/kernel/scripts/check-security.sh`.

## Expected failing tests before implementation

- No existe `go.mod`, entrypoint ni package domain.
- No existen entidades, enums, constructores, validaciones o transiciones.
- No existe `domain.Error` ni `docs/errors/aaa-map.md`.
- No existe la estructura hexagonal versionada.

Los tests se escribirán después de la aprobación de este plan y deberán fallar antes de incorporar la implementación correspondiente.

## Acceptance criteria covered

- Base Go 1.26 compilable y lista para ramas paralelas.
- Dominio completo del MVP, independiente de transporte y persistencia.
- Defaults e invariantes compatibles con ADRs y OpenAPI.
- Ciclos de estado y fronteras humanas verificables.
- Catálogo estable de errores de dominio.
- Ausencia estructural de secretos e imports de infraestructura.

## Open questions / human approval notes

- No quedan decisiones funcionales abiertas; se aplican las elecciones aprobadas en el plan de implementación.
- La estrategia para reabrir o asociar incidentes terminales que reaparecen permanece fuera del MVP; se crea un nuevo Incident.
- Se requiere aprobación humana explícita de este archivo antes de crear tests o implementar código.
