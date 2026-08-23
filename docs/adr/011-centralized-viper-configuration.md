# ADR-011 — Configuración centralizada mediante Viper

## Estado

Accepted

## Contexto

La configuración runtime de Akritas estaba distribuida entre varios archivos y
parte de los defaults operativos vivía en el composition root. Esto dificulta
descubrir, validar y mantener las variables soportadas.

## Decisión

La configuración runtime del backend se define en un único
`config/config.go` y se carga mediante `github.com/spf13/viper`.

- `Config` es el contrato tipado único para bootstrap y adapters.
- Las variables conservan el prefijo `AKRITAS_`.
- Un archivo `app.env` es opcional; las variables de entorno tienen precedencia.
- Los defaults no sensibles se declaran junto al loader.
- La validación devuelve errores y falla cerrada; no usa `panic` ni incluye
  valores secretos en mensajes.
- Los tests usan instancias Viper aisladas para evitar estado global compartido.

## Consecuencias

- La configuración disponible y sus defaults quedan visibles en un solo lugar.
- Bootstrap deja de contener valores configurables dispersos.
- Viper se convierte en dependencia estable del boundary de configuración.
- Agregar una variable exige actualizar `Config`, defaults, validación y
  `docs/configuration.md`.

