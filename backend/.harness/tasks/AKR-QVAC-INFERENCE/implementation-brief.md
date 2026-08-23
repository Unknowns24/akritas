# Implementation brief — AKR-QVAC-INFERENCE

## Estado inicial

`out.InvestigationRunner` está implementado solo por `qvac.StubRunner` (siempre
falla). `RunUseCase.Execute` ya orquesta assemble→run→Complete/Fail.
`InvestigationRunResult` y los enums `RootCauseStatus`/`ResolutionStatus` ya
existen. QVAC se corre como `qvac serve openai` (HTTP OpenAI-compatible en
loopback); no se usa sidecar Node propio ni SDK embebido.

## Estrategia

```text
RunUseCase.Execute
  -> qvac.Runner.Run(investigation)
       -> POST {base}/chat/completions  (model=akritas, response_format=json_schema, sin tools)
       -> parse + Validate enums/confidence/required fields
       -> out.InvestigationRunResult
  -> Investigation.Complete(...) | failInvestigation(err)
```

Base URL default: `http://127.0.0.1:11434/v1`, solo loopback/private; auth
opcional Bearer; model default `akritas`.

## Componentes

1. `Client` — HTTP hacia QVAC (`ChatCompletions`), timeouts, normalización de
   errores a `domain.ErrIntegrationUnavailable` / errores de paquete.
2. `Runner` — arma mensajes (system: rol investigativo; user: contexto de
   Investigation ID + instrucciones de JSON); pide schema alineado al
   contrato; parsea.
3. `parseResult` — JSON → `InvestigationRunResult`; rechaza valores que no
   matchean `Validate()` exacto.
4. Bootstrap: `qvac.NewRunner(cfg)` en lugar de `NewStubRunner()`.

## Seguridad / ADR-001

- Solo inferencia local (URL no pública).
- No loguear prompts completos ni API keys.
- Output inválido nunca llega a `Complete`.

## Test strategy

httptest simula QVAC; tests de client, parse, runner feliz/error, y un test
que ejecuta `RunUseCase` + runner real contra httptest y verifica campos
persistidos en el store fake (PB-035).
