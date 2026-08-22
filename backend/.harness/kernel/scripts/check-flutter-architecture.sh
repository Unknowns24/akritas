#!/usr/bin/env bash
set -euo pipefail

if [ ! -d lib ]; then
  echo "No lib/ directory found. Skipping Flutter architecture check."
  exit 0
fi

echo "Checking Flutter architecture..."

if [ -d lib/features ]; then
  presentation_files=$(find lib/features -path '*/presentation/*' -name '*.dart' -type f || true)
  if [ -n "$presentation_files" ]; then
    if echo "$presentation_files" | xargs grep -n "core/api\|ApiClient\|package:dio/dio.dart\|/data/dto/" >/dev/null 2>&1; then
      echo "ERROR: presentation layer imports ApiClient/Dio/DTOs or data internals."
      echo "$presentation_files" | xargs grep -n "core/api\|ApiClient\|package:dio/dio.dart\|/data/dto/" || true
      exit 1
    fi
  fi

  domain_files=$(find lib/features -path '*/domain/*' -name '*.dart' -type f || true)
  if [ -n "$domain_files" ]; then
    if echo "$domain_files" | xargs grep -n "package:flutter/\|package:dio/\|/data/\|/presentation/" >/dev/null 2>&1; then
      echo "ERROR: domain layer imports Flutter, Dio, data or presentation."
      echo "$domain_files" | xargs grep -n "package:flutter/\|package:dio/\|/data/\|/presentation/" || true
      exit 1
    fi
  fi

  if grep -R "Dio(" lib/features --include='*.dart' >/dev/null 2>&1; then
    echo "ERROR: Dio() instantiated inside features. Use the shared ApiClient/provider."
    grep -R "Dio(" lib/features --include='*.dart' || true
    exit 1
  fi
fi

echo "Flutter architecture check passed."
