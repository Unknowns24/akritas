# AKR-H3-FINAL-INTEGRATION — Finish H3 after H2/H3 merge

## Estado

complete — implementación y verificación finalizadas el 2026-08-23

## Tipo de tarea

backend-api-feature

## Modo de proyecto

existing_project

## Contexto

H2 e H3 están físicamente combinados, pero H3 conserva supuestos temporales de cuando Incident y LogEvent aún no existían. El catálogo PostgreSQL quedó roto por el merge y la composición productiva no compila por el port `IncidentReader.Exists` legado.

## Objetivo

Completar el flujo productivo `Incident → LogEvents → Evidence → QVAC local → GitHub read-only tools → structured result → PostgreSQL`, sin implementar H4/H5.

## Criterios de aceptación

- AKR-35 a AKR-44 satisfechos con wiring productivo y pruebas sobre H2 real.
- Evidence inicial acotada a 24 KiB máximos y reducida según `ctx_size=16384`; corpus persistido acotado a 128 KiB/25 items.
- H3 completa Investigation sin ejecutar `StartIssuePublication`.
- Tools limitadas a las cinco operaciones GitHub de lectura y scope fijado por Project.
- Build, suite, race, vet, integración PostgreSQL y gates del harness ejecutados sin skips relevantes.

## Fuera de alcance

GitHub Issues, ramas, escrituras de código, commits, Pull Requests y remediation.

## Instrucción para el harness

El implementation brief y el plan TDD fueron aprobados explícitamente por el usuario antes de implementar.
