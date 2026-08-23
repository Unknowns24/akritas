import { api } from "@/core/libs/api-client";
import { requireApiData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type SystemStatusResponse =
  components["schemas"]["SystemStatusResponse"];

export async function getSystemStatusService(): Promise<SystemStatusResponse> {
  const { data, error } = await api.GET("/system/status", {});
  return requireApiData(data, error);
}
