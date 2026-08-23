import React from "react";
import type { ValidationResult, RemediationStatus } from "../../../types/remediation.types";
import { ValidationResultItem } from "./ValidationResultItem";
import { ValidationFailureBanner } from "./ValidationFailureBanner";
import styles from "./ValidationResultsViewer.module.css";

interface ValidationResultsViewerProps {
  results?: ValidationResult[];
  remediationStatus?: RemediationStatus;
  failureUserMessage?: string;
}

export function ValidationResultsViewer({
  results = [],
  remediationStatus,
  failureUserMessage,
}: ValidationResultsViewerProps) {
  const hasFailed =
    remediationStatus === "failed" || results.some((r) => r.status === "failed");

  return (
    <div className={styles.container}>
      {hasFailed && <ValidationFailureBanner message={failureUserMessage} />}

      <div className={styles.header}>
        <span className={styles.title}>
          Validation Checks {results.length > 0 ? `(${results.length})` : ""}
        </span>
      </div>

      {results.length > 0 ? (
        <div className={styles.list}>
          {results.map((result) => (
            <ValidationResultItem key={result.id} result={result} />
          ))}
        </div>
      ) : (
        <div className={styles.emptyBox}>
          {remediationStatus === "planned"
            ? "Validation checks will be executed once code fixes are generated."
            : remediationStatus === "in_progress"
            ? "Executing test suite and build validation..."
            : "No detailed validation records available."}
        </div>
      )}
    </div>
  );
}
