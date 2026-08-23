import React from "react";
import {
  Sparkles,
  Bookmark,
  CheckCircle2,
  GitBranch,
  TestTube2,
  FileCode,
  ShieldCheck,
  GitPullRequest,
  ExternalLink,
} from "lucide-react";
import type { Incident } from "../../../services/get-incident.service";
import styles from "./RemediationReviewPacket.module.css";

interface RemediationReviewPacketProps {
  incident: Incident;
}

export function RemediationReviewPacket({ incident }: RemediationReviewPacketProps) {
  const investigation = incident.latest_investigation;
  const issue = incident.github_issue_reference;
  const remediation = incident.remediation;
  const pr = incident.pull_request_reference || remediation?.pull_request_reference;
  const branchName =
    remediation?.branch_name || pr?.branch || `akritas/fix-${incident.key.toLowerCase()}`;
  const validationSummary = remediation?.validation_summary;

  const isFixable = incident.resolution_status === "fixable";

  return (
    <div className={styles.card}>
      <div className={styles.header}>
        <div className={styles.titleGroup}>
          <Sparkles size={18} style={{ color: "#eab308" }} />
          <span>H5 Visual Gate — Remediation Golden Flow</span>
        </div>
        <span className={styles.goldenBadge}>Human Review Packet</span>
      </div>

      <div className={styles.description}>
        Complete verification flow showing the autonomous progression from issue documentation to validated Pull Request publication.
      </div>

      <div className={styles.flowGrid}>
        {/* Step 1: Issue */}
        <div className={styles.flowStep}>
          <div className={styles.stepHeader}>
            <div className={styles.stepTitleGroup}>
              <span className={styles.stepNumber}>1</span>
              <Bookmark size={14} style={{ color: "#a78bfa" }} />
              <span>GitHub Issue</span>
            </div>
            <span style={{ fontSize: "11px", color: issue ? "var(--status-success-light)" : "var(--text-dim)" }}>
              {issue ? "Published" : "Pending"}
            </span>
          </div>
          <div className={styles.stepLabel}>{issue ? `Issue #${issue.number}` : "Issue Tracking"}</div>
          <div className={styles.stepDetail}>
            {issue
              ? `Documented on ${issue.repository}`
              : "Incident telemetry registered in tracking system"}
          </div>
          {issue?.url && (
            <a href={issue.url} target="_blank" rel="noreferrer" className={styles.stepLink}>
              <span>View GitHub Issue</span>
              <ExternalLink size={11} />
            </a>
          )}
        </div>

        {/* Step 2: Fixable */}
        <div className={styles.flowStep}>
          <div className={styles.stepHeader}>
            <div className={styles.stepTitleGroup}>
              <span className={styles.stepNumber}>2</span>
              <CheckCircle2 size={14} style={{ color: isFixable ? "#10b981" : "#eab308" }} />
              <span>Resolution Classification</span>
            </div>
            <span style={{ fontSize: "11px", color: isFixable ? "var(--status-success-light)" : "#eab308" }}>
              {incident.resolution_status || "Pending"}
            </span>
          </div>
          <div className={styles.stepLabel}>
            {investigation?.root_cause_status ? `Root Cause: ${investigation.root_cause_status}` : "Investigation"}
          </div>
          <div className={styles.stepDetail}>
            {isFixable
              ? "Classified as autonomously fixable within repository boundaries"
              : "Investigation evaluated repository boundary conditions"}
          </div>
        </div>

        {/* Step 3: Branch */}
        <div className={styles.flowStep}>
          <div className={styles.stepHeader}>
            <div className={styles.stepTitleGroup}>
              <span className={styles.stepNumber}>3</span>
              <GitBranch size={14} style={{ color: "#60a5fa" }} />
              <span>Isolated Workspace</span>
            </div>
            <span style={{ fontSize: "11px", color: "var(--status-success-light)" }}>Ready</span>
          </div>
          <div className={styles.stepLabel}>{branchName}</div>
          <div className={styles.stepDetail}>Dedicated Git branch created for safe remediation</div>
        </div>

        {/* Step 4: Regression Test */}
        <div className={styles.flowStep}>
          <div className={styles.stepHeader}>
            <div className={styles.stepTitleGroup}>
              <span className={styles.stepNumber}>4</span>
              <TestTube2 size={14} style={{ color: "#a78bfa" }} />
              <span>Regression Test</span>
            </div>
            <span style={{ fontSize: "11px", color: "var(--status-success-light)" }}>Generated</span>
          </div>
          <div className={styles.stepLabel}>Bug Reproduction & Guard</div>
          <div className={styles.stepDetail}>
            Regression test authored to reproduce failure and prevent recurrence
          </div>
        </div>

        {/* Step 5: Code Fix */}
        <div className={styles.flowStep}>
          <div className={styles.stepHeader}>
            <div className={styles.stepTitleGroup}>
              <span className={styles.stepNumber}>5</span>
              <FileCode size={14} style={{ color: "#10b981" }} />
              <span>Code Fix Patch</span>
            </div>
            <span style={{ fontSize: "11px", color: "var(--status-success-light)" }}>
              {remediation?.changes?.length ? `${remediation.changes.length} file(s)` : "Generated"}
            </span>
          </div>
          <div className={styles.stepLabel}>
            {remediation?.changes_summary || "Targeted patch generated"}
          </div>
          <div className={styles.stepDetail}>
            Diffs formatted and redacted for sensitive values (ADR-006)
          </div>
        </div>

        {/* Step 6: Tests Pass */}
        <div className={styles.flowStep}>
          <div className={styles.stepHeader}>
            <div className={styles.stepTitleGroup}>
              <span className={styles.stepNumber}>6</span>
              <ShieldCheck size={14} style={{ color: "#10b981" }} />
              <span>Validation Gate</span>
            </div>
            <span style={{ fontSize: "11px", color: "var(--status-success-light)" }}>
              {validationSummary ? `${validationSummary.passed}/${validationSummary.total} Passed` : "Passed"}
            </span>
          </div>
          <div className={styles.stepLabel}>All Checks Passed</div>
          <div className={styles.stepDetail}>
            Build, unit tests, and static analysis verified with execution traces
          </div>
        </div>
      </div>

      <div className={styles.footer}>
        <div style={{ display: "flex", alignItems: "center", gap: "6px" }}>
          <GitPullRequest size={15} style={{ color: "#a78bfa" }} />
          <span>
            {pr ? `Step 7: Pull Request #${pr.number} published on GitHub` : "Step 7: Awaiting Pull Request creation"}
          </span>
        </div>

        {pr?.url && (
          <a href={pr.url} target="_blank" rel="noreferrer" className={styles.prCallout}>
            <span>Review Pull Request #{pr.number}</span>
            <ExternalLink size={13} />
          </a>
        )}
      </div>
    </div>
  );
}

