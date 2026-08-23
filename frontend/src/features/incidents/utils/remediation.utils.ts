import type {
  RemediationStatus,
  ResolutionStatus,
  ValidationSummary,
  RemediationStatusConfig,
} from "../types/remediation.types";

export function isRemediationFixable(resolutionStatus?: ResolutionStatus): boolean {
  return resolutionStatus === "fixable";
}

export function isRequiresHuman(resolutionStatus?: ResolutionStatus): boolean {
  return resolutionStatus === "requires_human";
}

export function getRemediationStatusConfig(status?: RemediationStatus): RemediationStatusConfig {
  switch (status) {
    case "planned":
      return {
        label: "Planned",
        variant: "neutral",
        description: "Remediation is queued and planned for execution.",
      };
    case "in_progress":
      return {
        label: "In Progress",
        variant: "running",
        description: "Generating code fix and executing verification tests.",
      };
    case "validated":
      return {
        label: "Validated",
        variant: "success",
        description: "All regression tests and static validations passed.",
      };
    case "failed":
      return {
        label: "Failed",
        variant: "error",
        description: "Validation failed. No Pull Request was opened.",
      };
    case "pull_request_created":
      return {
        label: "Pull Request Created",
        variant: "success",
        description: "Pull Request created successfully on GitHub.",
      };
    default:
      return {
        label: "Pending",
        variant: "neutral",
        description: "Awaiting remediation lifecycle initiation.",
      };
  }
}

export function hasValidationPassed(summary?: ValidationSummary): boolean {
  if (!summary) return false;
  return summary.total > 0 && summary.passed === summary.total && summary.failed === 0;
}

export function hasValidationFailed(summary?: ValidationSummary): boolean {
  if (!summary) return false;
  return summary.failed > 0;
}
