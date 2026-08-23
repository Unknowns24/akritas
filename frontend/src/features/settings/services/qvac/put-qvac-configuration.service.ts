import { api } from "@/core/libs/api-client";
import { requireApiData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";
import type { QvacConfiguration } from "./get-qvac-configuration.service";

export type PutQvacConfigurationRequest = components["schemas"]["PutQvacConfigurationRequest"];

export async function putQvacConfigurationService(
  body: PutQvacConfigurationRequest
): Promise<QvacConfiguration> {
  const { data, error } = await api.PUT("/integrations/qvac/configuration", {
    body,
  });
  return requireApiData(data, error).data;
}
