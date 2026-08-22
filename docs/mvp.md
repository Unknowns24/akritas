# Akritas — Hackathon MVP

## Objetivo

Construir una demostración funcional end-to-end del loop:

```text
DokployServer
    ↓
DokployApplication
    ↓
Project
    ↓
LogEvent
    ↓
Incident
    ↓
Investigation / QVAC
    ↓
GitHub Issue
    ↓
Optional Remediation
    ↓
Pull Request
```

El MVP debe demostrar que Akritas puede configurar integraciones reutilizables de GitHub y Dokploy, asociarlas a un Project, observar una aplicación desplegada, detectar un problema, investigarlo localmente y generar automáticamente los artefactos necesarios en GitHub.

## Must Have

### GitHub Account Management

- crear/configurar una `GitHubAccount`;
- soportar al menos una cuenta o instalación de GitHub;
- validar que la integración pueda autenticarse correctamente;
- listar o seleccionar repositorios accesibles mediante la cuenta configurada;
- reutilizar una misma `GitHubAccount` entre múltiples Projects.

Las credenciales y secretos utilizados para autenticarse contra GitHub deben mantenerse fuera del modelo de dominio.

### Dokploy Server Management

- crear/configurar un `DokployServer`;
- configurar su endpoint/base URL;
- validar conectividad con el servidor;
- listar o seleccionar aplicaciones disponibles;
- reutilizar un mismo `DokployServer` entre múltiples Projects.

Las API keys, tokens y secretos utilizados para conectarse a Dokploy deben mantenerse fuera del modelo de dominio.

### Project Management

- crear/configurar un `Project`;
- asociar un `GitHubRepository` perteneciente a una `GitHubAccount`;
- asociar una `DokployApplication` perteneciente a un `DokployServer`;
- configurar `MonitoringConfiguration`;
- activar/desactivar monitoreo.

El Project no administra directamente credenciales de GitHub ni Dokploy. Solo referencia recursos previamente configurados.

### Monitoring

- obtener logs de la `DokployApplication` asociada al Project;
- mantener cursor temporal o mecanismo equivalente para evitar reprocesar logs;
- detectar eventos relevantes mediante reglas determinísticas;
- generar `LogEvent`;
- crear fingerprints determinísticos;
- ignorar eventos según `ignored_patterns`;
- detectar errores según `error_patterns`;
- capturar contexto anterior y posterior cuando corresponda;
- agrupar ocurrencias equivalentes dentro de un mismo `Incident`;
- actualizar frecuencia, `first_seen_at` y `last_seen_at`.

La configuración debe provenir de `MonitoringConfiguration`.

### Incidents

- crear un `Incident` cuando corresponda;
- asociarlo al `Project`;
- asociar sus `LogEvent`;
- registrar fingerprint;
- registrar cantidad de ocurrencias;
- registrar `first_seen_at` y `last_seen_at`;
- mantener estado del incidente;
- mostrar estado de root cause;
- mostrar estado de resolución;
- almacenar confidence cuando exista;
- mantener referencia a la GitHub Issue creada;
- mantener referencia a la Pull Request cuando exista.

### Investigation

- crear una `Investigation` para cada ejecución de análisis sobre un Incident;
- registrar inicio y finalización;
- almacenar evidencia;
- registrar hipótesis;
- registrar root cause cuando pueda determinarse;
- registrar `root_cause_status`;
- registrar `resolution_status`;
- registrar confidence;
- identificar archivos relevantes;
- identificar commits relevantes cuando corresponda;
- registrar acciones recomendadas.

Un Incident puede tener múltiples Investigations.

### Evidence

La investigación debe poder utilizar evidencia proveniente de:

- fragmentos de logs;
- stack traces;
- archivos de código;
- símbolos o líneas relevantes;
- commits;
- diffs;
- resultados de validaciones o tests;
- metadata de deployment cuando esté disponible.

### QVAC

- inferencia 100% local;
- analizar un `Incident`;
- trabajar dentro del contexto de una `Investigation`;
- tool calling;
- inspeccionar el `GitHubRepository` asociado al Project;
- consultar evidencia proveniente de los `LogEvent`;
- identificar archivos y símbolos relevantes;
- generar hipótesis;
- clasificar `root_cause_status`;
- clasificar `resolution_status`;
- producir confidence;
- determinar si existe una remediación posible;
- producir output estructurado y determinístico en su contrato.

### GitHub Issue

