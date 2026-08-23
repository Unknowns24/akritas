import React from "react";
import { Wrench, CheckCircle2, ShieldCheck } from "lucide-react";
import type { Incident } from "../../../services/get-incident.service";
import styles from "../IncidentDetailView.module.css";

export function RemediationCard({ incident }: { incident: Incident }) {
  const remediation = incident.remediation;
  if (!remediation) return null;

  const diffLines = remediation.changes?.[0]?.patch?.split("\n") || [];

  return (
    <div className={styles.card} style={{ height: "100%", justifyContent: "space-between" }}>
      <div>
        <div className={styles.cardHeader}>
          <Wrench size={16} />
          Repository Analysis & Remediation
        </div>
        
        <div className={styles.remediationMeta} style={{ marginTop: "16px" }}>
          Files inspected: {incident.latest_investigation?.relevant_files?.join(", ")}
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
            <span style={{ fontFamily: "var(--font-mono)", fontSize: "12px", color: "var(--text-secondary)" }}>
              <GitMergeIcon /> akritas/fix-incident-184
            </span>
          </div>
        </div>
        
        <div style={{ display: "flex", gap: "12px", marginTop: "12px" }}>
          {remediation.validation_summary && remediation.validation_summary.passed > 0 && remediation.validation_summary.failed === 0 && (
            <div className={styles.validationTag}>
              <CheckCircle2 size={12} />
              Validation Passed
            </div>
          )}
        </div>

        <button className={styles.prButton}>
          Open Pull Request
        </button>
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
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ verticalAlign: "middle", marginRight: "4px" }}>
      <circle cx="18" cy="18" r="3"></circle>
      <circle cx="6" cy="6" r="3"></circle>
      <path d="M6 21V9a9 9 0 0 0 9 9"></path>
    </svg>
  );
}
