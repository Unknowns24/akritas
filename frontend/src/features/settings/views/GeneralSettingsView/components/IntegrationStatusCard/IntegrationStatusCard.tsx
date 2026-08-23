import { Code2, Server, Cpu } from "lucide-react";
import { Card } from "@/core/ui/primitives/Card";
import styles from "./IntegrationStatusCard.module.css";
import type { components } from "@/core/libs/api-client";

type SystemStatus = components["schemas"]["SystemStatus"];

interface IntegrationStatusCardProps {
  status: SystemStatus;
}

export function IntegrationStatusCard({ status }: IntegrationStatusCardProps) {
  return (
    <Card className={styles.card}>
      <div className={styles.header}>
        <h2 className={styles.title}>Integration Status</h2>
      </div>
      
      <div className={styles.content}>
        {/* GitHub */}
        <div className={styles.row}>
          <div className={styles.iconWrapper}>
            <Code2 size={20} />
          </div>
          <div className={styles.info}>
            <div className={styles.name}>GitHub</div>
            <div className={styles.details}>
              {status.github_account_count} account{status.github_account_count !== 1 ? 's' : ''} <span className={styles.dot}>•</span> <span className={styles.connected}>Connected</span>
            </div>
          </div>
        </div>

        {/* Dokploy */}
        <div className={styles.row}>
          <div className={styles.iconWrapper}>
            <Server size={20} />
          </div>
          <div className={styles.info}>
            <div className={styles.name}>Dokploy</div>
            <div className={styles.details}>
              {status.dokploy_server_count} server{status.dokploy_server_count !== 1 ? 's' : ''} <span className={styles.dot}>•</span> <span className={styles.connected}>Connected</span>
            </div>
          </div>
        </div>

        {/* QVAC */}
        <div className={styles.row}>
          <div className={styles.iconWrapper}>
            <Cpu size={20} />
          </div>
          <div className={styles.info}>
            <div className={styles.name}>QVAC</div>
            <div className={styles.details}>
              <span className={styles.mono}>{status.qvac_endpoint}</span> <span className={styles.dot}>•</span> <span className={styles.connected}>Connected</span>
            </div>
          </div>
        </div>
      </div>
    </Card>
  );
}
