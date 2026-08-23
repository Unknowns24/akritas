import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";
import type { DokployServer } from "./list-dokploy-servers.service";

export type CreateDokployServerRequest = components["schemas"]["CreateDokployServerRequest"];

export async function createDokployServerService(
  payload: CreateDokployServerRequest
): Promise<ServiceData<DokployServer>> {
  const { data, error } = await api.POST("/integrations/dokploy/servers", {
    body: payload,
  });

  return { data: requireApiData(data, error).data };
}
