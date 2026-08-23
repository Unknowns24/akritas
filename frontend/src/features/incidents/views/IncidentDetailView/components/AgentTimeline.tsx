import React from "react";
import { Bot, Terminal, PlayCircle, GitCommit, FileText, CheckCircle2, Search, AlertCircle } from "lucide-react";
import type { TimelineEvent } from "../../../services/get-incident-timeline.service";
import styles from "../IncidentDetailView.module.css";

interface AgentTimelineProps {
  timeline: TimelineEvent[];
}

export function AgentTimeline({ timeline }: AgentTimelineProps) {
  if (!timeline || timeline.length === 0) {
    return null;
  }

  // Filter events to only show agent reasoning/tool activity
  const agentEvents = timeline.filter(t => 
    ["investigation_started", "tool_used", "changes_generated", "root_cause_classified", "workflow_failed"].includes(t.type)
  );

  if (agentEvents.length === 0) return null;

  const getEventIcon = (type: string, summary: string) => {
    if (type === "investigation_started") return <PlayCircle size={14} />;
    if (type === "changes_generated") return <CheckCircle2 size={14} />;
    if (type === "root_cause_classified") return <Search size={14} />;
    if (type === "workflow_failed") return <AlertCircle size={14} />;
    
    // For tool_used, infer icon from summary context
    const s = summary.toLowerCase();
    if (s.includes("log")) return <FileText size={14} />;
    if (s.includes("commit") || s.includes("git")) return <GitCommit size={14} />;
    return <Terminal size={14} />;
  };

  return (
    <div className={styles.card} style={{ marginTop: "var(--space-6)" }}>
      <div className={styles.cardHeader}>
        <Bot size={18} style={{ color: "var(--accent-indigo-light)" }} />
        <span>Agent Activity Timeline</span>
      </div>
      <div style={{ fontSize: "13px", color: "var(--text-secondary)", marginBottom: "var(--space-6)" }}>
        Read-only sequence of tools utilized by the automated agent during this investigation.
      </div>

      <div style={{ position: "relative", paddingLeft: "var(--space-4)" }}>
        {/* Vertical line connecting the timeline dots */}
        <div 
          style={{ 
            position: "absolute", 
            left: "21px", 
            top: "14px", 
            bottom: "24px", 
            width: "2px", 
            backgroundColor: "var(--border-subtle)" 
          }} 
        />

        {agentEvents.map((event, idx) => {
          const isLast = idx === agentEvents.length - 1;
          return (
            <div key={event.id} style={{ position: "relative", paddingBottom: isLast ? 0 : "var(--space-6)" }}>
              {/* Timeline dot/icon */}
              <div 
                style={{
                  position: "absolute",
                  left: "-4px",
                  top: "2px",
                  width: "24px",
                  height: "24px",
                  borderRadius: "50%",
                  backgroundColor: "var(--surface-1)",
                  border: "1px solid var(--border-strong)",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  color: event.type === "workflow_failed" ? "var(--status-error)" : "var(--text-secondary)",
                  borderColor: event.type === "workflow_failed" ? "var(--status-error)" : "var(--border-strong)",
                  zIndex: 2
                }}
              >
                {getEventIcon(event.type, event.summary)}
              </div>

              <div style={{ paddingLeft: "var(--space-6)" }}>
                <div style={{ display: "flex", alignItems: "center", gap: "var(--space-3)", marginBottom: "var(--space-1)" }}>
                  <span style={{ fontSize: "14px", fontWeight: 500, color: event.type === "workflow_failed" ? "var(--status-error)" : "var(--text-primary)" }}>
                    {event.summary}
                  </span>
                  <span style={{ fontSize: "12px", fontFamily: "var(--font-mono)", color: "var(--text-muted)" }}>
                    {new Date(event.occurred_at).toLocaleTimeString()}
                  </span>
                </div>
                
                {event.detail && (
                  <div style={{ 
                    marginTop: "var(--space-2)", 
                    padding: "var(--space-2) var(--space-3)", 
                    backgroundColor: "var(--surface-2)", 
                    border: "1px solid var(--border-subtle)", 
                    borderRadius: "var(--radius-sm)",
                    fontFamily: "var(--font-mono)",
                    fontSize: "12px",
                    color: "var(--text-secondary)",
                    overflowX: "auto"
                  }}>
                    {event.detail}
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
