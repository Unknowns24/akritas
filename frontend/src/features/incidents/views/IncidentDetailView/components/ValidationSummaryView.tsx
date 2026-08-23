import React from "react";
import { CheckCircle2, XCircle, ShieldCheck } from "lucide-react";
import type { ValidationSummary } from "../../../types/remediation.types";
import styles from "./ValidationSummaryView.module.css";

interface ValidationSummaryViewProps {
  summary?: ValidationSummary;
}

export function ValidationSummaryView({ summary }: ValidationSummaryViewProps) {
  if (!summary || summary.total === 0) {
    return null;
  }

  const isPassed = summary.passed === summary.total && summary.failed === 0;
  const isFailed = summary.failed > 0;

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <span className={styles.label}>Automated Validation</span>
        <span className={styles.countLabel}>
          {summary.passed}/{summary.total} checks passed
        </span>
      </div>

      <div className={styles.metrics}>
        {isPassed && (
          <div className={`${styles.tag} ${styles.tagSuccess}`}>
            <CheckCircle2 size={13} />
            <span>Validation Passed</span>
          </div>
        )}

        {isFailed && (
          <div className={`${styles.tag} ${styles.tagError}`}>
            <XCircle size={13} />
            <span>{summary.failed} check{summary.failed > 1 ? "s" : ""} failed</span>
          </div>
        )}

        {!isPassed && !isFailed && (
          <div className={`${styles.tag} ${styles.tagNeutral}`}>
            <ShieldCheck size={13} />
            <span>Validations Pending</span>
          </div>
        )}
      </div>
    </div>
  );
}
