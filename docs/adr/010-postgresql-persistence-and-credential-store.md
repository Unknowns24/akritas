# ADR-010 — PostgreSQL como persistencia y Credential Store cifrado compartido

## Estado

Accepted

## Contexto

Akritas necesita persistir el dominio del Control Plane, estados operativos y
credenciales de integraciones. La arquitectura ya exige GORM dentro del adapter
de base de datos, migraciones explícitas mediante gormigrate y una master key
externa para cifrar secretos, pero no había seleccionado un motor relacional.

El Credential Store podría usar una base separada, por ejemplo SQLite, o vivir en
el mismo motor que el resto de la aplicación. Una segunda base aumentaría la
complejidad de transacciones, backups, restauración y despliegue sin crear una
frontera criptográfica: la confidencialidad de los secretos depende de
`AKRITAS_MASTER_KEY`, que permanece fuera de toda persistencia.

## Decisión

PostgreSQL será el motor relacional canónico del backend Akritas.

- El adapter vivirá en `internal/adapter/db/postgres/`.
- GORM se utilizará únicamente dentro del adapter de persistencia.
- La evolución de schema se realizará con
  `github.com/go-gormigrate/gormigrate/v2` y migraciones ordenadas e inmutables.
- No se utilizará `AutoMigrate` global como estrategia de schema.
- La configuración de conexión se recibirá mediante `AKRITAS_DATABASE_URL` y no
  se persistirá ni expondrá por API.

El Credential Store utilizará el mismo PostgreSQL, en tablas y repositorios
dedicados. Los valores almacenados serán ciphertext autenticado, nonce, versión y
metadata mínima de ownership. El plaintext y `AKRITAS_MASTER_KEY` nunca se
guardarán en PostgreSQL.

Para el MVP:

- `AKRITAS_MASTER_KEY` será una clave Base64 de 32 bytes;
- los secretos se cifrarán con AES-256-GCM;
- el AAD vinculará ciphertext, integración, clase de secreto y versión;
- metadata de integración y credenciales se modificarán en una misma transacción
  PostgreSQL cuando la operación requiera atomicidad;
- las transacciones se cerrarán antes de llamadas lentas o no confiables a
  GitHub o Dokploy;
- los tests de repositorios y migraciones usarán PostgreSQL real, no SQLite como
  sustituto semántico.

## Consecuencias

### Positivas

- un solo sistema de backup, migración, observabilidad y recuperación;
- transacciones atómicas entre metadata de integración y secretos cifrados;
- constraints e índices relacionales consistentes para todo el backend;
- comportamiento de tests equivalente al motor productivo;
- el Credential Store puede migrar detrás de su port a un secret manager futuro
  sin modificar dominio o usecases.

### Negativas

- el desarrollo y los tests de integración requieren una instancia PostgreSQL;
- comprometer simultáneamente PostgreSQL y la master key permite descifrar los
  secretos;
- perder la master key vuelve inaccesibles las credenciales almacenadas;
- un único cluster comparte blast radius operacional, aunque credenciales y
  repositorios permanezcan separados lógicamente.

## Alternativas descartadas

### SQLite separado para Credential Store

Se descarta porque introduce coordinación entre dos stores, backups separados y
fallos parciales sin mejorar la protección criptográfica frente a un compromiso
del host y la master key.

### Secret manager externo obligatorio

Vault, KMS y equivalentes continúan fuera del alcance del MVP según ADR-005.

### Persistencia en memoria

No satisface reinicios, migraciones ni el objetivo del Control Plane.

## Fuera de alcance

- alta disponibilidad y topología administrada de PostgreSQL;
- rotación automática de master key;
- envelope encryption con KMS;
- particionado, read replicas o multi-region;
- mover automáticamente credenciales existentes a un secret manager externo.
