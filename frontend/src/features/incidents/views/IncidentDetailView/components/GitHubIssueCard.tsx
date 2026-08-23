import React from "react";
import { CircleDot, ExternalLink, AlertCircle } from "lucide-react";
import type { Incident } from "../../../services/get-incident.service";
import styles from "../IncidentDetailView.module.css";

export function GitHubIssueCard({ incident }: { incident: Incident }) {
  const issue = incident.github_issue_reference;

  const issueFailed =
    incident.terminal_outcome === "issue_publication_failed";

  if (issueFailed) {
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
            Issue Publication Failed
          </h3>
          <p
            style={{
              fontSize: "14px",
              color: "var(--text-secondary)",
              lineHeight: "1.5",
              margin: 0,
            }}
          >
            QVAC was unable to publish the tracking issue to GitHub. Check your integration settings and permissions.
          </p>
        </div>
      </div>
    );
  }

  if (!issue) {
    return null;
  }

  return (
    <div className={styles.card}>
      <div className={styles.cardHeader}>
        <CircleDot size={18} style={{ color: "var(--text-primary)" }} />
        <span>GitHub Issue</span>
      </div>
      
      <div style={{ display: "flex", flexDirection: "column", gap: "var(--space-3)", marginTop: "var(--space-2)" }}>
        <div>
          <div className={styles.contextLabel} style={{ marginBottom: "2px" }}>Repository</div>
          <div style={{ fontSize: "14px", fontWeight: "500", color: "var(--text-primary)" }}>
            {issue.repository}
          </div>
        </div>

        <div>
          <div className={styles.contextLabel} style={{ marginBottom: "2px" }}>Issue Number</div>
          <div style={{ fontSize: "14px", fontWeight: "500", color: "var(--text-primary)" }}>
            #{issue.number}
          </div>
        </div>

        <div>
          <div className={styles.contextLabel} style={{ marginBottom: "2px" }}>Created At</div>
          <div style={{ fontSize: "14px", color: "var(--text-secondary)", fontFamily: "var(--font-mono)" }}>
            {new Date(issue.created_at).toLocaleString()}
          </div>
        </div>

        <div style={{ marginTop: "var(--space-2)" }}>
          <a
            href={issue.url}
            target="_blank"
            rel="noopener noreferrer"
            className={styles.actionButton}
            style={{ display: "inline-flex", textDecoration: "none" }}
          >
            View on GitHub <ExternalLink size={14} />
          </a>
        </div>
      </div>
    </div>
  );
}
