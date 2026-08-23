# Architecture review — AKR-H1-PROJECTS

## Veredicto

Conforme con el profile `backend_api`, las ADRs vigentes y el workflow
`backend-api-feature`.

## Revisión

- El core sólo conoce dominio y ports planos; no importa REST, GORM ni clientes
  externos.
- Los tags GORM del dominio son metadata pasiva permitida por ADR-012.
- `ProjectUseCase` concentra invariantes y orquesta ports estrechos; handlers y
  repositorios no contienen lógica de negocio.
- PostgreSQL, GitHub y Dokploy permanecen detrás de adapters existentes.
- Router, sesión, CSRF, response envelopes, errores y Uker se reutilizan desde
  la rama destino.
- Los DTOs están agrupados por feature y los mappers separados por dirección.
- La migración es incremental y no usa `AutoMigrate`.
- El validador OpenAPI sólo actualizó el pin del contrato autoritativo de 1.3.0
  a 1.4.0; profiles, policies, bases y dependencias no cambiaron.

No se detectaron dependencias invertidas ni duplicación de infraestructura de
`feat/project-handling`.
