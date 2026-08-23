import React from "react";
import { TerminalSquare } from "lucide-react";
import type { Incident } from "../../../services/get-incident.service";
import styles from "../IncidentDetailView.module.css";

export function StackTraceCard({ incident }: { incident: Incident }) {
  const trace = (incident.latest_investigation as any)?.stack_traces?.[0];
  if (!trace) return null;

  const lines = trace.raw_content.split("\n");

  return (
    <div className={styles.card}>
      <div className={styles.cardHeader}>
        <TerminalSquare size={16} />
        Stack Trace
      </div>
      <div className={styles.stackTraceContent}>
        {lines.map((line: string, i: number) => {
          if (line.trim().startsWith("//")) {
            return <div key={i} className={styles.traceComment}>{line}</div>;
          }
          if (line.includes("panic") || line.includes("error")) {
            return <div key={i} className={styles.traceError}>{line}</div>;
          }
          return <div key={i}>{line}</div>;
        })}
      </div>
    </div>
  );
}
