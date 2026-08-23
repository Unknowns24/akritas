import React from "react";
import { Sparkles } from "lucide-react";
import { Badge } from "@/core/ui/primitives/Badge";
import styles from "./InvestigationEngineBanner.module.css";

export const InvestigationEngineBanner: React.FC = () => {
  return (
    <div className={styles.aiBanner}>
      <div className={styles.aiHeader}>
        <div className={styles.aiTitle}>
          <Sparkles size={16} />
          <span>QVAC Investigation Engine</span>
        </div>
        <Badge variant="intel" dot>
          READY
        </Badge>
      </div>
      <p className={styles.aiText}>
        Deterministic log grouping active. Local AI is primed for root cause analysis, Issue
        generation, and sandbox fix validation.
      </p>
    </div>
  );
};
