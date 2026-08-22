# Ákritas

Akritas es un sistema autónomo de incident response para aplicaciones desplegadas en producción.

Su objetivo es observar errores de runtime, investigar su causa utilizando contexto de logs y código fuente, registrar siempre el incidente en GitHub y, cuando el problema sea solucionable de forma segura, generar una Pull Request con una propuesta de corrección.

La inferencia de IA de Akritas corre localmente mediante QVAC.

## Objetivo del MVP

Demostrar un ciclo completo y verificable:

1. Akritas observa los logs de una aplicación desplegada en Dokploy.
2. Detecta una anomalía o error relevante.
3. Agrupa eventos relacionados dentro de un incidente.
4. Analiza el incidente utilizando QVAC.
5. Inspecciona el repositorio GitHub asociado al proyecto.
6. Formula una hipótesis de causa raíz.
7. Crea siempre una Issue en GitHub documentando el incidente.
8. Si el incidente es solucionable automáticamente, genera una corrección y crea una Pull Request.

## Principios

- La Issue es el registro canónico de todo incidente detectado.
- La remediación automática es opcional y depende del nivel de confianza y del tipo de problema.
- Akritas no debe auto-mergear ni desplegar automáticamente en el MVP.
- El modelo local no debe procesar cada línea de logs individualmente; primero existe una etapa determinística de detección y agrupación.
- Logs, código fuente y contexto operativo no deben enviarse a proveedores externos de IA.

## Documentación

- `spec.md`: especificación funcional del sistema.
- `domain.md`: entidades y conceptos del dominio.
- `incident-lifecycle.md`: ciclo de vida de un incidente.
- `integrations.md`: responsabilidades de GitHub y Dokploy.
- `agent.md`: responsabilidades y límites del agente QVAC.
- `mvp.md`: alcance concreto para el hackathon.
- `demo.md`: historia de demo objetivo.
- `adr/`: decisiones de arquitectura aceptadas.
- `rfc/`: ideas futuras fuera del MVP.
