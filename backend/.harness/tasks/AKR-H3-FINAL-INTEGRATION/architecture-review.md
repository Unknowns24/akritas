# Architecture Review — AKR-H3-FINAL-INTEGRATION

## Veredicto

PASS

## Dirección de dependencias

- El core sólo declara contratos de aplicación/dominio.
- Evidence assembly depende de ports H2, no de GORM.
- El usecase resuelve Incident/Project/repository/Evidence antes de invocar `InvestigationRunner`.
- QVAC implementa el port runner y recibe contexto; no importa PostgreSQL.
- GitHub implementa `RepositoryInspector`; QVAC no permite que el modelo elija otro repositorio.
- REST usa los usecases compuestos y expone simultáneamente rutas H2/H3.

## SRP e interfaces

`IncidentReader` fue eliminado por duplicar H2. Las capacidades `IncidentGetter`, `IncidentLogEventLister` e `IncidentWorkflowStore` permiten que REST, assembly y workflow dependan de la mínima interfaz necesaria. `IncidentStore` sólo compone las lecturas H2 publicadas.

## Transacciones y lifecycle

Las transacciones rodean sólo cambios PostgreSQL; QVAC/GitHub nunca se ejecutan bajo una transacción. Start despacha post-commit. El resultado y Operation se cierran atómicamente; ante fallo de esa persistencia se conserva el estado en memoria `running` para intentar una transición failed separada. Recovery reconcilia trabajo durable antes de servir HTTP.

## Persistencia

Las migraciones son aditivas y ordenadas. Los FKs RESTRICT conservan audit history y el índice único parcial refuerza la regla activa frente a carreras. Evidence citada queda separada del count total mediante `evidence_ids`.

## Scope

No se implementaron H4/H5, publicación, branches, writes, remediation, commits, pushes ni PRs. `StartIssuePublication` sólo permanece en el dominio/tests de fases posteriores y no es llamado por wiring H3.

## Gate

`.harness/kernel/scripts/check-backend-architecture.sh`: PASS.
