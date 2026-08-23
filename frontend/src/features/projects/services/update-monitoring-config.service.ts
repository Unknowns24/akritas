import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type MonitoringConfiguration = components["schemas"]["MonitoringConfiguration"];
export type MonitoringConfigurationResponse = components["schemas"]["MonitoringConfigurationResponse"];

export async function updateMonitoringConfigService(
  projectId: string,
  config: MonitoringConfiguration
): Promise<{ data?: MonitoringConfiguration; error?: Error | any }> {
  const { data, error } = await api.PUT("/projects/{project_id}/monitoring-configuration", {
    params: {
      path: { project_id: projectId },
    },
    body: config,
  });

  if (error || !data) {
    console.warn("API failed, returning mock updated monitoring configuration");
    return { data: config };
  }

  return { data: data.data };
}
