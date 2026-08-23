# Architecture review — AKR-H2-DETECTION-INCIDENTS

## Resultado

Aprobado para el profile `backend_api`.

- La detección pura vive en service y depende sólo de dominio.
- Dokploy implementa el nuevo port `LogSource` reutilizando Credential Store,
  transporte SSRF-safe, timeout y `x-api-key`.
- Los usecases y servicios coordinan ports; GORM, SQL, JSONB y locks permanecen
  en adapters PostgreSQL.
- No se introdujo AutoMigrate ni un cliente/proceso HTTP paralelo.
- REST implementa exclusivamente los paths y schemas de `docs/openapi.yaml`.
- El runner pertenece al bootstrap, no a handlers, y termina por contexto.
- QVAC y los componentes H3+ no tienen dependencias ni participación en H2.

El gate `check-backend-architecture.sh` pasa sin excepciones nuevas.
