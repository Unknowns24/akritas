import React from "react";
import { getIncidentService } from "../../services/get-incident.service";
import { IncidentHeader } from "./components/IncidentHeader";
import { RootCauseCard } from "./components/RootCauseCard";
import { StackTraceCard } from "./components/StackTraceCard";
import { ContextCards } from "./components/ContextCards";
import { RemediationCard } from "./components/RemediationCard";
import styles from "./IncidentDetailView.module.css";

export const IncidentDetailView = async ({ id }: { id: string }) => {
  let incident;
  try {
    incident = await getIncidentService(id);
  } catch (error) {
    // Basic error handling for now
    return <div className={styles.container}><h1>Error loading incident {id}</h1></div>;
  }

  return (
    <div className={styles.container}>
      <IncidentHeader incident={incident} />

      <div className={styles.contentGrid}>
        <div className={styles.leftColumn}>
          <RootCauseCard incident={incident} />
          <StackTraceCard incident={incident} />
          <ContextCards incident={incident} />
        </div>
        
        <div className={styles.rightColumn}>
          <RemediationCard incident={incident} />
        </div>
      </div>
    </div>
  );
};
