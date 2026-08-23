# Akritas ERD Notes

Current persistence uses PostgreSQL migrations under `internal/adapter/db/postgres/migrations`.

## Incident and Investigation Publication

```text
projects
  id PK
    |
    | 1:N
    v
incidents
  id PK
  project_id FK -> projects.id
    |
    | 1:N
    v
investigations
  id PK
  incident_id FK -> incidents.id
  UNIQUE (id, incident_id)
    |
    | 1:N
    v
evidence
  id PK
  investigation_id FK -> investigations.id

github_issue_references
  investigation_id PK, FK -> investigations.id
  incident_id FK -> incidents.id
  FK (investigation_id, incident_id) -> investigations(id, incident_id)
  UNIQUE (repository, issue_number)
```

The composite FK prevents a GitHub IssueReference from combining an existing Incident with an Investigation that belongs to another Incident.

