import React from "react";
import { getIncidentsService } from "../../services/get-incidents.service";
import { IncidentsListClient } from "./IncidentsListClient";
import { ErrorState } from "@/core/ui/feedback";
import styles from "./IncidentsListView.module.css";

export const IncidentsListView = async () => {
  let incidents;
  try {
    const res = await getIncidentsService();
    incidents = res.data;
  } catch (error) {
    return (
      <div className={styles.container}>
        <div className={styles.header}>
          <div className={styles.titleGroup}>
            <h1 className={styles.title}>Incidents</h1>
            <p className={styles.subtitle}>
              Production incidents requiring autonomous remediation.
            </p>
          </div>
        </div>
        <ErrorState error={error as Error} />
      </div>
    );
  }

  return (
    <IncidentsListClient initialIncidents={incidents} />
  );
};
