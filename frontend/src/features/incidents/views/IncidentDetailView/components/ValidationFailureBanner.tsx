import React from "react";
import { AlertTriangle, ShieldX } from "lucide-react";
import styles from "./ValidationFailureBanner.module.css";

interface ValidationFailureBannerProps {
  message?: string;
}

export function ValidationFailureBanner({ message }: ValidationFailureBannerProps) {
  return (
    <div className={styles.banner}>
      <div className={styles.header}>
        <AlertTriangle size={16} />
        <span>Remediation Failed</span>
      </div>

      {message && <div className={styles.message}>{message}</div>}

      <div className={styles.blockageNotice}>
        <ShieldX size={14} style={{ flexShrink: 0 }} />
        <span>
          <strong>No Pull Request Created:</strong> In accordance with safety policies (ADR-004),
          Akritas halts the autonomous pipeline and will not open a PR for unvalidated code changes.
        </span>
      </div>
    </div>
  );
}

