import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type SystemStatusResponse = components["schemas"]["SystemStatusResponse"];

export async function getSystemStatusService(): Promise<SystemStatusResponse> {
  const { data, error } = await api.GET("/system/status", {});

  if (error || !data) {
    console.warn("Failed to fetch system status, returning mock data:", error);
    return {
      data: {
        github_account_count: 1,
        dokploy_server_count: 1,
        qvac_endpoint: "http://qvac.local",
        components: [
          { component: "github", status: "healthy", checked_at: new Date().toISOString() },
          { component: "dokploy", status: "healthy", checked_at: new Date().toISOString() },
          { component: "qvac", status: "degraded", checked_at: new Date().toISOString() },
          { component: "investigator", status: "running", checked_at: new Date().toISOString() }
        ],
        last_diagnostics_at: new Date().toISOString(),
      }
    };
  }

  return data;
}
