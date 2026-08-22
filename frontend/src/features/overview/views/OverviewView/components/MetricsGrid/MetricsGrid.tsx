import React from "react";
import { Activity, AlertTriangle, FolderGit2, GitPullRequest } from "lucide-react";
import styles from "./MetricsGrid.module.css";

export const MetricsGrid: React.FC = () => {
  return (
    <div className={styles.metricsGrid}>
      <div className={styles.metricCard}>
        <div className={styles.metricHeader}>
          <span className={styles.metricLabel}>Monitored Projects</span>
          <FolderGit2 size={16} />
        </div>
        <span className={styles.metricValue}>0</span>
        <span className={styles.metricFootnote}>Active Dokploy connections</span>
      </div>

      <div className={styles.metricCard}>
        <div className={styles.metricHeader}>
          <span className={styles.metricLabel}>Active Incidents</span>
          <AlertTriangle size={16} color="var(--status-warning)" />
        </div>
        <span className={styles.metricValue}>0</span>
        <span className={styles.metricFootnote}>Under investigation</span>
      </div>

      <div className={styles.metricCard}>
        <div className={styles.metricHeader}>
          <span className={styles.metricLabel}>Completed Workflows</span>
          <Activity size={16} color="var(--status-success)" />
        </div>
        <span className={styles.metricValue}>0</span>
        <span className={styles.metricFootnote}>Closed incident pipelines</span>
      </div>

      <div className={styles.metricCard}>
        <div className={styles.metricHeader}>
          <span className={styles.metricLabel}>PRs Created</span>
          <GitPullRequest size={16} color="var(--accent-indigo-light)" />
        </div>
        <span className={styles.metricValue}>0</span>
        <span className={styles.metricFootnote}>Validated autonomous fixes</span>
      </div>
    </div>
  );
};
