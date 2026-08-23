import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type MonitoringConfiguration = components["schemas"]["MonitoringConfiguration"];
export type MonitoringConfigurationResponse = components["schemas"]["MonitoringConfigurationResponse"];

export async function updateMonitoringConfigService(
  projectId: string,
  config: MonitoringConfiguration
): Promise<ServiceData<MonitoringConfiguration>> {
  const { data, error } = await api.PUT("/projects/{project_id}/monitoring-configuration", {
    params: {
      path: { project_id: projectId },
    },
    body: config,
  });

  return { data: requireApiData(data, error).data };
}
