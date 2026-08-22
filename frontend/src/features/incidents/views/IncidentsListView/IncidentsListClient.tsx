"use client";

import { useState } from "react";
import { AlertCircle, FileText, CheckCircle2 } from "lucide-react";
import { Button } from "@/core/ui/primitives/Button";
import { Badge } from "@/core/ui/primitives/Badge";
import styles from "./IncidentsListView.module.css";
import type { components } from "@/core/libs/api-client";

type IncidentSummary = components["schemas"]["IncidentSummary"];

interface IncidentsListClientProps {
  initialIncidents: IncidentSummary[];
}

export function IncidentsListClient({ initialIncidents }: IncidentsListClientProps) {
  const [filter, setFilter] = useState<"all" | "active" | "completed">("all");

  const filteredIncidents = initialIncidents.filter((inc) => {
    if (filter === "active") return inc.phase !== "completed";
    if (filter === "completed") return inc.phase === "completed";
    return true;
  });

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.titleGroup}>
          <h1 className={styles.title}>Incidents</h1>
          <p className={styles.subtitle}>Production incidents requiring autonomous remediation.</p>
        </div>
        
        <div className={styles.filters}>
          <Button
            variant={filter === "all" ? "secondary" : "ghost"}
            onClick={() => setFilter("all")}
          >
            All
          </Button>
          <Button
            variant={filter === "active" ? "secondary" : "ghost"}
            onClick={() => setFilter("active")}
          >
            Active
          </Button>
          <Button
            variant={filter === "completed" ? "secondary" : "ghost"}
            onClick={() => setFilter("completed")}
          >
            Completed
          </Button>
        </div>
      </div>

      <div className={styles.list}>
        {filteredIncidents.length === 0 ? (
          <div className={styles.emptyState}>
            <p>No incidents found matching the current filter.</p>
          </div>
        ) : (
          filteredIncidents.map((inc) => (
            <div key={inc.id} className={styles.incidentRow}>
              <div className={styles.incidentMain}>
                <div className={styles.incidentIcon}>
                  {inc.severity === "critical" ? (
                    <AlertCircle className={styles.iconCritical} size={20} />
                  ) : inc.severity === "warning" ? (
                    <AlertCircle className={styles.iconWarning} size={20} />
                  ) : (
                    <FileText className={styles.iconInfo} size={20} />
                  )}
                </div>
                <div className={styles.incidentDetails}>
                  <div className={styles.incidentTitle}>{inc.title}</div>
                  <div className={styles.incidentMeta}>
                    <span className={styles.projectId}>{inc.project.name}</span>
                    <span className={styles.dot}>•</span>
                    <span className={styles.time}>{new Date(inc.first_seen_at).toLocaleString()}</span>
                  </div>
                </div>
              </div>
              
              <div className={styles.incidentStatus}>
                <Badge variant={
                  inc.phase === "completed" ? "success" : 
                  inc.phase === "failed" ? "error" : "warning"
                }>
                  {inc.phase.toUpperCase()}
                </Badge>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
