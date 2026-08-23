import React from "react";
import { CircleDot, ExternalLink } from "lucide-react";
import type { Incident } from "../../../services/get-incident.service";
import styles from "../IncidentDetailView.module.css";

export function GitHubIssueCard({ incident }: { incident: Incident }) {
  const issue = incident.github_issue_reference;

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
