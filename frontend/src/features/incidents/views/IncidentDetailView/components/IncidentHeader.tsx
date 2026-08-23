import React from "react";
import { AlertTriangle, Clock, Activity, ExternalLink, GitPullRequest } from "lucide-react";
import { Button } from "@/core/ui/primitives/Button";
import type { Incident } from "../../../services/get-incident.service";
import styles from "../IncidentDetailView.module.css";

export function IncidentHeader({ incident }: { incident: Incident }) {
  const isCritical = incident.severity === "critical";

  return (
    <div className={styles.header}>
      <div className={styles.titleGroup}>
        <div className={styles.badges}>
          <div className={isCritical ? styles.badgeCritical : ""}>
            <AlertTriangle size={12} />
            {incident.severity.toUpperCase()}
          </div>
          <span>•</span>
          <span>{incident.project.name.toUpperCase()}</span>
          <span>•</span>
          <span>PRODUCTION</span>
        </div>
        <h1 className={styles.pageTitle}>{incident.title}</h1>
        <div className={styles.metaInfo}>
          <div className={styles.metaItem}>
            <Clock size={14} />
            <span>Started {Math.floor((Date.now() - new Date(incident.first_seen_at).getTime()) / 60000)} mins ago</span>
          </div>
          <div className={styles.metaItem}>
            <Activity size={14} />
            <span>{incident.occurrence_count} occurrences</span>
          </div>
        </div>
      </div>
      <div className={styles.headerActions}>
        <Button variant="ghost" className={styles.actionButton}>
          <ExternalLink size={14} />
          View GitHub Issue
        </Button>
        <Button variant="ghost" className={styles.actionButton}>
          <GitPullRequest size={14} />
          View Pull Request
        </Button>
      </div>
    </div>
  );
}
