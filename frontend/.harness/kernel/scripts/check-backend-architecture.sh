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
    echo "ERROR: core/usecase imports infrastructure packages"
    grep -RE "$PATTERN" "${CHECK_DIRS[@]}" --include='*.go' || true
    exit 1
  fi
fi

echo "Backend architecture check passed."
