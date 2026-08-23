import React from "react";
import { Sparkles, Loader2, AlertCircle, AlertTriangle } from "lucide-react";
import type { Incident } from "../../../services/get-incident.service";
import styles from "../IncidentDetailView.module.css";

export function RootCauseCard({ incident }: { incident: Incident }) {
  const investigation = incident.latest_investigation;
  if (!investigation) return null;

  const renderStatus = () => {
    switch (investigation.status) {
      case "pending":
      case "running":
        return (
          <div className={`${styles.card} ${styles.rootCauseCard}`} style={{ justifyContent: "center", alignItems: "center", padding: "var(--space-8)" }}>
            <Loader2 size={24} className={styles.spin} style={{ color: "var(--accent-indigo-light)", marginBottom: "var(--space-4)" }} />
            <span style={{ color: "var(--accent-indigo-light)", fontSize: "14px", fontWeight: 500 }}>
              {investigation.status === "pending" ? "Investigation pending..." : "Investigation running..."}
            </span>
          </div>
        );
      case "failed":
        return (
          <div className={`${styles.card} ${styles.rootCauseCard}`} style={{ borderColor: "var(--status-error-border)", background: "linear-gradient(180deg, rgba(244, 63, 94, 0.1) 0%, var(--surface-1) 100%)" }}>
            <div className={styles.cardHeader} style={{ color: "var(--status-error-light)" }}>
              <AlertTriangle size={16} />
              <span>Investigation failed (Execution error)</span>
            </div>
            <div className={styles.rootCauseText}>
              The investigation encountered an error during execution and could not complete the root cause analysis.
            </div>
          </div>
        );
      case "completed":
        if (investigation.root_cause_status === "unknown" || incident.root_cause_status === "unknown") {
          return (
            <div className={`${styles.card} ${styles.rootCauseCard}`} style={{ borderColor: "var(--border-strong)", background: "var(--surface-2)" }}>
              <div className={styles.cardHeader} style={{ color: "var(--text-secondary)" }}>
                <AlertCircle size={16} />
                <span>Investigation completed (Root cause unknown)</span>
              </div>
              <div className={styles.rootCauseText}>
                The automated investigation finished analyzing the available evidence but could not determine a definitive root cause.
              </div>
            </div>
          );
        }
        
        if (!investigation.root_cause) return null;

        // Custom simple syntax highlighter for the root cause text based on the screenshot format
        const formattedText = investigation.root_cause.split(' ').map((word: string, i: number) => {
          if (word.includes('.go') || word === 'FindByID' || word === 'user.Name') {
            return <span key={i} className={styles.inlineCode}>{word}</span>;
          }
          return word + ' ';
        });

        return (
          <div className={`${styles.card} ${styles.rootCauseCard}`}>
            <div className={styles.cardHeader}>
              <Sparkles size={16} className={styles.iconCritical} style={{ color: "var(--accent-indigo-light)" }} />
              <span style={{ color: "var(--accent-indigo-light)" }}>Root cause identified</span>
              <div className={styles.confidenceBadge}>
                CONFIDENCE: {Math.round((investigation.confidence ?? 0) * 100)}%
              </div>
            </div>
            <div className={styles.rootCauseText}>
              {formattedText}
            </div>
          </div>
        );
      default:
        return null;
    }
  };

  return renderStatus();
}
