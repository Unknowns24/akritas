# Task

## ID

`AKR-H1-INTEGRATIONS`

## Issues incluidas

- AKR-5 / PB-001 — Gestionar `GitHubAccount`.
- AKR-6 / PB-002 — Almacenar credenciales GitHub cifradas.
- AKR-7 / PB-003 — Validar autenticación GitHub.
- AKR-8 / PB-004 — Obtener y seleccionar repositorios accesibles.
- AKR-9 / PB-005 — Gestionar `DokployServer`.
- AKR-10 / PB-006 — Almacenar credenciales Dokploy cifradas.
- AKR-11 / PB-007 — Validar conectividad con Dokploy.
- AKR-12 / PB-008 — Obtener y seleccionar aplicaciones Dokploy.
- AKR-21 / PB-066 — Conectar GitHub mediante App Manifest.

## Objetivo

Implementar como un único incremento del Hito 1 las integraciones reutilizables
de GitHub y Dokploy, con credenciales cifradas fuera de `Project`, validación de
conectividad y discovery de los recursos que los futuros casos de uso de Project
podrán seleccionar.

## Profile y workflow

- Profile: `backend_api` (`.harness/kernel/profiles/go-hexagonal-api.yaml`).
- Workflow: `.harness/kernel/workflows/backend-api-feature.yaml`.
- Estado: `awaiting_tdd_approval`.

## Alcance

- CRUD seguro de GitHubAccount mediante PAT.
- GitHub App Manifest con registro, conversión e instalación verificada.
- Credential Store AES-256-GCM persistido en PostgreSQL.
- Prueba de conexión y discovery de repositorios GitHub.
- CRUD seguro de DokployServer.
- Prueba de conexión y discovery de aplicaciones Dokploy.
- Persistencia GORM y migraciones gormigrate versionadas.
- Puertos, usecases, adapters externos y REST de las operaciones publicadas en
  `docs/openapi.yaml`.
- Ajustes de OpenAPI y documentación necesarios para alinear el contrato con las
  policies de seguridad y el backlog P0.

## Fuera de alcance

- PB-009..PB-013: CRUD de Project, asociación persistente, configuración y
  activación de monitoring.
- PB-061..PB-063: setup, TOTP, login, sesión y middleware concretos de auth.
- Logs, Detection Engine, Incidents, QVAC, GitHub Issues, Remediation y Pull
  Requests.
- GitHub Enterprise, webhooks, OAuth App, rotación automática de master key,
  Vault/KMS, auto-merge o deploy.

## Dependencias externas al incremento

- PB-061..PB-063 debe proveer el middleware de administrador antes de habilitar
  los endpoints de integración en runtime.
- PB-010/PB-011 consumirá el discovery implementado aquí y completará el adapter
  de referencias desde Project para impedir el borrado de integraciones en uso.

## Artefactos de esta etapa

- `implementation-brief.md`.
- `tdd-test-plan.md`.

No se crearán tests ni implementación antes de la aprobación humana explícita
del plan TDD.
