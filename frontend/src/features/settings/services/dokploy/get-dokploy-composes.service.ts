import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type DokployCompose = components["schemas"]["DokployCompose"];

export async function getDokployComposesService(
  serverId: string
): Promise<ServiceData<DokployCompose[]>> {
  const { data, error } = await api.GET("/integrations/dokploy/servers/{server_id}/composes", {
    params: {
      path: { server_id: serverId },
    },
  });

  return { data: requireApiData(data, error).data };
}
