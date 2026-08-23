import React from "react";
import { Wrench, CheckCircle2, ShieldCheck, User, AlertCircle } from "lucide-react";
import type { Incident } from "../../../services/get-incident.service";
import styles from "../IncidentDetailView.module.css";

export function RemediationCard({ incident }: { incident: Incident }) {
  const requiresHuman =
    incident.latest_investigation?.resolution_status === "requires_human" || incident.terminal_outcome === "requires_human";
  const remediationFailed =
    incident.terminal_outcome === "remediation_failed";

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
              fontWeight: "600",
              marginBottom: "8px",
              color: "var(--status-error)",
              margin: "0 0 8px 0",
            }}
          >
            Remediation Failed
          </h3>
          <p
            style={{
              fontSize: "14px",
              color: "var(--text-secondary)",
              lineHeight: "1.5",
              margin: 0,
            }}
          >
            The automated agent was unable to generate a valid patch that passes validation. Intervention is required.
          </p>
        </div>
      </div>
    );
  }

  if (requiresHuman) {
    return (
      <div
        className={styles.card}
        style={{
          height: "100%",
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          textAlign: "center",
          gap: "16px",
          padding: "32px 16px",
          backgroundColor: "var(--bg-secondary)",
          border: "1px dashed var(--border-color)",
        }}
      >
        <div
          style={{
            padding: "12px",
            borderRadius: "50%",
            backgroundColor: "var(--bg-tertiary)",
            color: "var(--text-secondary)",
          }}
        >
          <User size={24} />
        </div>
        <div>
          <h3
            style={{
              fontSize: "16px",
              fontWeight: "600",
              marginBottom: "8px",
              color: "var(--text-primary)",
              margin: "0 0 8px 0",
            }}
          >
            Human Action Required
          </h3>
          <p
            style={{
              fontSize: "14px",
              color: "var(--text-secondary)",
              lineHeight: "1.5",
              margin: 0,
            }}
          >
            The automatic flow ended after creating the GitHub Issue. QVAC
            determined that this incident requires human intervention. No pull
            request will be created.
          </p>
        </div>
      </div>
    );
  }

  const remediation = incident.remediation;
  if (!remediation) return null;

  const diffLines = remediation.changes?.[0]?.patch?.split("\n") || [];

  return (
    <div
      className={styles.card}
      style={{ height: "100%", justifyContent: "space-between" }}
    >
      <div>
        <div className={styles.cardHeader}>
          <Wrench size={16} />
          Repository Analysis & Remediation
        </div>

        <div className={styles.remediationMeta} style={{ marginTop: "16px" }}>
          Files inspected:{" "}
          {incident.latest_investigation?.relevant_files?.join(", ")}
        </div>

        {diffLines.length > 0 && (
          <div className={styles.diffViewer}>
            {diffLines.map((line: string, i: number) => {
              const isMeta = line.startsWith("@@");
              const isAdded = line.startsWith("+");

              let lineClass = styles.diffLine;
              if (isMeta) lineClass += ` ${styles.diffMeta}`;
              else if (isAdded) lineClass += ` ${styles.diffAdded}`;
              else lineClass += ` ${styles.diffUnchanged}`;

              return (
                <div key={i} className={lineClass}>
                  {line}
                </div>
              );
            })}
          </div>
        )}
      </div>

      <div>
        <div className={styles.validationArea}>
          <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
            <span
              style={{
                fontFamily: "var(--font-mono)",
                fontSize: "12px",
                color: "var(--text-secondary)",
              }}
            >
              <GitMergeIcon /> akritas/fix-incident-184
            </span>
          </div>
        </div>

        <div style={{ display: "flex", gap: "12px", marginTop: "12px" }}>
          {remediation.validation_summary &&
            remediation.validation_summary.passed > 0 &&
            remediation.validation_summary.failed === 0 && (
              <div className={styles.validationTag}>
                <CheckCircle2 size={12} />
                Validation Passed
              </div>
            )}
        </div>

        <button className={styles.prButton}>Open Pull Request</button>
        <div className={styles.prDisclaimer}>
          <ShieldCheck size={12} />
          Akritas never merges changes automatically.
        </div>
      </div>
    </div>
  );
}

function GitMergeIcon() {
  return (
    <svg
      width="14"
      height="14"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      style={{ verticalAlign: "middle", marginRight: "4px" }}
    >
      <circle cx="18" cy="18" r="3"></circle>
      <circle cx="6" cy="6" r="3"></circle>
      <path d="M6 21V9a9 9 0 0 0 9 9"></path>
    </svg>
  );
}
