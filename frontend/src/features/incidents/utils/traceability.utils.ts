import type { Incident } from "../services/get-incident.service";
import type { TraceabilityStep } from "../types/traceability.types";

export function buildIncidentTraceabilityChain(incident: Incident): TraceabilityStep[] {
  const steps: TraceabilityStep[] = [];
  const investigation = incident.latest_investigation;
  const issue = incident.github_issue_reference;
  const remediation = incident.remediation;
  const pr = incident.pull_request_reference || remediation?.pull_request_reference;
  const isFixable = incident.resolution_status === "fixable";
  const isHuman = incident.resolution_status === "requires_human";
  const hasRemediationFailed = remediation?.status === "failed";
  const hasPR = Boolean(pr?.url);

  // 1. Incident Node
  steps.push({
    id: "incident",
    title: "Incident Detected",
    label: incident.key,
    status: "completed",
    detail: incident.title,
    badgeText: incident.severity?.toUpperCase(),
    badgeVariant: incident.severity === "error" ? "error" : "warning",
  });

  // 2. Investigation Node
  if (investigation) {
    const isInvDone = investigation.status === "completed";
    steps.push({
      id: "investigation",
      title: "Root Cause Investigation",
      label: investigation.root_cause_status
        ? `Cause: ${investigation.root_cause_status}`
        : "Investigating",
      status: isInvDone ? "completed" : investigation.status === "failed" ? "failed" : "running",
      detail: investigation.summary || investigation.root_cause,
      badgeText: isInvDone ? "Analyzed" : "In Progress",
      badgeVariant: isInvDone ? "success" : "running",
    });
  } else {
    steps.push({
      id: "investigation",
      title: "Root Cause Investigation",
      label: "Pending Investigation",
      status: "pending",
      detail: "Waiting for agent investigation to start",
      badgeText: "Pending",
      badgeVariant: "neutral",
    });
  }

  // 3. GitHub Issue Node
  if (issue) {
    steps.push({
      id: "issue",
      title: "GitHub Issue",
      label: `Issue #${issue.number}`,
      status: "completed",
      detail: issue.repository,
      url: issue.url,
      externalLabel: "View Issue",
      badgeText: "Documented",
      badgeVariant: "success",
    });
  } else {
    steps.push({
      id: "issue",
      title: "GitHub Issue",
      label: "Issue Tracking",
      status: investigation?.status === "completed" ? "running" : "pending",
      detail: "Creating tracking issue in repository",
      badgeText: "Pending",
      badgeVariant: "neutral",
    });
  }

  // 4. Remediation Node
  if (isHuman) {
    steps.push({
      id: "remediation",
      title: "Remediation Boundary",
      label: "Requires Human",
      status: "halted",
      detail: "Automated remediation halted; engineer intervention required",
      badgeText: "Halted",
      badgeVariant: "warning",
    });
  } else if (remediation) {
    const isRemDone = remediation.status === "validated" || remediation.status === "pull_request_created";
    steps.push({
      id: "remediation",
      title: "Autonomous Remediation",
      label: `Status: ${remediation.status}`,
      status: hasRemediationFailed
        ? "failed"
        : isRemDone
        ? "completed"
        : remediation.status === "in_progress"
        ? "running"
        : "pending",
      detail: remediation.changes_summary || "Autonomous patch and verification cycle",
      badgeText: remediation.status.toUpperCase(),
      badgeVariant: hasRemediationFailed ? "error" : isRemDone ? "success" : "running",
    });
  } else if (isFixable) {
    steps.push({
      id: "remediation",
      title: "Autonomous Remediation",
      label: "Planned Remediation",
      status: "pending",
      detail: "Queued for workspace execution",
      badgeText: "Planned",
      badgeVariant: "neutral",
    });
  } else {
    steps.push({
      id: "remediation",
      title: "Autonomous Remediation",
      label: "Pending Classification",
      status: "pending",
      detail: "Waiting for fixable classification",
      badgeText: "Pending",
      badgeVariant: "neutral",
    });
  }

  // 5. Branch Node
  const branchName =
    remediation?.branch_name || pr?.branch || (isFixable ? `akritas/fix-${incident.key.toLowerCase()}` : undefined);
  if (isHuman) {
    steps.push({
      id: "branch",
      title: "Remediation Branch",
      label: "No Branch Created",
      status: "not_applicable",
      detail: "Branch creation skipped for manual incidents",
      badgeText: "N/A",
      badgeVariant: "neutral",
    });
  } else if (branchName && (remediation?.status === "in_progress" || remediation?.status === "validated" || hasPR || hasRemediationFailed)) {
    steps.push({
      id: "branch",
      title: "Remediation Branch",
      label: branchName,
      status: "completed",
      detail: "Isolated Git workspace branch",
      badgeText: "Created",
      badgeVariant: "success",
    });
  } else {
    steps.push({
      id: "branch",
      title: "Remediation Branch",
      label: isFixable ? `akritas/fix-${incident.key.toLowerCase()}` : "Branch Creation",
      status: "pending",
      detail: "Branch will be created during execution",
      badgeText: "Pending",
      badgeVariant: "neutral",
    });
  }

  // 6. Commit Node
  const commitSha =
    incident.deployment_correlation?.commit_sha ||
    investigation?.relevant_commits?.[0] ||
    (hasPR ? "validated-fix" : undefined);

  if (isHuman) {
    steps.push({
      id: "commit",
      title: "Validated Commit",
      label: "No Commit Created",
      status: "not_applicable",
      detail: "Commit generation skipped for manual incidents",
      badgeText: "N/A",
      badgeVariant: "neutral",
    });
  } else if (hasRemediationFailed) {
    steps.push({
      id: "commit",
      title: "Validated Commit",
      label: "Commit Blocked",
      status: "failed",
      detail: "Unvalidated code changes were not committed (ADR-004)",
      badgeText: "Blocked",
      badgeVariant: "error",
    });
  } else if (hasPR || remediation?.status === "validated") {
    steps.push({
      id: "commit",
      title: "Validated Commit",
      label: commitSha ? `Commit: ${commitSha.slice(0, 7)}` : "Verified Commit",
      status: "completed",
      detail: "Code fixes and regression tests committed",
      badgeText: "Committed",
      badgeVariant: "success",
    });
  } else {
    steps.push({
      id: "commit",
      title: "Validated Commit",
      label: "Awaiting Validation",
      status: "pending",
      detail: "Commit is generated after tests pass",
      badgeText: "Pending",
      badgeVariant: "neutral",
    });
  }

  // 7. Pull Request Node
  if (isHuman) {
    steps.push({
      id: "pull_request",
      title: "GitHub Pull Request",
      label: "No PR Created",
      status: "not_applicable",
      detail: "Incident requires manual human intervention",
      badgeText: "N/A",
      badgeVariant: "neutral",
    });
  } else if (hasRemediationFailed) {
    steps.push({
      id: "pull_request",
      title: "GitHub Pull Request",
      label: "PR Blocked",
      status: "failed",
      detail: "Validation failed; PR creation halted by safety policy",
      badgeText: "Blocked",
      badgeVariant: "error",
    });
  } else if (hasPR && pr) {
    steps.push({
      id: "pull_request",
      title: "GitHub Pull Request",
      label: `PR #${pr.number}`,
      status: "completed",
      detail: pr.repository,
      url: pr.url,
      externalLabel: "View Pull Request",
      badgeText: "Opened",
      badgeVariant: "success",
    });
  } else {
    steps.push({
      id: "pull_request",
      title: "GitHub Pull Request",
      label: "Pull Request Creation",
      status: "pending",
      detail: "Will open PR once validation succeeds",
      badgeText: "Pending",
      badgeVariant: "neutral",
    });
  }

  return steps;
}

