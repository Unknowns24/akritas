import React from "react";
import { FileText, Code2, GitCommit, GitPullRequest, Search, Activity, ShieldCheck } from "lucide-react";
import type { Evidence } from "../../../services/get-investigation-evidence.service";
import styles from "../IncidentDetailView.module.css";

interface EvidenceListProps {
  evidence: Evidence[];
}

export function EvidenceList({ evidence }: EvidenceListProps) {
  if (!evidence || evidence.length === 0) {
    return null;
  }

  const getEvidenceIcon = (type: string) => {
    switch (type) {
      case "log_excerpt":
        return <FileText size={14} />;
      case "code_location":
      case "stack_trace":
        return <Code2 size={14} />;
      case "commit":
        return <GitCommit size={14} />;
      case "diff":
        return <GitPullRequest size={14} />;
      case "validation_result":
        return <ShieldCheck size={14} />;
      case "deployment_metadata":
        return <Activity size={14} />;
      default:
        return <Search size={14} />;
    }
  };

  const getEvidenceTitle = (type: string) => {
    return type.split("_").map((word) => word.charAt(0).toUpperCase() + word.slice(1)).join(" ");
  };

  const renderDiffContent = (patch: string) => {
    const lines = patch.split('\n');
    return (
      <div className={styles.diffViewer}>
        {lines.map((line, idx) => {
          let lineClass = styles.diffUnchanged;
          if (line.startsWith('+')) lineClass = styles.diffAdded;
          else if (line.startsWith('-')) lineClass = styles.diffUnchanged; // or define a diffRemoved class
          else if (line.startsWith('@@')) lineClass = styles.diffMeta;

          return (
            <div key={idx} className={`${styles.diffLine} ${lineClass}`}>
              {line}
            </div>
          );
        })}
      </div>
    );
  };

  return (
    <div className={styles.card} style={{ marginTop: "var(--space-6)" }}>
      <div className={styles.cardHeader}>
        <Search size={18} style={{ color: "var(--text-secondary)" }} />
        <span>Investigation Evidence</span>
        <div style={{ marginLeft: "auto", fontSize: "12px", color: "var(--status-success-light)", display: "flex", alignItems: "center", gap: "var(--space-1)" }}>
          <ShieldCheck size={14} />
          <span>Sanitized</span>
        </div>
      </div>
      <div style={{ fontSize: "14px", color: "var(--text-secondary)", marginBottom: "var(--space-4)" }}>
        {evidence.length} pieces of evidence were analyzed during this investigation.
      </div>

      <div className={styles.logEventsList}>
        {evidence.map((item) => (
          <div key={item.id} className={styles.logEventContainer}>
            <div className={styles.logEventHeader}>
              <div style={{ display: "flex", alignItems: "center", gap: "var(--space-2)" }}>
                <span style={{ color: "var(--text-secondary)" }}>
                  {getEvidenceIcon(item.type)}
                </span>
                <span className={styles.logEventSeverity}>
                  {getEvidenceTitle(item.type)}
                </span>
                {item.file_path && (
                  <>
                    <span style={{ color: "var(--border-strong)" }}>•</span>
                    <span className={styles.logEventTimestamp}>
                      {item.file_path}
                      {item.line_start ? `:${item.line_start}` : ''}
                      {item.line_end && item.line_end !== item.line_start ? `-${item.line_end}` : ''}
                    </span>
                  </>
                )}
                {item.commit_sha && (
                  <>
                    <span style={{ color: "var(--border-strong)" }}>•</span>
                    <span className={styles.logEventTimestamp}>
                      commit: {item.commit_sha.substring(0, 7)}
                    </span>
                  </>
                )}
              </div>
              <div className={styles.logEventTimestamp}>
                {new Date(item.occurred_at || item.created_at).toLocaleString()}
              </div>
            </div>
            
            {item.summary && (
              <div className={styles.logContextSection}>
                <span style={{ fontSize: "13px", color: "var(--text-primary)" }}>{item.summary}</span>
              </div>
            )}
            
            {item.content && (
              <div className={styles.logEvidence}>
                <pre>{item.content}</pre>
              </div>
            )}

            {item.patch && (
              <div className={styles.logEvidence} style={{ padding: 0 }}>
                {renderDiffContent(item.patch)}
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
