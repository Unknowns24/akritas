import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";
import type { DokployServer } from "./list-dokploy-servers.service";

export type UpdateDokployServerRequest = components["schemas"]["UpdateDokployServerRequest"];

export async function updateDokployServerService(
  serverId: string,
  payload: UpdateDokployServerRequest
): Promise<ServiceData<DokployServer>> {
  const { data, error } = await api.PATCH("/integrations/dokploy/servers/{server_id}", {
    params: {
      path: { server_id: serverId },
    },
    body: payload,
  });

  return { data: requireApiData(data, error).data };
}
