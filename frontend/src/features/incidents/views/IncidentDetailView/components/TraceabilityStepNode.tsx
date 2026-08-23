import React from "react";
import {
  AlertCircle,
  Search,
  Bookmark,
  Wrench,
  GitBranch,
  GitCommit,
  GitPullRequest,
  ExternalLink,
} from "lucide-react";
import type { TraceabilityStep } from "../../../types/traceability.types";
import styles from "./TraceabilityStepNode.module.css";

interface TraceabilityStepNodeProps {
  step: TraceabilityStep;
}

export function TraceabilityStepNode({ step }: TraceabilityStepNodeProps) {
  const renderIcon = () => {
    switch (step.id) {
      case "incident":
        return <AlertCircle size={16} />;
      case "investigation":
        return <Search size={16} />;
      case "issue":
        return <Bookmark size={16} />;
      case "remediation":
        return <Wrench size={16} />;
      case "branch":
        return <GitBranch size={16} />;
      case "commit":
        return <GitCommit size={16} />;
      case "pull_request":
        return <GitPullRequest size={16} />;
    }
  };

  const getStatusNodeClass = () => {
    switch (step.status) {
      case "completed":
        return styles.nodeCompleted;
      case "failed":
        return styles.nodeFailed;
      case "halted":
        return styles.nodeHalted;
      case "running":
        return styles.nodeRunning;
      case "pending":
      case "not_applicable":
      default:
        return styles.nodePending;
    }
  };

  const getIconWrapperClass = () => {
    switch (step.status) {
      case "completed":
        return styles.iconCompleted;
      case "failed":
        return styles.iconFailed;
      case "halted":
        return styles.iconHalted;
      case "running":
        return styles.iconRunning;
      default:
        return "";
    }
  };

  const getBadgeClass = () => {
    switch (step.badgeVariant) {
      case "success":
        return styles.badgeSuccess;
      case "error":
        return styles.badgeError;
      case "warning":
        return styles.badgeWarning;
      case "running":
        return styles.badgeRunning;
      case "neutral":
      default:
        return styles.badgeNeutral;
    }
  };

  return (
    <div className={`${styles.node} ${getStatusNodeClass()}`}>
      <div className={`${styles.iconWrapper} ${getIconWrapperClass()}`}>
        {renderIcon()}
      </div>

      <div className={styles.content}>
        <div className={styles.topRow}>
          <span className={styles.title}>{step.title}</span>
          {step.badgeText && (
            <span className={`${styles.badge} ${getBadgeClass()}`}>
              {step.badgeText}
            </span>
          )}
        </div>

        <div className={styles.label}>{step.label}</div>

        {step.detail && <div className={styles.detail}>{step.detail}</div>}

        {step.url && (
          <a
            href={step.url}
            target="_blank"
            rel="noreferrer"
            className={styles.linkBtn}
          >
            <span>{step.externalLabel || "Open Link"}</span>
            <ExternalLink size={11} />
          </a>
        )}
      </div>
    </div>
  );
}

