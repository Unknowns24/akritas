"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { AlertCircle, FileText } from "lucide-react";
import { Button } from "@/core/ui/primitives/Button";
import { Badge } from "@/core/ui/primitives/Badge";
import { EmptyState } from "@/core/ui/feedback";
import styles from "./IncidentsListView.module.css";
import type { components } from "@/core/libs/api-client";

type IncidentSummary = components["schemas"]["IncidentSummary"];

interface IncidentsListClientProps {
  initialIncidents: IncidentSummary[];
}

export function IncidentsListClient({
  initialIncidents,
}: IncidentsListClientProps) {
  const router = useRouter();
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
          <p className={styles.subtitle}>
            Production incidents requiring autonomous remediation.
          </p>
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
          <EmptyState 
            icon={<AlertCircle size={32} />}
            title="No incidents found"
            description="No incidents found matching the current filter."
          />
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
                  <div className={styles.incidentHeader}>
                    <span className={styles.incidentKey}>{inc.key}</span>
                    <div className={styles.incidentTitle}>{inc.title}</div>
                    <span className={styles.incidentFingerprint} title={inc.fingerprint}>
                      {inc.fingerprint}
                    </span>
                  </div>
                  
                  {inc.summary && (
                    <div className={styles.incidentSummary}>
                      {inc.summary}
                    </div>
                  )}

                  <div className={styles.incidentMeta}>
                    <div className={styles.metaItem}>
                      <span className={styles.projectId}>{inc.project.name}</span>
                    </div>
                    <span className={styles.dot}>•</span>
                    <div className={styles.metaItem}>
                      <span>{inc.occurrence_count} occurrences</span>
                    </div>
                    <span className={styles.dot}>•</span>
                    <div className={styles.metaItem}>
                      <span className={styles.time} title={new Date(inc.first_seen_at).toLocaleString()}>
                        First seen: {new Date(inc.first_seen_at).toLocaleDateString()}
                      </span>
                    </div>
                    {inc.last_seen_at && (
                      <>
                        <span className={styles.dot}>•</span>
                        <div className={styles.metaItem}>
                          <span className={styles.time} title={new Date(inc.last_seen_at).toLocaleString()}>
                            Last seen: {new Date(inc.last_seen_at).toLocaleDateString()}
                          </span>
                        </div>
                      </>
                    )}
                  </div>
                </div>
              </div>

              <div className={styles.incidentStatusGroup}>
                <Badge
                  variant={
                    inc.phase === "completed"
                      ? "success"
                      : inc.phase === "failed"
                        ? "error"
                        : "warning"
                  }
                >
                  {inc.phase.toUpperCase()}
                </Badge>
                {inc.root_cause_status && (
                  <div className={styles.rootCauseStatus}>
                    {inc.root_cause_status.replace("_", " ")}
                  </div>
                )}
              </div>

              <div className={styles.incidentActions}>
                {inc.phase !== "completed" && inc.phase !== "failed" && (
                  <Button
                    variant="primary"
                    size="sm"
                    className={styles.actionBtn}
                    onClick={() => router.push(`/incidents/${inc.id}`)}
                  >
                    Examine
                  </Button>
                )}
                {inc.phase === "completed" && (
                  <Button
                    variant="secondary"
                    size="sm"
                    className={styles.actionBtn}
                    onClick={() => router.push(`/incidents/${inc.id}`)}
                  >
                    View Report
                  </Button>
                )}
                {inc.phase === "failed" && (
                  <Button
                    variant="secondary"
                    size="sm"
                    className={styles.actionBtn}
                    onClick={() => router.push(`/incidents/${inc.id}`)}
                  >
                    View Details
                  </Button>
                )}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
