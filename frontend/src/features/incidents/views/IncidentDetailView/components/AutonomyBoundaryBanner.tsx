import React from "react";
import { ShieldCheck, Lock, UserCheck, ShieldAlert } from "lucide-react";
import styles from "./AutonomyBoundaryBanner.module.css";

export function AutonomyBoundaryBanner() {
  return (
    <div className={styles.banner}>
      <div className={styles.header}>
        <div className={styles.titleGroup}>
          <ShieldCheck size={16} className={styles.shieldIcon} />
          <span>Autonomous Remediation Completed</span>
        </div>
        <span className={styles.badge}>Workflow Finished</span>
      </div>

      <div className={styles.message}>
        Akritas has successfully prepared, validated, and published the code changes as a GitHub Pull Request.
        In accordance with safety principles (ADR-004), autonomous execution stops here.
        Code review, merge, and deployment must be performed by human engineers through standard GitHub processes.
      </div>

      <div className={styles.policyGrid}>
        <div className={`${styles.policyTag} ${styles.policyHighlight}`}>
          <UserCheck size={12} />
          <span>Human Review Required</span>
        </div>
        <div className={styles.policyTag}>
          <Lock size={12} />
          <span>No Auto-Merge</span>
        </div>
        <div className={styles.policyTag}>
          <ShieldAlert size={12} />
          <span>No Auto-Deploy</span>
        </div>
      </div>
    </div>
  );
}
