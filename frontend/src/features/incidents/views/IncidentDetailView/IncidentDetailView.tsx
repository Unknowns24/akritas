import React from "react";
import { getIncidentService } from "../../services/get-incident.service";
import { getInvestigationEvidenceService } from "../../services/get-investigation-evidence.service";
import { getIncidentTimelineService } from "../../services/get-incident-timeline.service";
import type { Evidence } from "../../services/get-investigation-evidence.service";
import type { TimelineEvent } from "../../services/get-incident-timeline.service";
import { IncidentHeader } from "./components/IncidentHeader";
import { RootCauseCard } from "./components/RootCauseCard";
import { AgentTimeline } from "./components/AgentTimeline";
import { EvidenceList } from "./components/EvidenceList";
import { StackTraceCard } from "./components/StackTraceCard";
import { ContextCards } from "./components/ContextCards";
import { LogEventsCard } from "./components/LogEventsCard";
import { RemediationCard } from "./components/RemediationCard";
import { TraceabilityChainView } from "./components/TraceabilityChainView";
import { RemediationReviewPacket } from "./components/RemediationReviewPacket";
import { GitHubIssueCard } from "./components/GitHubIssueCard";
import styles from "./IncidentDetailView.module.css";

export const IncidentDetailView = async ({ id }: { id: string }) => {
  let incident;

  try {
    incident = await getIncidentService(id);
  } catch {
    return (
      <div className={styles.container}>
        <h1>Error loading incident {id}</h1>
      </div>
    );
  }

  let evidence: Evidence[] = [];
  let timeline: TimelineEvent[] = [];

  if (incident.latest_investigation?.id) {
    try {
      const [evidenceRes, timelineRes] = await Promise.all([
        getInvestigationEvidenceService(incident.latest_investigation.id),
        getIncidentTimelineService(id),
      ]);

      evidence = evidenceRes.data || [];
      timeline = timelineRes.data || [];
    } catch (error) {
      console.error("Failed to load investigation extra data", error);
    }
  }

  return (
    <div className={styles.container}>
      <IncidentHeader incident={incident} />

      <div className={styles.sectionHeader}>
        <h2>Deterministic Detection</h2>
        <p>Hard evidence collected from monitoring</p>
      </div>

      <div className={styles.sequentialLayout}>
        <StackTraceCard incident={incident} />
        <LogEventsCard incidentId={id} />
        <ContextCards incident={incident} />
        <TraceabilityChainView incident={incident} />
        <RemediationReviewPacket incident={incident} />
      </div>

      <div className={styles.sectionHeader}>
        <h2>QVAC Investigation</h2>
        <p>AI-driven analysis and resolution</p>
      </div>

      <div className={styles.sequentialLayout}>
        <RootCauseCard incident={incident} />

        {timeline.length > 0 && <AgentTimeline timeline={timeline} />}
        {evidence.length > 0 && <EvidenceList evidence={evidence} />}

        <GitHubIssueCard incident={incident} />
        <RemediationCard incident={incident} />
      </div>
    </div>
  );
};