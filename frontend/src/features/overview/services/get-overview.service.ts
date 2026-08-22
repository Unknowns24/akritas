import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type OverviewResponse = components["schemas"]["OverviewResponse"];

export async function getOverviewService(): Promise<OverviewResponse> {
  const { data, error } = await api.GET("/", {});

  if (error || !data) {
    console.warn("Failed to fetch overview, returning mock data:", error);
    return {
      data: {
        monitored_projects: 2,
        active_incidents: 1,
        workflow_completed_incidents: 5,
        pull_requests_created: 3,
        active_investigations: [
          {
            id: "inc-1",
            key: "AKR-1",
            project: { id: "1", name: "E-Commerce Platform" },
            fingerprint: "db_conn_timeout",
            severity: "critical",
            title: "Database connection timeout in production",
            summary: "Multiple instances reporting timeouts when connecting to the primary DB.",
            phase: "investigating",
            occurrence_count: 42,
            first_seen_at: new Date(Date.now() - 3600000).toISOString(),
            last_seen_at: new Date().toISOString()
          }
        ]
      }
    };
  }

  return data;
}
