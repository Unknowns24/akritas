import { api } from "@/core/libs/api-client";
import { requireApiData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type OverviewResponse = components["schemas"]["OverviewResponse"];

export async function getOverviewService(): Promise<OverviewResponse> {
  const { data, error } = await api.GET("/overview");
  return requireApiData(data, error);
}
