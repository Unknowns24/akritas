import React from "react";
import { UserCheck, ExternalLink, ShieldAlert, ListChecks } from "lucide-react";
import type { Incident } from "../../../services/get-incident.service";
import styles from "./RequiresHumanCard.module.css";

interface RequiresHumanCardProps {
  incident: Incident;
}

export function RequiresHumanCard({ incident }: RequiresHumanCardProps) {
  const investigation = incident.latest_investigation;
  const issue = incident.github_issue_reference;
  const recommendedActions = investigation?.recommended_actions || [];

  return (
    <div className={styles.card}>
      <div className={styles.content}>
        <div className={styles.header}>
          <div className={styles.titleGroup}>
            <UserCheck size={18} style={{ color: "#eab308" }} />
            <span>Remediation Boundary</span>
          </div>
          <span className={styles.badge}>Requires Human</span>
        </div>

        <div className={styles.explanation}>
          Akritas investigated this incident and determined that automatic remediation
          is not safe or feasible within the repository boundary. Manual intervention
          by an engineer is required to resolve this issue.
        </div>

        {issue && (
          <div className={styles.section}>
            <span className={styles.sectionLabel}>Documented GitHub Issue</span>
            <div className={styles.issueBox}>
              <div className={styles.issueMeta}>
                <span className={styles.issueNumber}>Issue #{issue.number}</span>
                <span className={styles.issueRepo}>{issue.repository}</span>
              </div>
              <a
                href={issue.url}
                target="_blank"
                rel="noreferrer"
                className={styles.issueLink}
              >
                View on GitHub
                <ExternalLink size={12} />
              </a>
            </div>
          </div>
        )}

        {recommendedActions.length > 0 && (
          <div className={styles.section}>
            <span className={styles.sectionLabel}>
              <ListChecks size={13} style={{ verticalAlign: "middle", marginRight: "4px" }} />
              Recommended Engineer Actions
            </span>
            <ul className={styles.actionList}>
              {recommendedActions.map((action, index) => (
                <li key={index}>{action}</li>
              ))}
            </ul>
          </div>
        )}
      </div>

      <div className={styles.footer}>
        <ShieldAlert size={14} />
        <span>Autonomous remediation stopped at GitHub Issue (ADR-004 boundary).</span>
      </div>
    </div>
  );
}
