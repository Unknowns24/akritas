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
        const isUnknown = investigation.root_cause_status === "unknown" || incident.root_cause_status === "unknown";
        
        if (isUnknown) {
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
        
        // Custom simple syntax highlighter for the root cause text
        const formatText = (text: string) => {
          return text.split(' ').map((word: string, i: number) => {
            if (word.includes('.go') || word === 'FindByID' || word === 'user.Name') {
              return <span key={i} className={styles.inlineCode}>{word}</span>;
            }
            return word + ' ';
          });
        };

        return (
          <div className={`${styles.card} ${styles.rootCauseCard}`}>
            <div className={styles.cardHeader}>
              <Sparkles size={16} className={styles.iconCritical} style={{ color: "var(--accent-indigo-light)" }} />
              <span style={{ color: "var(--accent-indigo-light)" }}>
                Root cause {investigation.root_cause_status}
              </span>
              <div style={{ display: "flex", gap: "var(--space-2)", marginLeft: "auto" }}>
                {investigation.resolution_status && (
                  <div className={styles.confidenceBadge} style={{ background: "var(--surface-2)", color: "var(--text-secondary)", borderColor: "var(--border-strong)" }}>
                    RESOLUTION: {investigation.resolution_status.toUpperCase().replace(/_/g, " ")}
                  </div>
                )}
                {investigation.confidence !== undefined && (
                  <div className={styles.confidenceBadge}>
                    CONFIDENCE: {Math.round(investigation.confidence * 100)}%
                  </div>
                )}
              </div>
            </div>

            {investigation.summary && (
              <div>
                <div className={styles.contextLabel}>Summary</div>
                <div className={styles.rootCauseText}>
                  {investigation.summary}
                </div>
              </div>
            )}

            {investigation.root_cause && (
              <div>
                <div className={styles.contextLabel}>Root Cause Analysis</div>
                <div className={styles.rootCauseText}>
                  {formatText(investigation.root_cause)}
                </div>
              </div>
            )}

            {investigation.relevant_files && investigation.relevant_files.length > 0 && (
              <div>
                <div className={styles.contextLabel}>Relevant Files</div>
                <div style={{ display: "flex", gap: "var(--space-2)", flexWrap: "wrap", marginTop: "var(--space-2)" }}>
                  {investigation.relevant_files.map((file, idx) => (
                    <span key={idx} className={styles.inlineCode}>{file}</span>
                  ))}
                </div>
              </div>
            )}

            {investigation.recommended_actions && investigation.recommended_actions.length > 0 && (
              <div>
                <div className={styles.contextLabel}>Recommended Actions</div>
                <ul style={{ margin: "var(--space-2) 0 0", paddingLeft: "var(--space-4)", color: "var(--text-secondary)", fontSize: "14px", lineHeight: 1.6 }}>
                  {investigation.recommended_actions.map((action, idx) => (
                    <li key={idx} style={{ marginBottom: "var(--space-1)" }}>{action}</li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        );
      default:
        return null;
    }
  };

  return renderStatus();
}
