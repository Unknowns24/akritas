import React from "react";
import { AlertTriangle, Clock, Activity, ExternalLink, GitPullRequest, Hash, Zap } from "lucide-react";
import { Button } from "@/core/ui/primitives/Button";
import { Badge } from "@/core/ui/primitives/Badge";
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
          <span style={{ fontWeight: 600 }}>{incident.key}</span>
          <span>•</span>
          <span>{incident.project.name.toUpperCase()}</span>
          <span>•</span>
          <Badge 
            variant={incident.phase === "completed" ? "success" : incident.phase === "failed" ? "error" : "warning"}
          >
            {incident.phase.toUpperCase()}
          </Badge>
          {incident.root_cause_status && (
            <>
              <span>•</span>
              <span style={{ fontSize: "11px", textTransform: "uppercase", letterSpacing: "0.05em", color: "var(--text-muted)" }}>
                {incident.root_cause_status.replace("_", " ")}
              </span>
            </>
          )}
        </div>
        
        <h1 className={styles.pageTitle}>{incident.title}</h1>
        
        <div className={styles.metaInfo}>
          <div className={styles.metaItem} title="Fingerprint">
            <Hash size={14} />
            <span style={{ fontFamily: "var(--font-mono)", fontSize: "12px", background: "var(--surface-2)", padding: "2px 6px", borderRadius: "4px" }}>
              {incident.fingerprint}
            </span>
          </div>
          <div className={styles.metaItem}>
            <Activity size={14} />
            <span>{incident.occurrence_count} occurrences</span>
          </div>
          <div className={styles.metaItem}>
            <Clock size={14} />
            <span>First seen: {new Date(incident.first_seen_at).toLocaleString()}</span>
          </div>
          {incident.last_seen_at && (
            <div className={styles.metaItem}>
              <Zap size={14} />
              <span>Last seen: {new Date(incident.last_seen_at).toLocaleString()}</span>
            </div>
          )}
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
