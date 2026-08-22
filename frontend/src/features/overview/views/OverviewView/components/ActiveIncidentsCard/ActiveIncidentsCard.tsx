import React from "react";
import { AlertTriangle, ShieldCheck } from "lucide-react";
import { Badge } from "@/core/ui/primitives/Badge";
import { Card, CardBody, CardHeader } from "@/core/ui/primitives/Card";
import styles from "./ActiveIncidentsCard.module.css";

export const ActiveIncidentsCard: React.FC = () => {
  return (
    <Card>
      <CardHeader>
        <div className={styles.sectionTitle}>
          <AlertTriangle size={16} />
          <span>Active Incidents & Diagnostics</span>
        </div>
        <Badge variant="neutral">0 ACTIVE</Badge>
      </CardHeader>
      <CardBody>
        <div className={styles.emptyState}>
          <ShieldCheck size={36} className={styles.emptyIcon} />
          <h2 className={styles.emptyTitle}>All Systems Nominal</h2>
          <p className={styles.emptyText}>
            No active incidents detected. When an application log anomaly triggers a rule,
            investigations will appear here in real time.
          </p>
        </div>
      </CardBody>
    </Card>
  );
};
