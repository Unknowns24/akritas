# Akritas — Domain Model

Este documento define las entidades conceptuales del dominio. No prescribe tablas, ORM ni detalles de persistencia.

## Administrator

Representa la única identidad humana autorizada a administrar una instalación
Akritas durante el MVP.

Campos conceptuales públicos:

- `id`
- `email`
- `display_name`
- `created_at`
- `updated_at`

El password hash, el secreto TOTP, el bootstrap token y los identificadores de
sesión no forman parte de esta proyección pública. Password y sesión pertenecen al
subdominio de autenticación; el secreto TOTP se cifra mediante infraestructura.

## AdministratorSession

Representa una sesión opaca server-side.

Campos conceptuales seguros:

- `administrator_id`
- `authenticated_at`
- `idle_expires_at`
- `absolute_expires_at`
- `revoked_at` opcional

El token entregado al navegador se persiste únicamente en forma no reversible y
nunca se devuelve dentro de un DTO JSON.

## Project

Representa una aplicación o sistema que Akritas debe observar.

Campos conceptuales:

- `id`
- `name`
- `description`
- `monitoring_status`
- `github_repository`
- `dokploy_application`
- `monitoring_configuration`
- `created_at`
- `updated_at`

Un Project referencia recursos previamente configurados en Akritas. No contiene credenciales ni secretos de acceso a GitHub o Dokploy.

## GitHubAccount

Representa una cuenta, organización o instalación de GitHub configurada en Akritas.
Permite reutilizar una misma integración de GitHub entre múltiples Projects y repositorios.
Campos conceptuales:

- `id`
- `display_name`
- `account_type`
- `authentication_method`
- `account_identifier`
- `authentication_status`
- `created_at`
- `updated_at`

Valores conceptuales posibles para `account_type`:

- `personal`
- `organization`

Valores conceptuales posibles para `authentication_method`:

- `personal_access_token`
- `github_app`

Las credenciales, tokens, private keys y secretos de autenticación no forman parte conceptual de esta entidad; pertenecen a la capa de integración y configuración segura.

## GitHubRepository

Referencia al repositorio que contiene el código del Project.
Campos conceptuales:

- `github_account_id`
- `owner`
- `name`
- `default_branch`
- `repository_identifier`

El repositorio pertenece conceptualmente a una `GitHubAccount` configurada en Akritas.
Las credenciales no forman parte conceptual del repositorio; pertenecen a la capa de integración/configuración segura.

## DokployServer

Representa una instancia de Dokploy configurada en Akritas.
Permite reutilizar una misma conexión a Dokploy entre múltiples Projects y aplicaciones desplegadas.
Campos conceptuales:

- `id`
- `name`
- `base_url`
- `server_identifier`
- `connection_status`
- `created_at`
- `updated_at`

Las API keys, tokens y demás secretos necesarios para conectarse al servidor no forman parte conceptual de esta entidad; pertenecen a la capa de integración y configuración segura.

## DokployApplication

Referencia a la aplicación desplegada que produce los logs observados.
Campos conceptuales:

- `dokploy_server_id`
- `application_identifier`
- `instance_identifier`
- `display_name`

La aplicación pertenece conceptualmente a un `DokployServer` configurado en Akritas.

## MonitoringConfiguration

Reglas aplicadas durante la observación del Project.
Campos conceptuales:

- `enabled`
- `error_patterns`
- `ignored_patterns`
- `grouping_window`
- `context_before`
- `context_after`

## LogEvent

Evento individual relevante detectado en los logs.
No representa necesariamente un incidente.
Campos conceptuales:

- `id`
- `project_id`
- `timestamp`
- `severity`
- `message`
- `raw_context`
- `fingerprint`

## Incident

Agrupa uno o más LogEvents que Akritas considera manifestaciones del mismo problema.
Campos conceptuales:

- `id`
- `project_id`
- `fingerprint`
- `phase`
- `terminal_outcome` opcional
- `first_seen_at`
- `last_seen_at`
- `occurrence_count`
- `title`
- `summary`
- `root_cause_status`
- `resolution_status`
- `confidence`
- `github_issue_reference`
- `pull_request_reference` opcional

Valores conceptuales posibles para `phase`:

- `detected`
- `investigating`
- `publishing_issue`
- `remediating`
- `completed`
- `failed`

`completed` significa que el workflow de Akritas alcanzó un resultado terminal;
no afirma que el cambio fue mergeado, desplegado o resuelto en producción.

## Investigation

Representa el proceso y resultado del análisis de un Incident.
Campos conceptuales:

- `id`
- `incident_id`
- `started_at`
- `finished_at`
- `evidence`
- `hypotheses`
- `root_cause`
- `root_cause_status`
- `resolution_status`
- `confidence`
- `relevant_files`
- `relevant_commits`
- `recommended_actions`

## Evidence

Unidad de información usada para justificar una conclusión.
Puede representar:

- fragmento de logs;
- stack trace;
- archivo de código;
- línea o símbolo;
- commit;
- diff;
- resultado de test;
- metadata del deployment.

## Remediation

Intento de resolver un Incident mediante cambios de código.
Campos conceptuales:

- `id`
- `incident_id`
- `status`
- `branch_name`
- `changes_summary`
- `validation_results`
- `pull_request_reference`

Estados posibles conceptuales:

- `planned`
- `in_progress`
- `validated`
- `failed`
- `pull_request_created`

## GitHubIssueReference

Referencia a la Issue creada para el Incident.
Campos conceptuales:

- `number`
- `url`
- `repository`

## PullRequestReference

Referencia a la Pull Request creada para una Remediation.
Campos conceptuales:

- `number`
- `url`
- `repository`
- `branch`

## Relaciones

```text
GitHubAccount *
 └── GitHubRepository *

DokployServer *
 └── DokployApplication *

Project
 ├── GitHubRepository 1
 ├── DokployApplication 1
 ├── MonitoringConfiguration 1
 └── Incident *
       ├── LogEvent *
       ├── Investigation 1..*
       ├── GitHubIssueReference 1
       └── Remediation 0..1
              └── PullRequestReference 0..1
```

## Reglas conceptuales

- Una instalación MVP posee exactamente un `Administrator` activo.
- Confirmar recovery invalida las sesiones y el enrollment TOTP anteriores.
- Un `GitHubAccount` puede ser utilizado por múltiples `GitHubRepository`.
- El tipo de cuenta GitHub y su método de autenticación son independientes.
- Un `DokployServer` puede contener múltiples `DokployApplication`.
- Múltiples Projects pueden utilizar recursos pertenecientes a la misma cuenta de GitHub o al mismo servidor Dokploy.
- Un `Project` observa una aplicación desplegada y conoce el repositorio de código asociado.
- Las credenciales y secretos de GitHub y Dokploy no pertenecen al modelo de dominio.
- Las integraciones externas deben poder configurarse independientemente de los Projects.
- Un `LogEvent` no implica necesariamente la existencia de un nuevo `Incident`.
- Los LogEvents equivalentes pueden agruparse dentro de un mismo `Incident` mediante su fingerprint.
- Un `Incident` puede ser investigado múltiples veces.
- Todo Incident debe producir una GitHub Issue.
- Una `Remediation` solo existe cuando Akritas intenta resolver el Incident mediante cambios de código.
- Una Pull Request solo existe como resultado de una Remediation.
