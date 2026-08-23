# Harness Tasks

| ID | Estado | Tipo | Descripción |
| --- | --- | --- | --- |
| AKR-GITHUB-APP-MANIFEST-STATE | complete | backend-api-feature | Incluir el state en form_action del handoff GitHub App Manifest |
| AKR-68-69-AUTH-RECOVERY-HARDENING | complete | backend-api-feature | Recovery password/TOTP, rate limiting acotado y sesiones robustas |
| AKR-AUTH-INTEGRATION | complete | backend-api-feature | Integración limpia de autenticación sobre el backend milestone |
| AKR-H1-INTEGRATIONS | complete | backend-api-feature | GitHub/Dokploy, Credential Store y discovery del Hito 1 |
| AKR-OPENAPI-MVP | complete | backend-api-feature | OpenAPI v1, autenticacion TOTP y contratos del MVP |
| AKR-BACKEND-FOUNDATION | complete | backend-api-feature | Estructura hexagonal y dominio completo del MVP |
| AKR-AUTH-BOOTSTRAP | complete | backend-api-feature | Bootstrap del único Administrator: setup-status y setup con TOTP |
| AKR-AUTH-TOTP-VERIFY | complete | backend-api-feature | Enrollment y verificación TOTP: activa al Administrator y abre sesión |
| AKR-AUTH-LOGIN-SESSION | complete | backend-api-feature | Login, sesión opaca (idle TTL deslizante) y logout |
| AKR-REST-CHI | complete | backend-api-feature | Migración modular del adaptador REST a Chi |
| AKR-H1-PROJECTS | complete | backend-api-feature | Project CRUD, snapshots verificados y MonitoringConfiguration |
| AKR-INVESTIGATION-LIFECYCLE | complete | backend-api-feature | Investigation + Operation (infra genérica async), frontera con H2 |
| AKR-INVESTIGATION-EVIDENCE | complete | backend-api-feature | Evidence real (deployment_metadata) ensamblada en el pipeline async |
| AKR-QVAC-INFERENCE | complete | backend-service-feature | PB-028/032/033/034/035: runner QVAC local + resultado estructurado |
| AKR-QVAC-TOOL-LOOP | complete | backend-service-feature | PB-029: loop de tool calling allowlisted read-only |
| AKR-GITHUB-REPO-TOOLS | complete | backend-service-feature | PB-030/031: search_code/read_file/commits/diffs vía adapter GitHub |
| AKR-H2-DETECTION-INCIDENTS | complete | backend-api-feature | Ingesta incremental, Detection Engine determinístico, LogEvents e Incidents |
| AKR-REST-CORS | complete | backend-api-feature | CORS credentialed global mediante middleware oficial de Chi |
