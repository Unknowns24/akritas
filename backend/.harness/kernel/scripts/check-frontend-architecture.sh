#!/usr/bin/env bash
set -euo pipefail

if [ ! -d src ]; then
  echo "No src/ directory found. Skipping frontend architecture check."
  exit 0
fi

echo "Checking frontend architecture..."

if [ -d src/core ] && grep -R "from ['\"]@/features" src/core --include='*.ts' --include='*.tsx' >/dev/null 2>&1; then
  echo "ERROR: src/core imports src/features"
  grep -R "from ['\"]@/features" src/core --include='*.ts' --include='*.tsx' || true
  exit 1
fi

# Runtime environment reads should remain centralized. Allow conventional config/runtime-env modules.
if [ -d src/features ]; then
  if grep -R "process\.env" src/features --include='*.ts' --include='*.tsx' >/dev/null 2>&1; then
    echo "ERROR: feature code reads process.env directly; use centralized runtime configuration"
    grep -R "process\.env" src/features --include='*.ts' --include='*.tsx' || true
    exit 1
  fi
fi

# Feature code should not instantiate its own HTTP client.
if [ -d src/features ]; then
  if grep -R "axios\.create" src/features --include='*.ts' --include='*.tsx' >/dev/null 2>&1; then
    echo "ERROR: feature code creates an ad-hoc Axios client; use the shared API client"
    grep -R "axios\.create" src/features --include='*.ts' --include='*.tsx' || true
    exit 1
  fi
fi

echo "Frontend architecture check passed."
