# ADR-004 — La remediación automática termina en Pull Request

## Estado

Accepted

## Contexto

Akritas puede generar correcciones, pero aplicar cambios directamente en producción aumenta significativamente el riesgo del MVP.

## Decisión

Akritas puede crear branches, commits y Pull Requests, pero no puede mergear ni desplegar automáticamente durante el MVP.

## Consecuencias

- existe supervisión humana antes de modificar producción;
- la demo sigue mostrando trabajo real;
- se reduce el riesgo de remediaciones incorrectas;
- auto-merge y auto-deploy pueden evaluarse en versiones futuras.
