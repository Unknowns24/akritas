# TDD test plan — AKR-QVAC-TOOL-LOOP

*(Aprobado vía instrucción de ejecución del plan H3 restante.)*

- Loop: modelo pide tool fake → se ejecuta → segundo turno sin tools → JSON final OK.
- Tool desconocida → mensaje de error en role=tool, sin panic; investigación puede completar o fallar según turno final.
- Args JSON inválidos → error controlado a la tool response.
- Límite de iteraciones excedido → error del runner (Investigation failed).
- Registry vacío: comportamiento de AKR-QVAC-INFERENCE (solo fase estructurada) se preserva.
