# TDD test plan — AKR-QVAC-INFERENCE

Plan sujeto a aprobación explícita antes de escribir tests o implementación.
*(Aprobado vía instrucción de ejecución del plan H3 restante.)*

## Client

- POST chat/completions feliz → decodifica content del assistant.
- Timeout / connection refused → error propagable (no resultado inventado).
- HTTP 404/5xx / modelo ausente → error.
- URL no loopback/private rechazada en construcción.

## Parse / validación (PB-032/033/034)

- JSON completo válido → InvestigationRunResult con enums exactos.
- `root_cause_status` o `resolution_status` desconocido → error.
- confidence fuera de [0,1] → error.
- summary vacío / campos requeridos ausentes → error.
- JSON malformado → error.

## Runner

- Happy path httptest → result poblado.
- Respuesta inválida → error (RunUseCase marcará failed).

## Persistencia (PB-035)

- RunUseCase + Runner real (httptest) + stores fake: tras Execute exitoso,
  Investigation actualizada contiene summary, root_cause, statuses,
  confidence, relevant_files/commits, recommended_actions, hypotheses.
- Fallo QVAC → Investigation failed, sin Complete parcial.

## Validación final

- `go test ./...`, `go vet`, `gofmt`, check-backend-architecture,
  check-security (y check-openapi si aplica; esta tarea no toca OpenAPI).
- Solo fallos preexistentes: TestCreatePersistsAdministrator /
  TestSavePersistsAdministratorSession.
