import { api } from "@/core/libs/api-client";
import { requireApiData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type QvacConfiguration = components["schemas"]["QvacConfiguration"];

export async function getQvacConfigurationService(): Promise<QvacConfiguration> {
  const { data, error } = await api.GET("/integrations/qvac/configuration");
  return requireApiData(data, error).data;
}
