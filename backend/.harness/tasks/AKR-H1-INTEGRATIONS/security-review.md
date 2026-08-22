# Security Review

## Veredicto

**Aprobado sin hallazgos bloqueantes para el código entregado. El runtime
permanece fail-closed.**

## Configuración y secretos

- Viper usa instancia local, environment sobre `app.env`, defaults no sensibles
  y errores que no imprimen valores.
- Public URL, PostgreSQL DSN, master key Base64 de 32 bytes, pagination secret,
  TTL y pool se validan antes de bootstrap.
- Los strings raw de master key/pagination se limpian de la configuración
  devuelta; adapters reciben bytes ya validados.
- AES-256-GCM usa nonce aleatorio y AAD; ciphertext/nonce sólo existen en el
  record DB privado.

## Credenciales y providers

- PAT/API key son write-only; private key/webhook secret están cifrados;
  installation tokens sólo viven en memoria.
- GitHub valida credenciales e instalación antes de persistir, usa API versionada
  y rechaza redirects no confiables.
- Dokploy bloquea userinfo, query, fragment, metadata/link-local/multicast,
  redirects prohibidos y DNS rebinding; `x-api-key` no cruza de origen.
- Bodies externos y tiempos están acotados; errores normalizados no exponen
  responses de proveedor.

## Uker y REST

- Cursores Uker usan firma HMAC y TTL; tests cubren manipulación, expiración,
  overrides, allowlists y reutilización entre integraciones.
- Scope, filtros, sort, límite y provider boundary viajan firmados.
- DTOs/envelopes seguros no contienen secretos, ciphertext ni nonce.
- Los callbacks usan state persistido de un uso y `Cache-Control: no-store`.

## Auth fail-closed

- Router y bootstrap rechazan middleware administrativo ausente o nulo.
- El chequeo ocurre antes de conectar PostgreSQL o ejecutar migraciones.
- La implementación concreta de sesión y Origin es PB-061..063; por eso el
  módulo no se monta desde `cmd/main.go`.
- Los deletes también quedan bloqueados por el reader Project fail-closed hasta
  PB-010/011.

## Validación

- `go test ./...`, `go test -race ./...`, `go vet ./...` y el gate de seguridad
  pasan.
- `go test -tags=integration ./...` pasa; el caso Testcontainers PostgreSQL se
  omite explícitamente porque Docker no está activo.
