# Implementation Summary

## Task

`AKR-OPENAPI-MVP` — contrato OpenAPI v1 completo, autenticación TOTP y contratos del MVP.

## Implemented

- Se creó `docs/openapi.yaml` como contrato canónico OpenAPI 3.1.0, API `1.0.0`, con base `/api/v1`.
- El contrato contiene 45 paths, 59 operaciones y 112 schemas reutilizables para autenticación, sistema, dashboard, operaciones asíncronas, integraciones, proyectos, automatización e incident lifecycle.
- Se agregaron envelopes individuales y paginados, cursores, errores seguros, UUID, timestamps RFC 3339, duraciones ISO 8601, ejemplos y `operationId` únicos.
- Se incorporó autenticación single-admin con bootstrap, password, TOTP, sesión opaca y recovery con rotación del factor.
- Se modelaron GitHub PAT y GitHub App Manifest como métodos independientes del tipo de cuenta.
- Se preservó el límite humano en Pull Request: el contrato no expone merge, deploy ni un falso estado de resolución productiva.
- Se reforzó el gate OpenAPI para validar estructura, referencias, seguridad pública, secretos, defaults, enums, idempotencia, exclusiones y ejemplos.
- Se corrigió el gate de seguridad para que su propia expresión de detección no produzca un falso positivo.
- Se actualizaron spec, MVP, backlog, dominio, demo, diseño, integraciones, configuración y arquitecturas; además se agregaron ADR-008 y ADR-009.
- Se agregó una matriz de cobertura para PB-001–PB-055, PB-061–PB-067 y las pantallas compatibles de Stitch.

## Deliberately not implemented

- Handlers, casos de uso, persistencia, workers y adapters Go.
- Frontend Next.js.
- Team, RBAC, multi-tenancy, incidentes manuales, detección por IA, merge y deploy.

## Notes

- Se actualizó la fuente `docs/img/diagrams/erd.plantuml`. El PNG existente no se regeneró porque PlantUML no está instalado en el entorno.
- El archivo no relacionado `qvac.config.json` fue preservado sin cambios.
