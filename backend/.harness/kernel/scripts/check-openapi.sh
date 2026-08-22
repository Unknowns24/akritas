#!/usr/bin/env bash
set -euo pipefail

SPEC=""
for f in internal/adapter/rest/swagger/openapi.yaml internal/adapter/rest/docs/openapi.yaml docs/openapi.yaml openapi.yaml; do
  if [ -f "$f" ]; then SPEC="$f"; break; fi
done

if [ -z "$SPEC" ]; then
  echo "No OpenAPI spec found. Skipping OpenAPI check."
  exit 0
fi

echo "Found OpenAPI spec: $SPEC"
python3 - <<PY
import sys, pathlib
try:
    import yaml
except Exception:
    print('PyYAML not installed; basic file existence check only.')
    sys.exit(0)
path = pathlib.Path('$SPEC')
with path.open('r', encoding='utf-8') as f:
    data = yaml.safe_load(f)
if not isinstance(data, dict):
    raise SystemExit('OpenAPI YAML is not an object')
if 'openapi' not in data:
    raise SystemExit('Missing openapi field')
if 'info' not in data or 'version' not in data['info']:
    raise SystemExit('Missing info.version')
print('OpenAPI check passed.')
PY
