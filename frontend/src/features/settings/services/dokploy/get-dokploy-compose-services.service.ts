import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type DokployComposeService = components["schemas"]["DokployComposeService"];

export async function getDokployComposeServicesService(
  serverId: string,
  composeId: string
): Promise<ServiceData<DokployComposeService[]>> {
  const { data, error } = await api.GET("/integrations/dokploy/servers/{server_id}/composes/{compose_id}/services", {
    params: {
      path: { server_id: serverId, compose_id: composeId },
    },
  });

  return { data: requireApiData(data, error).data };
}
