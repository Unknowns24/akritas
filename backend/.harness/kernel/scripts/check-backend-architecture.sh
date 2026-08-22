#!/usr/bin/env bash
set -euo pipefail

if [ ! -d internal ]; then
  echo "No internal/ directory found. Skipping backend architecture check."
  exit 0
fi

echo "Checking backend architecture imports..."

if [ -d internal/core ] && grep -R "internal/adapter" internal/core --include='*.go' >/dev/null 2>&1; then
  echo "ERROR: internal/core imports internal/adapter"
  grep -R "internal/adapter" internal/core --include='*.go' || true
  exit 1
fi

CHECK_DIRS=()
[ -d internal/core ] && CHECK_DIRS+=(internal/core)
[ -d internal/usecase ] && CHECK_DIRS+=(internal/usecase)

if [ ${#CHECK_DIRS[@]} -gt 0 ]; then
  PATTERN='gorm.io|github.com/go-chi/chi|net/http|os/exec'
  if grep -RE "$PATTERN" "${CHECK_DIRS[@]}" --include='*.go' >/dev/null 2>&1; then
    echo "ERROR: core/usecase imports infrastructure packages (declarative gorm tags are allowed; gorm packages are not)"
    grep -RE "$PATTERN" "${CHECK_DIRS[@]}" --include='*.go' || true
    exit 1
  fi
fi

if [ -d internal/core ] && grep -RE '"[0-9A-F]x[123][0-9A-F]{5}[VUFNCI]"' internal/core --include='*.go' >/dev/null 2>&1; then
  echo "ERROR: core declares REST/DB/external-adapter error codes"
  grep -RE '"[0-9A-F]x[123][0-9A-F]{5}[VUFNCI]"' internal/core --include='*.go' || true
  exit 1
fi

if [ -d internal/adapter/rest/dto ]; then
  root_dtos=$(find internal/adapter/rest/dto -maxdepth 1 -type f -name '*.go' ! -name '*_test.go' | sort)
  if [ -n "$root_dtos" ]; then
    echo "ERROR: REST DTOs must be grouped under dto/<feature> or dto/common"
    printf '%s\n' "$root_dtos"
    exit 1
  fi
  while IFS= read -r file; do
    names=$(sed -nE 's/^type ([A-Za-z0-9_]+)(\[[^]]+\])? struct.*/\1/p' "$file")
    count=$(printf '%s\n' "$names" | sed '/^$/d' | wc -l | tr -d ' ')
    if [ "$count" -gt 1 ]; then
      echo "ERROR: REST DTO file declares more than one structure: $file"
      exit 1
    fi
    for name in $names; do
      case "$name" in
        *DTO) ;;
        *) echo "ERROR: REST contract structure must use DTO suffix: $file ($name)"; exit 1 ;;
      esac
    done
  done < <(find internal/adapter/rest/dto -type f -name '*.go' ! -name '*_test.go' | sort)
fi

if [ -d internal/adapter/rest/mapper ]; then
  while IFS= read -r file; do
    count=$(sed -nE 's/^func ([A-Z][A-Za-z0-9_]*)\(.*/\1/p' "$file" | wc -l | tr -d ' ')
    if [ "$count" -gt 1 ]; then
      echo "ERROR: REST mapper file exposes more than one mapping responsibility: $file"
      exit 1
    fi
  done < <(find internal/adapter/rest/mapper -type f -name '*.go' ! -name '*_test.go' | sort)
fi

echo "Backend architecture check passed."
