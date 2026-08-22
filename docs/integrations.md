# Akritas — Integrations

## GitHub

GitHub cumple dos funciones principales:

1. fuente de verdad del código;
2. destino auditable del trabajo generado por Akritas.

### Capacidades necesarias en el MVP

Lectura:

- obtener metadata del repositorio;
- leer archivos;
- buscar código;
- inspeccionar branches;
- listar commits recientes;
- leer commits y diffs.

Escritura:

- crear Issues;
- crear branches;
- crear/modificar archivos o commits;
- crear Pull Requests.

### Reglas

- Akritas nunca debe hacer merge automático en el MVP.
- Toda PR creada debe estar asociada a una Issue existente.
- Toda Issue debe poder entenderse sin acceder al estado interno de Akritas.

### Métodos de autenticación

El MVP soporta:

```text
Personal Access Token
GitHub App Manifest
```

Un PAT se recibe una vez mediante un campo write-only, se valida y se cifra en el
Credential Store. Nunca se devuelve, ni siquiera al administrar la conexión.

El flujo GitHub App Manifest genera un `state` impredecible y de un solo uso,
intercambia el `code` en el backend y cifra la private key y webhook secret antes
de persistir la referencia. El `installation_id` recibido por callback debe
verificarse con autenticación de la App.

Permisos de la GitHub App del MVP:

- metadata read;
- contents read/write;
- issues write;
- pull requests write.

Los webhooks permanecen desactivados. GitHub Enterprise queda fuera de alcance.

## Dokploy

Dokploy es la fuente inicial de observabilidad runtime del MVP.

### Capacidades necesarias

- identificar una aplicación configurada;
- obtener logs recientes;
- limitar la cantidad de logs consultados;
- obtener logs posteriores a un punto temporal conocido;
- opcionalmente realizar búsquedas textuales sobre logs;
- distinguir logs runtime de logs de deployment cuando sea necesario.

### Uso dentro de Akritas

Akritas no requiere que Dokploy detecte incidentes.
Dokploy proporciona los logs; el Detection Engine de Akritas interpreta las señales.

## Credenciales

Las integraciones requieren secretos de acceso.

Principios:

- secretos cifrados en reposo;
- nunca incluir tokens en prompts;
- nunca incluir secretos en Issues o Pull Requests;
- minimizar permisos;
- separar permisos de lectura y escritura cuando sea práctico.
- responder únicamente con proyecciones como `credential_configured` y estados normalizados.

## QVAC

QVAC es el único runtime de inferencia y debe permanecer en loopback o una red
privada. Su configuración puede usar autenticación `none`, `bearer` o `basic`,
pero esos secretos se tratan igual que las demás credenciales.

El adapter debe resolver y validar el host antes de conectarse, rechazar destinos
públicos y no seguir redirects hacia una dirección pública. El frontend sólo
recibe endpoint seguro, tipo de autenticación, `credential_configured`, modelo,
versión, latencia y estado normalizado.

## Integraciones futuras

Fuera del MVP:

- GitLab;
- Bitbucket;
- Docker directo;
- Kubernetes;
- Loki;
- Grafana;
- OpenTelemetry;
- Sentry;
- cloud providers;
- Slack/Teams/Discord.
