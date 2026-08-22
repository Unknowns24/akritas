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
