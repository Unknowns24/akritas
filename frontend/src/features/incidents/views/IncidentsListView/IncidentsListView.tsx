"use client";

import React, { useState } from "react";
import { Filter, ShieldCheck } from "lucide-react";
import { Badge } from "@/core/ui/primitives/Badge";
import styles from "./IncidentsListView.module.css";

export const IncidentsListView: React.FC = () => {
  const [activeFilter, setActiveFilter] = useState<"all" | "active" | "completed">("all");

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.titleGroup}>
          <h1 className={styles.pageTitle}>Incidents Stream</h1>
          <p className={styles.pageSubtitle}>
            Real-time incident detection, log evidence clustering, and remediation pipelines.
          </p>
        </div>
      </div>

      {/* Filter and controls bar */}
      <div className={styles.filterBar}>
        <div className={styles.filterTabs}>
          <button
            type="button"
            className={`${styles.tabButton} ${activeFilter === "all" ? styles.tabButtonActive : ""}`.trim()}
            onClick={() => setActiveFilter("all")}
          >
            All Incidents
          </button>
          <button
            type="button"
            className={`${styles.tabButton} ${activeFilter === "active" ? styles.tabButtonActive : ""}`.trim()}
            onClick={() => setActiveFilter("active")}
          >
            Active / Investigating
          </button>
          <button
            type="button"
            className={`${styles.tabButton} ${activeFilter === "completed" ? styles.tabButtonActive : ""}`.trim()}
            onClick={() => setActiveFilter("completed")}
          >
            Workflow Completed
          </button>
        </div>

        <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
          <Filter size={14} color="var(--text-dim)" />
          <Badge variant="neutral">0 Total</Badge>
        </div>
      </div>

      {/* Empty State / Incidents table container */}
      <div className={styles.emptyCard}>
        <ShieldCheck size={40} className={styles.emptyIcon} />
        <h2 className={styles.emptyTitle}>No Incidents Recorded</h2>
        <p className={styles.emptyText}>
          Dokploy log monitoring has not captured any error signatures matching positive detection
          rules. When incidents occur, they will be triaged here automatically.
        </p>
      </div>
    </div>
  );
};
