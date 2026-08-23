# AKR-68/69 - Recovery y hardening de autenticación

## Estado

complete

## Tipo de tarea

backend-api-feature

## Modo de proyecto

existing_project

## Contexto

H1 ya implementa setup, confirmación TOTP, login, sesión opaca y logout. El
OpenAPI 1.5.0 ya define `POST /auth/recovery` y
`POST /auth/recovery/verify`, pero el router no los implementa. Los límites de
autenticación son fixed-window en memoria, con parámetros hardcodeados y sin
cota estricta de buckets.

## Objetivo

Implementar PB-064 y PB-065: recovery en dos pasos que rota password/TOTP y
revoca sesiones previas, más rate limiting acotado/configurable y validación de
sesión resistente a carreras con revocación.

## Criterios de aceptación

- Recovery sólo se autoriza con email del Administrator y bootstrap token.
- La confirmación rota password y TOTP, consume el período confirmado, revoca
  sesiones previas y crea una sesión nueva en una transacción PostgreSQL.
- Login iniciado con credenciales anteriores no puede persistir una sesión tras
  completar recovery.
- Sesiones expiradas, revocadas o aleatorias se rechazan en servidor.
- Setup, verify, login y ambos endpoints de recovery tienen budgets aislados.
- El limiter tiene memoria acotada y configuración Viper documentada.
- Errores externos sensibles son genéricos y ningún secreto se serializa o
  registra.

## Restricciones técnicas

- Profile `backend_api`; workflow `backend-api-feature`.
- Reutilizar Argon2id, TOTP, Credential Store, Transactor, GORM/gormigrate y
  patrones REST actuales.
- No cambiar paths ni DTOs del OpenAPI 1.5.0.
- No introducir Redis, AutoMigrate ni un segundo modelo de autenticación.

## Fuera de alcance

Incidents, Investigation, QVAC, GitHub Issues, Remediation y Pull Requests.

## Aprobación humana

El usuario aprobó explícitamente el plan TDD e implementación el 2026-08-23.
