# Implementation brief — AKR-QVAC-TOOL-LOOP

## Estrategia

```text
Runner.Run:
  messages = [system, user]
  loop (maxIterations):
    resp = ChatCompletions(messages, tools=allowlistedDefs)  # sin response_format
    if tool_calls:
      for each call:
        if not allowlisted -> tool error message (no crash)
        else execute -> append role=tool
      continue
    break
  final = ChatCompletions(messages, response_format=json_schema)  # sin tools
  return parseResult(final)
```

- `Tool` interface: Name, Description, Parameters schema, Execute(ctx, argsJSON) (string, error)
- Registry map allowlist; unknown tool → resultado de error controlado hacia el modelo
- Max tool rounds / max tool calls configurables (defaults seguros, p.ej. 8 rounds)
- System prompt: enmarcar cualquier contenido de logs/código como DATA, no instructions
