import React from "react";
import {
  Wrench,
  GitBranch,
  ExternalLink,
  ShieldCheck,
  AlertTriangle,
  Clock,
  Loader2,
  FileSearch,
} from "lucide-react";
import type { Incident } from "../../../services/get-incident.service";
import { isRemediationFixable, isRequiresHuman } from "../../../utils/remediation.utils";
import { RemediationStatusBadge } from "./RemediationStatusBadge";
import { ValidationSummaryView } from "./ValidationSummaryView";
import { CodeChangesDiffViewer } from "./CodeChangesDiffViewer";
import { RequiresHumanCard } from "./RequiresHumanCard";
import styles from "./RemediationCard.module.css";

interface RemediationCardProps {
  incident: Incident;
}

export function RemediationCard({ incident }: RemediationCardProps) {
  const resolutionStatus = incident.resolution_status;

  // 1. Boundary: If requires_human, render dedicated manual intervention card
  if (isRequiresHuman(resolutionStatus)) {
    return <RequiresHumanCard incident={incident} />;
  }

  // 2. If resolution status is not yet fixable (e.g. investigation pending or unknown)
  if (!isRemediationFixable(resolutionStatus)) {
    return (
      <div className={styles.pendingCard}>
        <FileSearch size={28} style={{ color: "var(--text-dim)" }} />
        <div>
          <div style={{ fontWeight: 600, color: "var(--text-primary)", marginBottom: "4px" }}>
            Remediation Pending
          </div>
          <div>
            Remediation evaluation will start once the incident investigation completes.
          </div>
        </div>
      </div>
    );
  }

  // 3. Fixable Incident: Render full Remediation lifecycle
  const remediation = incident.remediation;
  const status = remediation?.status || "planned";
  const branchName =
    remediation?.branch_name ||
    incident.pull_request_reference?.branch ||
    `akritas/fix-${incident.key.toLowerCase()}`;
  const prRef = incident.pull_request_reference || remediation?.pull_request_reference;
  const relevantFiles = incident.latest_investigation?.relevant_files || [];
  const failureMessage = remediation?.failure_user_message;

  return (
    <div className={styles.card}>
      <div className={styles.content}>
        <div className={styles.header}>
          <div className={styles.titleGroup}>
            <Wrench size={18} style={{ color: "#60a5fa" }} />
            <span>Autonomous Remediation</span>
          </div>
          <RemediationStatusBadge status={status} />
        </div>

        {/* Branch Reference */}
        <div className={styles.metaRow}>
          <span className={styles.metaLabel}>Target Remediation Branch</span>
          <div className={styles.branchBox}>
            <GitBranch size={14} className={styles.branchIcon} />
            <span>{branchName}</span>
          </div>
        </div>

        {/* Files Inspected / Touched */}
        {relevantFiles.length > 0 && (
          <div className={styles.metaRow}>
            <span className={styles.metaLabel}>Relevant Files</span>
            <span className={styles.metaValue}>{relevantFiles.join(", ")}</span>
          </div>
        )}

        {/* Changes Summary if available */}
        {remediation?.changes_summary && (
          <div className={styles.metaRow}>
            <span className={styles.metaLabel}>Changes Summary</span>
            <span style={{ fontSize: "13px", color: "var(--text-secondary)", lineHeight: 1.5 }}>
              {remediation.changes_summary}
            </span>
          </div>
        )}

        {/* Status: Planned */}
        {status === "planned" && (!remediation?.changes || remediation.changes.length === 0) && (
          <div className={styles.plannedBox}>
            <Clock size={16} style={{ color: "var(--text-dim)", flexShrink: 0 }} />
            <span>
              Remediation workspace is queued. Code fixes and regression tests will be generated automatically.
            </span>
          </div>
        )}

        {/* Status: In Progress */}
        {status === "in_progress" && (
          <div className={styles.inProgressBox}>
            <Loader2 size={16} className={styles.spin} style={{ color: "#60a5fa", flexShrink: 0 }} />
            <span>
              Generating patch, regression tests, and executing build validation checks...
            </span>
          </div>
        )}

        {/* Status: Failed */}
        {status === "failed" && (
          <div className={styles.failureAlert}>
            <div className={styles.failureAlertHeader}>
              <AlertTriangle size={16} />
              <span>Remediation Failed</span>
            </div>
            <div>{failureMessage || "Validation checks failed during the remediation process."}</div>
            <div className={styles.failureNote}>
              In accordance with safety policies (ADR-004), no Pull Request was opened for unvalidated changes.
            </div>
          </div>
        )}

        {/* Multi-file Code Changes Diff Viewer */}
        {remediation?.changes && remediation.changes.length > 0 && (
          <div className={styles.metaRow}>
            <span className={styles.metaLabel}>Proposed Code Changes ({remediation.changes.length} file{remediation.changes.length > 1 ? "s" : ""})</span>
            <CodeChangesDiffViewer changes={remediation.changes} />
          </div>
        )}

        {/* Automated Validation Summary */}
        <ValidationSummaryView summary={remediation?.validation_summary} />

        {/* Pull Request section if created */}
        {prRef && (
          <div className={styles.prSection}>
            <div className={styles.prDetails}>
              <span>GitHub Pull Request</span>
              <span className={styles.prNumber}>PR #{prRef.number}</span>
            </div>
            <a
              href={prRef.url}
              target="_blank"
              rel="noreferrer"
              className={styles.prButton}
            >
              <span>View Pull Request #{prRef.number}</span>
              <ExternalLink size={14} />
            </a>
          </div>
        )}
      </div>

      <div className={styles.footer}>
        <ShieldCheck size={14} />
        <span>Akritas never merges changes automatically. Human review is required.</span>
      </div>
    </div>
  );
}
