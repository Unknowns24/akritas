#!/usr/bin/env bash
set -euo pipefail

echo "Running lightweight security checks..."

if grep -R "AKIA[0-9A-Z]\{16\}\|BEGIN RSA PRIVATE KEY\|BEGIN OPENSSH PRIVATE KEY" . \
  --exclude=check-security.sh \
  --exclude-dir=.git --exclude-dir=node_modules --exclude-dir=.next --exclude-dir=.harness/runtime >/dev/null 2>&1; then
  echo "ERROR: potential secret detected"
  grep -R "AKIA[0-9A-Z]\{16\}\|BEGIN RSA PRIVATE KEY\|BEGIN OPENSSH PRIVATE KEY" . \
    --exclude=check-security.sh \
    --exclude-dir=.git --exclude-dir=node_modules --exclude-dir=.next --exclude-dir=.harness/runtime || true
  exit 1
fi

echo "Security check passed."
