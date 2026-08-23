import React from "react";
import { Clock, Loader2, CheckCircle2, XCircle, GitPullRequest } from "lucide-react";
import type { RemediationStatus } from "../../../types/remediation.types";
import { getRemediationStatusConfig } from "../../../utils/remediation.utils";
import styles from "./RemediationStatusBadge.module.css";

interface RemediationStatusBadgeProps {
  status?: RemediationStatus;
}

export function RemediationStatusBadge({ status }: RemediationStatusBadgeProps) {
  const config = getRemediationStatusConfig(status);

  const renderIcon = () => {
    switch (status) {
      case "planned":
        return <Clock size={13} />;
      case "in_progress":
        return <Loader2 size={13} className={styles.spin} />;
      case "validated":
        return <CheckCircle2 size={13} />;
      case "failed":
        return <XCircle size={13} />;
      case "pull_request_created":
        return <GitPullRequest size={13} />;
      default:
        return <Clock size={13} />;
    }
  };

  const variantClass = styles[config.variant] || styles.neutral;

  return (
    <span className={`${styles.badge} ${variantClass}`}>
      {renderIcon()}
      {config.label}
    </span>
  );
}
