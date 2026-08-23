import { describe, it } from "node:test";
import assert from "node:assert/strict";
import {
  isRemediationFixable,
  isRequiresHuman,
  getRemediationStatusConfig,
  hasValidationPassed,
  hasValidationFailed,
} from "./remediation.utils.ts";
describe("Remediation Utilities - FE-H5-01", () => {
  it("should identify fixable vs requires_human resolution status correctly", () => {
    assert.equal(isRemediationFixable("fixable"), true);
    assert.equal(isRemediationFixable("requires_human"), false);
    assert.equal(isRemediationFixable(undefined), false);

    assert.equal(isRequiresHuman("requires_human"), true);
    assert.equal(isRequiresHuman("fixable"), false);
    assert.equal(isRequiresHuman(undefined), false);
  });

  it("should return valid status configuration for all OpenAPI RemediationStatus values", () => {
    const planned = getRemediationStatusConfig("planned");
    assert.equal(planned.label, "Planned");
    assert.equal(planned.variant, "neutral");

    const inProgress = getRemediationStatusConfig("in_progress");
    assert.equal(inProgress.label, "In Progress");
    assert.equal(inProgress.variant, "running");

    const validated = getRemediationStatusConfig("validated");
    assert.equal(validated.label, "Validated");
    assert.equal(validated.variant, "success");

    const failed = getRemediationStatusConfig("failed");
    assert.equal(failed.label, "Failed");
    assert.equal(failed.variant, "error");

    const prCreated = getRemediationStatusConfig("pull_request_created");
    assert.equal(prCreated.label, "Pull Request Created");
    assert.equal(prCreated.variant, "success");

    const unknown = getRemediationStatusConfig(undefined);
    assert.equal(unknown.label, "Pending");
    assert.equal(unknown.variant, "neutral");
  });

  it("should evaluate validation summary correctly", () => {
    assert.equal(hasValidationPassed(undefined), false);
    assert.equal(hasValidationPassed({ total: 0, passed: 0, failed: 0 }), false);
    assert.equal(hasValidationPassed({ total: 3, passed: 3, failed: 0 }), true);
    assert.equal(hasValidationPassed({ total: 3, passed: 2, failed: 1 }), false);

    assert.equal(hasValidationFailed(undefined), false);
    assert.equal(hasValidationFailed({ total: 3, passed: 3, failed: 0 }), false);
    assert.equal(hasValidationFailed({ total: 3, passed: 2, failed: 1 }), true);
  });
});
