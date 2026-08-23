import React from "react";
import Link from "next/link";
import { AlertTriangle, ShieldCheck } from "lucide-react";
import { Badge } from "@/core/ui/primitives/Badge";
import { Card, CardBody, CardHeader } from "@/core/ui/primitives/Card";
import styles from "./ActiveIncidentsCard.module.css";
import type { OverviewResponse } from "@/features/overview/services/get-overview.service";

interface ActiveIncidentsCardProps {
  investigations?: OverviewResponse["data"]["active_investigations"];
}

export const ActiveIncidentsCard: React.FC<ActiveIncidentsCardProps> = ({ investigations = [] }) => {
  return (
    <Card>
      <CardHeader>
        <div className={styles.sectionTitle}>
          <AlertTriangle size={16} />
          <span>Active Incidents & Diagnostics</span>
        </div>
        <Badge variant={investigations.length > 0 ? "error" : "neutral"}>
          {investigations.length} ACTIVE
        </Badge>
      </CardHeader>
      <CardBody style={{ padding: 0 }}>
        {investigations.length === 0 ? (
          <div className={styles.emptyState}>
            <ShieldCheck size={36} className={styles.emptyIcon} />
            <h2 className={styles.emptyTitle}>All Systems Nominal</h2>
            <p className={styles.emptyText}>
              No active incidents detected. When an application log anomaly triggers a rule,
              investigations will appear here in real time.
            </p>
          </div>
        ) : (
          <div className={styles.list}>
            {investigations.map((inv) => (
              <Link href={`/incidents/${inv.id}`} key={inv.id} className={styles.item}>
                <div className={styles.itemHeader}>
                  <span className={styles.itemTitle}>{inv.title}</span>
                  <Badge variant={inv.severity === "critical" || inv.severity === "error" ? "error" : inv.severity === "warning" ? "warning" : "neutral"}>
                    {inv.severity}
                  </Badge>
                </div>
                <div className={styles.itemMeta}>
                  <span>{inv.key}</span>
                  <span>•</span>
                  <span className={styles.itemProject}>{inv.project?.name}</span>
                  <span>•</span>
                  <span>{inv.phase}</span>
                </div>
              </Link>
            ))}
          </div>
        )}
      </CardBody>
    </Card>
  );
};
