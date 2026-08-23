import React from "react";
import styles from "./EmptyState.module.css";

export interface EmptyStateProps {
  icon?: React.ReactNode;
  title: string;
  description?: React.ReactNode;
  action?: React.ReactNode;
  className?: string;
}

export const EmptyState: React.FC<EmptyStateProps> = ({
  icon,
  title,
  description,
  action,
  className = "",
}) => {
  return (
    <div className={`${styles.emptyState} ${className}`.trim()}>
      {icon && <div className={styles.iconWrapper}>{icon}</div>}
      <h3 className={styles.title}>{title}</h3>
      {description && <div className={styles.description}>{description}</div>}
      {action && <div className={styles.actionWrapper}>{action}</div>}
    </div>
  );
};
