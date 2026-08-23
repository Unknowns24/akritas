# Security Review

## Summary

El cambio repara el transporte del nonce público requerido por GitHub sin
alterar su generación, almacenamiento o validación.

## Auth / permissions

No se modifican middleware, autenticación administrativa, callbacks públicos ni
permisos de la GitHub App.

## Input validation

`CompleteManifest` conserva límites de `code` y `state`. Los tests confirman que
state ausente y vencido fallan antes del exchange y que el replay falla después
del primer consumo.

## Data exposure

El nonce aparece en la URL porque forma parte pública del protocolo. Persistencia
continúa almacenando sólo SHA-256. No se exponen private key, webhook secret o
client secret, y el Manifest mantiene webhooks inactivos y permisos mínimos
existentes.

## Error leakage

No se agregan ni modifican errores. Los fallos conservan los errores de dominio
normalizados actuales.

## Findings

Sin findings bloqueantes ni no bloqueantes.

## Result

pass
