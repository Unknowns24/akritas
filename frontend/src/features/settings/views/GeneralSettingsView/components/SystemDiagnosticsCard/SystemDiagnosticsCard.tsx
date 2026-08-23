import { Play, Code2, Server, Cpu, Activity, LucideIcon } from "lucide-react";
import { Card } from "@/core/ui/primitives/Card";
import { Button } from "@/core/ui/primitives/Button";
import { StatusPip } from "@/core/ui/primitives/StatusPip";
import styles from "./SystemDiagnosticsCard.module.css";
import type { components } from "@/core/libs/api-client";

type SystemStatus = components["schemas"]["SystemStatus"];
type ComponentHealth = components["schemas"]["ComponentHealth"];

interface SystemDiagnosticsCardProps {
  status: SystemStatus;
}

const COMPONENT_INFO: Record<ComponentHealth["component"], { name: string, icon: LucideIcon }> = {
  github: { name: "GITHUB API", icon: Code2 },
  qvac: { name: "QVAC", icon: Cpu },
  dokploy: { name: "DOKPLOY", icon: Server },
  investigator: { name: "INVESTIGATOR", icon: Activity },
};

export function SystemDiagnosticsCard({ status }: SystemDiagnosticsCardProps) {
  // Format the last diagnostics time relative to now (e.g., "46S AGO")
  // For simplicity, we just format it as a localized string or static placeholder if it's complex,
  // but a real implementation would use a relative time formatter.
  const timeStr = status.last_diagnostics_at ? new Date(status.last_diagnostics_at).toLocaleTimeString() : 'UNKNOWN';

  return (
    <Card className={styles.card}>
      <div className={styles.header}>
        <h2 className={styles.title}>System Diagnostics</h2>
        <div className={styles.lastCheck}>LAST SYSTEM CHECK {timeStr}</div>
      </div>

      <div className={styles.content}>
        <div className={styles.grid}>
          {status.components.map((comp) => {
            const info = COMPONENT_INFO[comp.component];
            const Icon = info.icon;
            
            return (
              <div key={comp.component} className={`${styles.item} ${comp.status === "healthy" ? '' : styles.itemError}`}>
                <div className={styles.itemHeader}>
                  <Icon size={14} className={styles.itemIcon} />
                  <span className={styles.itemName}>{info.name}</span>
                </div>
                <div className={styles.itemStatus}>
                  <StatusPip status={comp.status === "healthy" ? "healthy" : "offline"} />
                  <span className={styles.statusText}>{comp.status}</span>
                </div>
              </div>
            );
          })}
        </div>

        <div className={styles.footer}>
          <Button variant="secondary" className={styles.button}>
            <Play size={16} className={styles.buttonIcon} />
            Run Diagnostics
          </Button>
        </div>
      </div>
    </Card>
  );
}
