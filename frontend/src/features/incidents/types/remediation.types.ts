import type { components } from "@/core/libs/api-client";

export type Remediation = components["schemas"]["Remediation"];
export type RemediationStatus = components["schemas"]["RemediationStatus"];
export type CodeChange = components["schemas"]["CodeChange"];
export type ValidationSummary = components["schemas"]["Remediation"]["validation_summary"];
export type PullRequestReference = components["schemas"]["PullRequestReference"];
export type ResolutionStatus = components["schemas"]["ResolutionStatus"];
export type RootCauseStatus = components["schemas"]["RootCauseStatus"];

export type RemediationStatusVariant = "neutral" | "running" | "success" | "error";

export interface RemediationStatusConfig {
  label: string;
  variant: RemediationStatusVariant;
  description: string;
}
