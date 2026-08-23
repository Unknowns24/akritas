import React from "react";
import type { Incident } from "../../../services/get-incident.service";
import styles from "../IncidentDetailView.module.css";
import { Rocket } from "lucide-react";

export function ContextCards({ incident }: { incident: Incident }) {
  return (
    <div className={styles.contextRow}>
      <div className={`${styles.card} ${styles.contextCard}`}>
        <div className={styles.contextLabel}>Occurrences</div>
        <div className={styles.contextValue}>{incident.occurrence_count}</div>
        <div className={styles.contextSubtext}>HTTP 500</div>
      </div>

      {incident.deployment_correlation && (
        <div
          className={`${styles.card} ${styles.contextCard}`}
          style={{ position: "relative", overflow: "hidden" }}
        >
          <div
            style={{
              position: "absolute",
              right: "-10px",
              bottom: "-10px",
              opacity: 0.05,
            }}
          >
            <Rocket size={100} />
          </div>
          <div className={styles.contextLabel}>Deployment Correlation</div>
          <div
            className={styles.contextValue}
            style={{ fontSize: "14px", fontWeight: "600", marginBottom: "4px" }}
          >
            Deployment {incident.deployment_correlation.deployment_identifier}
          </div>
          <div className={styles.contextSubtext}>
            Occurred{" "}
            {incident.deployment_correlation.occurred_before_incident_seconds}{" "}
            seconds before first incident.
          </div>
        </div>
      )}
    </div>
  );
}
