import { api } from "@/core/libs/api-client";
import { requireApiData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type ConnectionTestResult = components["schemas"]["ConnectionTestResult"];

export async function testQvacConnectionService(): Promise<ConnectionTestResult> {
  const { data, error } = await api.POST("/integrations/qvac/connection-test");
  return requireApiData(data, error).data;
}
