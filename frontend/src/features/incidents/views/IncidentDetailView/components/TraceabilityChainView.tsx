import React from "react";
import { Layers, ShieldCheck } from "lucide-react";
import type { Incident } from "../../../services/get-incident.service";
import { buildIncidentTraceabilityChain } from "../../../utils/traceability.utils";
import { TraceabilityStepNode } from "./TraceabilityStepNode";
import styles from "./TraceabilityChainView.module.css";

interface TraceabilityChainViewProps {
  incident: Incident;
}

export function TraceabilityChainView({ incident }: TraceabilityChainViewProps) {
  const steps = buildIncidentTraceabilityChain(incident);
  const isPrCreated = steps.some((s) => s.id === "pull_request" && s.status === "completed");
  const isHuman = incident.resolution_status === "requires_human";
  const isFailed = incident.remediation?.status === "failed";

  const badgeText = isPrCreated
    ? "PR Linked"
    : isHuman
    ? "Manual Boundary"
    : isFailed
    ? "Validation Failed"
    : "Active Pipeline";

  return (
    <div className={styles.card}>
      <div className={styles.header}>
        <div className={styles.titleGroup}>
          <Layers size={18} className={styles.titleIcon} />
          <span>Incident Lineage & Traceability</span>
        </div>
        <span className={styles.badge}>{badgeText}</span>
      </div>

      <div className={styles.description}>
        End-to-end traceability connecting the telemetry incident to investigation, tracking issue, remediation, branch, commit, and Pull Request.
      </div>

      <div className={styles.chainContainer}>
        {steps.map((step, index) => {
          const isLast = index === steps.length - 1;
          const isNextCompleted =
            !isLast && (steps[index + 1].status === "completed" || steps[index + 1].status === "running");

          return (
            <div key={step.id} className={styles.stepWrapper}>
              <TraceabilityStepNode step={step} />
              {!isLast && (
                <div
                  className={`${styles.connector} ${
                    isNextCompleted ? styles.connectorActive : ""
                  }`}
                />
              )}
            </div>
          );
        })}
      </div>

      <div className={styles.footer}>
        <ShieldCheck size={14} />
        <span>
          Traceability log is immutable and generated deterministically from verified incident events.
        </span>
      </div>
    </div>
  );
}
