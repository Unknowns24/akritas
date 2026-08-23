# Implementation Brief

## Task

`AKR-68/69`: recovery con rotación password/TOTP y hardening de rate limits y
sesiones.

## Current project context

El backend usa arquitectura hexagonal, PostgreSQL/GORM, migraciones gormigrate,
Credential Store AES-GCM compartido, configuración Viper central y sesiones
opacas con idle/absolute expiry y `revoked_at`. Recovery existe sólo en el
OpenAPI canónico.

## Proposed approach

- Reutilizar `PendingEnrollment` para recovery sin crear otro modelo.
- Start recovery valida trabajo de credenciales de forma genérica y sólo
  reemplaza el enrollment pendiente.
- Verify recovery consume el enrollment y, dentro de una transacción, rota el
  hash, reemplaza el seed, revoca sesiones y crea la nueva sesión.
- Hacer login stale-safe condicionando el consumo TOTP al hash observado.
- Resolver y extender sesión mediante un único update condicional.
- Acotar el limiter fixed-window y mover attempts/window/max keys a Config.

## Architecture and persistence impact

Se agregan ports/usecases/handlers de recovery y operaciones cohesivas a los
repositorios existentes. No se requieren columnas ni tablas nuevas. GORM y SQL
permanecen dentro del adapter PostgreSQL.

## API/OpenAPI impact

Se implementan los dos endpoints ya documentados. No cambia OpenAPI ni su
versión 1.5.0.

## Error handling impact

Recovery usa `ErrInvalidCredentials` para fallos valid-shape de email,
bootstrap, enrollment o TOTP. Setup deja de serializar un error específico del
bootstrap. 429 conserva el error estable vigente.

## Risks

- Carrera login/recovery: CAS ligado al password hash y bulk revoke atómico.
- Carrera authenticate/revoke: update condicional con `RETURNING`.
- Saturación del limiter: cap fijo y fail-closed para keys nuevas.
- Estado parcial: una sola transacción para todos los cambios persistentes.