- crear una `GitHubIssueReference` para todo Incident investigado;
- crear la Issue en el `GitHubRepository` asociado al Project;
- incluir información suficiente para comprender el incidente;
- incluir evidencia relevante;
- incluir resultado de la Investigation;
- vincular la Issue al Incident.

La creación de una Issue ocurre independientemente de si Akritas puede resolver automáticamente el problema.

### Remediation

Cuando el Incident sea considerado resoluble mediante cambios de código:

- crear una `Remediation`;
- crear una branch dedicada;
- producir cambios sobre el repositorio;
- registrar un resumen de cambios;
- ejecutar al menos una validación;
- almacenar resultados de validación;
- marcar la Remediation como fallida si los cambios no pueden validarse;
- crear una Pull Request únicamente cuando la validación sea satisfactoria.

Estados conceptuales mínimos:

- `planned`;
- `in_progress`;
- `validated`;
- `failed`;
- `pull_request_created`.

### Pull Request

Cuando una Remediation validada pueda proponerse:

- crear una `PullRequestReference`;
- crear la Pull Request en el repositorio correspondiente;
- utilizar la branch generada por la Remediation;
- vincularla a la Issue correspondiente;
- vincularla a la Remediation;
- hacer visible su referencia desde el Incident.

### UI

Dashboard mínimo para administrar y visualizar:

#### Integrations

- GitHub Accounts;
- estado de autenticación de cada GitHub Account;
- Dokploy Servers;
- estado de conexión de cada Dokploy Server.

#### Projects

- Projects;
- GitHub Repository asociado;
- Dokploy Application asociada;
- estado de monitoring;
- configuración básica de monitoring.

#### Incidents

- Incidents;
- cantidad de ocurrencias;
- first seen;
- last seen;
- estado;
- investigación;
- root cause status;
- resolution status;
- confidence;
- link a Issue;
- estado de Remediation;
- link a Pull Request si existe.

## Should Have

- correlación con commits recientes;
- detectar que un Incident comenzó después de un cambio;
- identificar commits relevantes dentro de una Investigation;
- mostrar timeline de Investigation;
- mostrar tools utilizadas por QVAC;
- mostrar evidence utilizada para llegar a una conclusión;
- confidence score;
- agrupación visual de LogEvents y ocurrencias;
- visualizar GitHub Account y Dokploy Server asociados indirectamente al Project;
- distinguir entre fallo de monitoreo, fallo de integración y fallo de la aplicación observada.

## Nice to Have

- metadata del último deployment;
- relacionar deployment con commits;
- identificar el diff que probablemente introdujo una regresión;
- actualización automática de la Issue cuando falla una Remediation;
- actualización automática de la Issue cuando se crea una Pull Request;
- configuración de patrones custom por Project;
- prueba manual de conexión para `GitHubAccount`;
- prueba manual de conexión para `DokployServer`;
- descubrimiento automático de repositorios disponibles;
- descubrimiento automático de aplicaciones Dokploy disponibles.

## Explicitly Not in MVP

- auto merge;
- auto deploy;
- rollback;
- métricas;
- traces;
- Slack;
- Kubernetes;
- multi-provider observability;
- múltiples agentes;
- permisos empresariales complejos;
- billing;
- multi-tenancy avanzado;
- gestión avanzada de secretos;
- rotación automática de credenciales;
- múltiples proveedores de Git;
- múltiples proveedores de deployment.

## Success Criteria

El MVP se considera exitoso si durante una demo controlada Akritas puede:

1. configurar una `GitHubAccount`;
2. configurar un `DokployServer`;
3. crear un `Project`;
4. asociar al Project un `GitHubRepository` accesible mediante la GitHub Account;
5. asociar al Project una `DokployApplication` perteneciente al Dokploy Server;
6. activar el monitoreo;
7. detectar un error real producido por una aplicación de prueba;
8. generar uno o más `LogEvent`;
9. agrupar eventos equivalentes en un `Incident`;
10. iniciar una `Investigation`;
11. usar QVAC local para investigar el Incident;
12. consultar código relevante del repositorio;
13. registrar evidencia, hipótesis y root cause;
14. crear una GitHub Issue útil;
15. identificar el Incident como resoluble mediante código;
16. crear una `Remediation`;
17. modificar el repositorio en una branch dedicada;
18. ejecutar al menos una validación;
19. crear una Pull Request vinculada a la Issue;
20. visualizar desde Akritas el recorrido completo del Incident desde su detección hasta la Pull Request.
