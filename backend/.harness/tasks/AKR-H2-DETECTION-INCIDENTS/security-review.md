# Security review — AKR-H2-DETECTION-INCIDENTS

## Resultado

Aprobado.

- Las credenciales Dokploy se resuelven por el store existente y se limpian
  después del request; no entran en LogEvent ni checkpoint.
- Bearer/Basic auth, secretos etiquetados, PATs conocidos y PEM privados se
  redactan antes de construir evidencia o estado durable.
- La respuesta del proveedor y el evento lógico conservan límites explícitos.
- Los endpoints permanecen dentro del grupo autenticado y con política Origin
  ya establecida por el router.
- La API omite source interno, occurrence key y estado operacional.
- No se agregaron logs con payloads, tokens o causas internas.

El gate `check-security.sh` y la inspección de patrones de secretos pasan. Los
únicos valores con forma de secreto hallados son fixtures deliberados de tests
de redacción/credenciales.
