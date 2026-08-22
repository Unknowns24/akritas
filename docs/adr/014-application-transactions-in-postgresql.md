# ADR-014 — Transacciones de aplicación sobre PostgreSQL compartido

## Estado

Accepted

## Contexto

Algunas invariantes de autenticación e integraciones abarcan más de un
repositorio: por ejemplo, activar un administrador requiere crear la entidad,
crear su sesión, trasladar el seed TOTP cifrado y consumir el enrollment como
una única operación. Resolver cada paso con su propia transacción permite
estados parciales y duplicar métodos transaccionales en cada puerto.

ADR-010 ya establece que metadata y Credential Store comparten PostgreSQL, por
lo que esas invariantes sí pueden tener un boundary atómico común.

## Decisión

`out.Transactor.WithinTransaction(ctx, fn)` es el único boundary genérico de
atomicidad a nivel de aplicación. El adapter PostgreSQL adjunta la transacción
activa al `context.Context`; los repositorios y el Credential Store participantes
la resuelven automáticamente y conservan sus puertos libres de GORM.

Se utiliza únicamente cuando una invariante corta requiere coordinar dos o más
operaciones sobre el mismo PostgreSQL. En particular:

- reemplazo de enrollment, password hash y seed TOTP;
- activación del administrador, sesión, traslado del seed y consumo del
  enrollment;
- compare-and-set del período TOTP y creación de sesión;
- metadata de integración y secretos cifrados cuando deban cambiar juntos.

Dentro de `WithinTransaction` está prohibido ejecutar llamadas HTTP, acceder a
proveedores externos, esperar interacción humana, realizar cómputo criptográfico
costoso o cualquier trabajo lento/no acotado. Esos pasos deben completarse antes
de abrir la transacción; sólo se persisten sus resultados dentro del callback.

Los repositorios no exponen variantes `SaveTx` ni reciben `*gorm.DB` por puertos.
Una operación aislada puede mantener su transacción local cuando no participa de
una invariant application-level.

## Consecuencias

- commit y rollback abarcan repositorios y Credential Store sin contaminar el
  core con infraestructura;
- el contexto pasa a transportar una capacidad interna del adapter y nunca debe
  reutilizarse después de finalizar el callback;
- las transacciones permanecen breves, reduciendo locks y agotamiento del pool;
- los tests PostgreSQL deben comprobar commit/rollback compartido y los tests de
  usecase deben comprobar que ningún paso posterior ocurre tras un fallo.

## Alternativas descartadas

- Coordinar compensaciones manuales para operaciones que caben en un único
  PostgreSQL: deja ventanas de inconsistencia innecesarias.
- Agregar métodos transaccionales a cada repositorio: duplica la API y filtra la
  decisión de infraestructura hacia application/core.
- Mantener transacciones abiertas durante llamadas externas: aumenta locks y no
  puede volver atómico un sistema remoto con PostgreSQL.
