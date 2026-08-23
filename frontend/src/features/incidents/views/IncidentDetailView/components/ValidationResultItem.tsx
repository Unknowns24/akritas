"use client";

import React, { useState } from "react";
import {
  TestTube2,
  Hammer,
  Search,
  CheckCircle2,
  XCircle,
  Loader2,
  Clock,
  ChevronDown,
  ChevronRight,
  Terminal,
} from "lucide-react";
import type { ValidationResult } from "../../../types/remediation.types";
import styles from "./ValidationResultItem.module.css";

interface ValidationResultItemProps {
  result: ValidationResult;
  defaultExpanded?: boolean;
}

export function ValidationResultItem({
  result,
  defaultExpanded = false,
}: ValidationResultItemProps) {
  const [expanded, setExpanded] = useState<boolean>(
    defaultExpanded || result.status === "failed"
  );

  const renderTypeIcon = () => {
    switch (result.type) {
      case "test":
        return <TestTube2 size={16} className={styles.typeIcon} />;
      case "build":
        return <Hammer size={16} className={styles.typeIcon} />;
      case "static_analysis":
        return <Search size={16} className={styles.typeIcon} />;
      default:
        return <Terminal size={16} className={styles.typeIcon} />;
    }
  };

  const renderStatusBadge = () => {
    switch (result.status) {
      case "passed":
        return (
          <span className={`${styles.badge} ${styles.passed}`}>
            <CheckCircle2 size={12} />
            Passed
          </span>
        );
      case "failed":
        return (
          <span className={`${styles.badge} ${styles.failed}`}>
            <XCircle size={12} />
            Failed
          </span>
        );
      case "running":
        return (
          <span className={`${styles.badge} ${styles.running}`}>
            <Loader2 size={12} className={styles.spin} />
            Running
          </span>
        );
      case "pending":
      default:
        return (
          <span className={`${styles.badge} ${styles.pending}`}>
            <Clock size={12} />
            Pending
          </span>
        );
    }
  };

  return (
    <div className={styles.item}>
      <div
        className={styles.itemHeader}
        onClick={() => setExpanded((prev) => !prev)}
        role="button"
        tabIndex={0}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            setExpanded((prev) => !prev);
          }
        }}
      >
        <div className={styles.leftGroup}>
          {renderTypeIcon()}
          <div className={styles.nameGroup}>
            <span className={styles.name}>{result.name}</span>
            {result.summary && (
              <span className={styles.summary}>{result.summary}</span>
            )}
          </div>
        </div>

        <div className={styles.rightGroup}>
          {renderStatusBadge()}
          {expanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
        </div>
      </div>

      {expanded && result.output_excerpt && (
        <div className={styles.outputArea}>
          <div className={styles.outputHeader}>
            <span>Execution Trace & Output</span>
            <span className={styles.redactedNote}>Output sanitized (redacted: true)</span>
          </div>
          <pre className={styles.outputPre}>{result.output_excerpt}</pre>
        </div>
      )}
    </div>
  );
}

