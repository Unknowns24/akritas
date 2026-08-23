import React from "react";
import { Sparkles } from "lucide-react";
import type { Incident } from "../../../services/get-incident.service";
import styles from "../IncidentDetailView.module.css";

export function RootCauseCard({ incident }: { incident: Incident }) {
  if (!incident.latest_investigation?.root_cause) return null;
  
  // Custom simple syntax highlighter for the root cause text based on the screenshot format
  const formattedText = incident.latest_investigation.root_cause.split(' ').map((word: string, i: number) => {
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
          CONFIDENCE: {Math.round((incident.latest_investigation.confidence ?? 0) * 100)}%
        </div>
      </div>
      <div className={styles.rootCauseText}>
        {formattedText}
      </div>
    </div>
  );
}
