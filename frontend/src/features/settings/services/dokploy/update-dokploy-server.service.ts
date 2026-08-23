import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";
import { DokployServer } from "./list-dokploy-servers.service";

export type UpdateDokployServerRequest = components["schemas"]["UpdateDokployServerRequest"];

export async function updateDokployServerService(
  serverId: string,
  payload: UpdateDokployServerRequest
): Promise<{ data?: DokployServer; error?: Error | any }> {
  const { data, error } = await api.PATCH("/integrations/dokploy/servers/{server_id}", {
    params: {
      path: { server_id: serverId },
    },
    body: payload,
  });

  if (error || !data) throw error || new Error("No data returned");
  

  return { data: data.data };
}
