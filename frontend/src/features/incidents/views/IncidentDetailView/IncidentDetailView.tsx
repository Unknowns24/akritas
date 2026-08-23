import React from "react";
import { getIncidentService } from "../../services/get-incident.service";
import { getInvestigationEvidenceService } from "../../services/get-investigation-evidence.service";
import type { Evidence } from "../../services/get-investigation-evidence.service";
import { IncidentHeader } from "./components/IncidentHeader";
import { RootCauseCard } from "./components/RootCauseCard";
import { EvidenceList } from "./components/EvidenceList";
import { StackTraceCard } from "./components/StackTraceCard";
import { ContextCards } from "./components/ContextCards";
import { LogEventsCard } from "./components/LogEventsCard";
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

  let evidence: Evidence[] = [];
  if (incident.latest_investigation?.id) {
    try {
      const response = await getInvestigationEvidenceService(incident.latest_investigation.id);
      evidence = response.data || [];
    } catch (error) {
      console.error("Failed to load investigation evidence", error);
    }
  }

  return (
    <div className={styles.container}>
      <IncidentHeader incident={incident} />

      <div className={styles.contentGrid}>
        <div className={styles.leftColumn}>
          <RootCauseCard incident={incident} />
          {evidence.length > 0 && <EvidenceList evidence={evidence} />}
          <StackTraceCard incident={incident} />
          <LogEventsCard incidentId={id} />
        </div>
        
        <div className={styles.rightColumn}>
          <ContextCards incident={incident} />
          <RemediationCard incident={incident} />
        </div>
      </div>
    </div>
  );
};
