export type TraceabilityStepId =
  | "incident"
  | "investigation"
  | "issue"
  | "remediation"
  | "branch"
  | "commit"
  | "pull_request";

export type TraceabilityStepStatus =
  | "completed"
  | "running"
  | "failed"
  | "halted"
  | "pending"
  | "not_applicable";

export interface TraceabilityStep {
  id: TraceabilityStepId;
  title: string;
  label: string;
  status: TraceabilityStepStatus;
  detail?: string;
  url?: string;
  externalLabel?: string;
  badgeText?: string;
  badgeVariant?: "neutral" | "running" | "success" | "error" | "warning";
}

