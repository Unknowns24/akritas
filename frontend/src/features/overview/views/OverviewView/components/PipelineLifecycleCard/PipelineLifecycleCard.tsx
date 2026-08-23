import React from "react";
import { Card, CardBody, CardHeader } from "@/core/ui/primitives/Card";
import { StatusPip } from "@/core/ui/primitives/StatusPip";
import styles from "./PipelineLifecycleCard.module.css";

export const PipelineLifecycleCard: React.FC = () => {
  return (
    <Card>
      <CardHeader>
        <span style={{ fontSize: "13px", fontWeight: 600 }}>Autonomous Pipeline Lifecycle</span>
      </CardHeader>
      <CardBody>
        <div className={styles.pipelineList}>
          <div className={styles.pipelineStep}>
            <span>1. Dokploy Log Stream</span>
            <StatusPip status="healthy" />
          </div>
          <div className={styles.pipelineStep}>
            <span>2. Grouping & Fingerprint</span>
            <StatusPip status="healthy" />
          </div>
          <div className={styles.pipelineStep}>
            <span>3. QVAC Root Cause Analysis</span>
            <StatusPip status="healthy" />
          </div>
          <div className={styles.pipelineStep}>
            <span>4. GitHub Issue Created</span>
            <StatusPip status="healthy" />
          </div>
          <div className={styles.pipelineStep}>
            <span>5. Fix Validation & PR</span>
            <StatusPip status="healthy" />
          </div>
        </div>
      </CardBody>
    </Card>
  );
};
