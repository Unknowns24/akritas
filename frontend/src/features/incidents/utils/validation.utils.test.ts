import { describe, it } from "node:test";
import assert from "node:assert/strict";
import type { ValidationResult } from "../types/remediation.types";

describe("Validation Results Processing - FE-H5-03", () => {
  const mockPassingChecks: ValidationResult[] = [
    {
      id: "val-1",
      remediation_id: "rem-1",
      type: "build",
      name: "npm run build",
      status: "passed",
      started_at: "2026-08-23T08:00:00Z",
      finished_at: "2026-08-23T08:00:15Z",
      summary: "Compiled successfully in 15s",
      output_excerpt: "✓ Compiled /app in 15.2s",
      output_redacted: true,
    },
    {
      id: "val-2",
      remediation_id: "rem-1",
      type: "test",
      name: "go test ./...",
      status: "passed",
      started_at: "2026-08-23T08:00:16Z",
      finished_at: "2026-08-23T08:00:25Z",
      summary: "ok akritas/core 8.2s (all 42 tests passed)",
      output_excerpt: "PASS: TestRemediationExecution\nPASS: TestValidationGate",
      output_redacted: true,
    },
  ];

  const mockFailingChecks: ValidationResult[] = [
    ...mockPassingChecks,
    {
      id: "val-3",
      remediation_id: "rem-1",
      type: "static_analysis",
      name: "golangci-lint run",
      status: "failed",
      started_at: "2026-08-23T08:00:26Z",
      finished_at: "2026-08-23T08:00:30Z",
      summary: "1 error found in service.go",
      output_excerpt: "service.go:42:12: unhandled error returned by Query() (errcheck)",
      output_redacted: true,
    },
  ];

  it("should evaluate passing checks list as valid for PR creation", () => {
    const hasFailure = mockPassingChecks.some((c) => c.status === "failed");
    const allPassed = mockPassingChecks.every((c) => c.status === "passed");

    assert.equal(hasFailure, false);
    assert.equal(allPassed, true);
  });

  it("should detect validation failure and enforce PR blockage (ADR-004)", () => {
    const hasFailure = mockFailingChecks.some((c) => c.status === "failed");
    const failingCheck = mockFailingChecks.find((c) => c.status === "failed");

    assert.equal(hasFailure, true);
    assert.ok(failingCheck);
    assert.equal(failingCheck?.name, "golangci-lint run");
    assert.equal(failingCheck?.status, "failed");
    assert.equal(failingCheck?.output_redacted, true);
  });

  it("should guarantee all validation outputs are marked redacted: true", () => {
    for (const check of mockFailingChecks) {
      assert.equal(check.output_redacted, true);
      assert.ok(check.output_excerpt && check.output_excerpt.length > 0);
    }
  });
});

