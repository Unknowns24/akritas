import React from "react";
import {
  AlertCircle,
  Clock,
  ExternalLink,
  FileSearch,
  GitBranch,
  Loader2,
  ShieldCheck,
  Wrench,
} from "lucide-react";
import type { Incident } from "../../../services/get-incident.service";
import {
  isRemediationFixable,
  isRequiresHuman,
} from "../../../utils/remediation.utils";
import { RemediationStatusBadge } from "./RemediationStatusBadge";
import { ValidationSummaryView } from "./ValidationSummaryView";
import { ValidationResultsViewer } from "./ValidationResultsViewer";
import { CodeChangesDiffViewer } from "./CodeChangesDiffViewer";
import { RequiresHumanCard } from "./RequiresHumanCard";
import { AutonomyBoundaryBanner } from "./AutonomyBoundaryBanner";
import styles from "./RemediationCard.module.css";

interface RemediationCardProps {
  incident: Incident;
}

export function RemediationCard({ incident }: RemediationCardProps) {
  const resolutionStatus =
    incident.resolution_status ??
    incident.latest_investigation?.resolution_status;

  const remediationFailed =
    incident.terminal_outcome === "remediation_failed";

  const requiresHuman =
    isRequiresHuman(resolutionStatus) ||
    incident.terminal_outcome === "requires_human";

  if (remediationFailed) {
    return (
      <div
        className={styles.card}
        style={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          textAlign: "center",
          gap: "16px",
          padding: "32px 16px",
          backgroundColor: "rgba(var(--status-error-rgb), 0.1)",
          border: "1px dashed var(--status-error)",
        }}
      >
        <div
          style={{
            padding: "12px",
            borderRadius: "50%",
            backgroundColor: "rgba(var(--status-error-rgb), 0.2)",
            color: "var(--status-error)",
          }}
        >
          <AlertCircle size={24} />
        </div>

        <div>
          <h3
            style={{
              fontSize: "16px",
              fontWeight: 600,
              margin: "0 0 8px",
              color: "var(--status-error)",
            }}
          >
            Remediation Failed
          </h3>

          <p
            style={{
              fontSize: "14px",
              color: "var(--text-secondary)",
              lineHeight: 1.5,
              margin: 0,
            }}
          >
            The automated agent was unable to generate a valid patch that
            passes validation. Intervention is required.
          </p>
        </div>
      </div>
    );
  }

  if (requiresHuman) {
    return <RequiresHumanCard incident={incident} />;
  }

  if (!isRemediationFixable(resolutionStatus)) {
    return (
      <div className={styles.pendingCard}>
        <FileSearch size={28} style={{ color: "var(--text-dim)" }} />

        <div>
          <div
            style={{
              fontWeight: 600,
              color: "var(--text-primary)",
              marginBottom: "4px",
            }}
          >
            Remediation Pending
          </div>

          <div>
            Remediation evaluation will start once the incident investigation
            completes.
          </div>
        </div>
      </div>
    );
  }

  const remediation = incident.remediation;
  const status = remediation?.status || "planned";

  const branchName =
    remediation?.branch_name ||
    incident.pull_request_reference?.branch ||
    `akritas/fix-${incident.key.toLowerCase()}`;

  const prRef =
    incident.pull_request_reference ||
    remediation?.pull_request_reference;

  const relevantFiles =
    incident.latest_investigation?.relevant_files || [];

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

        <div className={styles.metaRow}>
          <span className={styles.metaLabel}>
            Target Remediation Branch
          </span>

          <div className={styles.branchBox}>
            <GitBranch size={14} className={styles.branchIcon} />
            <span>{branchName}</span>
          </div>
        </div>

        {relevantFiles.length > 0 && (
          <div className={styles.metaRow}>
            <span className={styles.metaLabel}>Relevant Files</span>
            <span className={styles.metaValue}>
              {relevantFiles.join(", ")}
            </span>
          </div>
        )}

        {remediation?.changes_summary && (
          <div className={styles.metaRow}>
            <span className={styles.metaLabel}>Changes Summary</span>

            <span
              style={{
                fontSize: "13px",
                color: "var(--text-secondary)",
                lineHeight: 1.5,
              }}
            >
              {remediation.changes_summary}
            </span>
          </div>
        )}

        {status === "planned" &&
          (!remediation?.changes || remediation.changes.length === 0) && (
            <div className={styles.plannedBox}>
              <Clock
                size={16}
                style={{
                  color: "var(--text-dim)",
                  flexShrink: 0,
                }}
              />

              <span>
                Remediation workspace is queued. Code fixes and regression tests
                will be generated automatically.
              </span>
            </div>
          )}

        {status === "in_progress" && (
          <div className={styles.inProgressBox}>
            <Loader2
              size={16}
              className={styles.spin}
              style={{
                color: "#60a5fa",
                flexShrink: 0,
              }}
            />

            <span>
              Generating patch, regression tests, and executing build validation
              checks...
            </span>
          </div>
        )}

        {remediation?.changes && remediation.changes.length > 0 && (
          <div className={styles.metaRow}>
            <span className={styles.metaLabel}>
              Proposed Code Changes ({remediation.changes.length}{" "}
              {remediation.changes.length === 1 ? "file" : "files"})
            </span>

            <CodeChangesDiffViewer changes={remediation.changes} />
          </div>
        )}

        <ValidationSummaryView
          summary={remediation?.validation_summary}
        />

        <ValidationResultsViewer
          remediationStatus={status}
          failureUserMessage={failureMessage}
        />

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

        {(status === "pull_request_created" || Boolean(prRef)) && (
          <AutonomyBoundaryBanner />
        )}
      </div>

      <div className={styles.footer}>
        <ShieldCheck size={14} />
        <span>
          Akritas never merges changes automatically. Human review is required.
        </span>
      </div>
    </div>
  );
}