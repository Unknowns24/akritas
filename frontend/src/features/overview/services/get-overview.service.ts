import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type OverviewResponse = components["schemas"]["OverviewResponse"];

export async function getOverviewService(): Promise<OverviewResponse> {
  try {
    const { data, error } = await api.GET("/overview");
    if (!error && data) return data;
  } catch (e) { console.error("OVERVIEW ERROR:", e); }

  return {
    data: {
      monitored_projects: 0,
      active_incidents: 0,
      workflow_completed_incidents: 0,
      pull_requests_created: 0,
      active_investigations: [],
    },
  };
}
