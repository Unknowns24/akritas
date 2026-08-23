# Final summary — AKR-H2-DETECTION-INCIDENTS

AKR-23 a AKR-34 quedaron implementadas como un pipeline H2 único y
determinístico, sin QVAC.

## Verificación

- `go fmt ./...`: pass
- `go test ./...`: pass
- `go test -race ./...`: pass
- `go test -tags=integration ./...`: pass; los casos PostgreSQL reales se
  declararon `SKIP` porque no hay PostgreSQL local ni Docker/Testcontainers
  disponible en el entorno
- `go vet ./...`: pass
- arquitectura, OpenAPI y seguridad: pass
- `git diff --check`: pass

## Contrato y discrepancias

OpenAPI es 1.5.0 e incorpora únicamente `initial_log_ingestion`. Los tres
endpoints H2 de Incident ya estaban publicados y fueron implementados sin
campos frontend adicionales. La única discrepancia externa confirmada sigue
siendo la limitación de Dokploy: `application.readLogs` ofrece `tail` y `since`
relativo, no cursor absoluto ni streams separados; se resolvió con checkpoint
timestamp/ordinal/hash, overlap y continuidad fail-closed según el plan
aprobado.
