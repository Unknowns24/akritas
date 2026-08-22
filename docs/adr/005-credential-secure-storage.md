# ADR-005 — Almacenar credenciales en un Credential Store cifrado

## Estado

Accepted

## Contexto

Akritas necesita autenticarse contra integraciones externas como GitHub y Dokploy.

Estas credenciales son necesarias para operar el sistema, pero no forman parte del modelo de dominio. Un `GitHubAccount` representa una integración configurada con GitHub y un `DokployServer` representa una instancia configurada de Dokploy; ninguno de ellos debe contener conceptualmente tokens, API keys, private keys u otros secretos.

Guardar secretos en texto plano junto con los datos normales de la aplicación aumentaría innecesariamente el impacto de una filtración de la base de datos. Al mismo tiempo, incorporar un sistema externo de gestión de secretos como Vault excede el alcance del MVP y dificultaría la ejecución local y la demo.

## Decisión

Akritas tendrá un `Credential Store` perteneciente a la capa de infraestructura.

El Credential Store almacenará los secretos requeridos por las integraciones cifrados en reposo y asociados mediante identificadores internos a la integración correspondiente.

Los objetos del dominio sólo mantendrán su identidad y estado de integración. No contendrán el valor de las credenciales ni dependerán de su representación persistida.

El material necesario para descifrar las credenciales se obtendrá desde una master key proporcionada al proceso de Akritas desde fuera de la persistencia de la aplicación.

La master key:

- no se almacenará en la base de datos;
- no formará parte del modelo de dominio;
- no será expuesta mediante la API o UI;
- podrá ser provista mediante configuración segura del runtime, por ejemplo una variable de entorno o secret del entorno de despliegue.

Los adapters de integración accederán a las credenciales únicamente cuando necesiten realizar una operación contra el proveedor externo.

Conceptualmente, el acceso seguirá el flujo:

```text
GitHubAccount / DokployServer
            ↓
      integration id
            ↓
      Credential Store
            ↓
   decrypt at point of use
            ↓
 GitHub / Dokploy adapter
```

La API y la UI podrán indicar si una credencial está configurada y si la autenticación es válida, pero no devolverán el secreto almacenado una vez creado.

## Reglas de seguridad

- los secretos deben estar cifrados en reposo;
- el texto plano sólo debe existir en memoria durante el tiempo necesario para utilizarlo;
- secretos y master keys nunca deben aparecer en logs;
- secretos nunca deben incorporarse al contexto enviado a QVAC;
- secretos nunca deben incluirse en Issues, Pull Requests, evidencias o mensajes de error persistidos;
- los permisos solicitados a cada integración deben ser los mínimos necesarios para las capacidades habilitadas;
- actualizar una credencial reemplaza su valor almacenado sin modificar la identidad de la integración del dominio.

## Consecuencias

### Positivas

- las credenciales permanecen fuera del modelo de dominio;
- una copia de la base de datos no contiene secretos directamente utilizables sin la master key;
- GitHub Accounts y Dokploy Servers pueden reutilizarse entre múltiples Projects sin duplicar credenciales;
- la solución funciona tanto en ejecución local como desplegada;
- los adapters de integración tienen un único mecanismo para solicitar secretos;
- permite sustituir el Credential Store por un secret manager externo en el futuro sin modificar el dominio.

### Negativas

- Akritas debe gestionar cifrado y descifrado de secretos;
- perder la master key hace inaccesibles las credenciales almacenadas;
- comprometer simultáneamente la persistencia y la master key permite recuperar los secretos;
- rotación automática, versionado avanzado y secret managers externos quedan fuera del MVP.

## Fuera de alcance del MVP

- HashiCorp Vault u otros secret managers externos obligatorios;
- rotación automática de credenciales;
- rotación automática de la master key;
- envelope encryption con KMS;
- sincronización de secretos entre múltiples instancias de Akritas;
- recuperación del valor original de una credencial desde la UI.
