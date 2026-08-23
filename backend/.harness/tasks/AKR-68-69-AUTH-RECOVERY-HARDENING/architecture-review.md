# Architecture Review — AKR-68/69

## Veredicto

Aprobado para `backend_api`, sin hallazgos bloqueantes.

## Límites y dependencias

- Los usecases coordinan puertos de repositorio, crypto, limiter y transacción;
  no importan HTTP, GORM ni PostgreSQL.
- SQL, `RETURNING`, locks y operaciones compare-and-set permanecen en adapters
  PostgreSQL.
- REST implementa exclusivamente `POST /auth/recovery` y
  `POST /auth/recovery/verify` con los DTOs existentes en OpenAPI 1.5.0.
- Password, TOTP, Credential Store, errores, cookies y sesiones reutilizan el
  modelo H1; no existe un segundo subsistema de autenticación.

## Consistencia

- El reemplazo del slot pendiente se serializa y la confirmación usa consumo
  atómico exactamente una vez.
- Rotación, reemplazo cifrado, revocación global y nueva sesión comparten una
  transacción PostgreSQL.
- Login usa CAS por password hash observado y período TOTP; session refresh usa
  update condicional. Los ordenamientos concurrentes quedan linealizados por
  locks de fila y verificados contra PostgreSQL 17.

## Alcance

No se modificó comportamiento de Incidents, Investigation, QVAC, GitHub Issues,
Remediation ni Pull Requests. El único ajuste fuera de auth fue truncar a
segundos un fixture de paginación del test integrado, alineándolo con la
precisión del cursor para eliminar una falla preexistente revelada por el run.

`check-backend-architecture.sh` y `check-openapi.sh` pasan; OpenAPI permanece en
1.5.0 con 60 operaciones y 112 schemas.
