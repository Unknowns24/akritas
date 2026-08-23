# AKR-H1-PROJECTS

Integrar la historia de `feat/project-handling` sobre `feat/backend-milestone-1`
sin reemplazar infraestructura vigente y completar PB-009 a PB-013.

## Alcance aprobado

- CRUD autenticado de Project, incluido `DELETE /projects/{project_id}`.
- Asociación verificada con GitHubRepository y DokployApplication.
- MonitoringConfiguration completa y activación/desactivación segura.
- Persistencia PostgreSQL versionada, paginación Uker y contrato OpenAPI 1.4.0.

## Reglas de producto aprobadas

- nombre único sin distinguir mayúsculas;
- DokployApplication exclusiva por Project;
- asociaciones y borrado requieren monitoring desactivado;
- snapshots siempre resueltos contra proveedores reales;
- `default_branch` debe coincidir con GitHub;
- todo PUT habilitado revalida y deja el estado en `starting`.
