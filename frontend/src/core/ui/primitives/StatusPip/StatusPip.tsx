import React from "react";
import styles from "./StatusPip.module.css";

export type SystemHealthStatus = "healthy" | "degraded" | "incident" | "offline";

export interface StatusPipProps {
  status: SystemHealthStatus;
  pulse?: boolean;
  label?: string;
  className?: string;
}

export const StatusPip: React.FC<StatusPipProps> = ({
  status,
  pulse = false,
  label,
  className = "",
}) => {
  const statusClass = styles[`status-${status}`] || styles["status-offline"];
  const pulseClass = pulse ? styles.pulse : "";

  return (
    <span className={`${styles.container} ${className}`.trim()}>
      <span className={`${styles.pip} ${statusClass} ${pulseClass}`.trim()} aria-hidden="true" />
      {label && <span className={styles.label}>{label}</span>}
    </span>
  );
};
