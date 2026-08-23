import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type DokployApplication = components["schemas"]["DokployApplication"];

export async function getDokployApplicationsService(
  serverId: string
): Promise<ServiceData<DokployApplication[]>> {
  const { data, error } = await api.GET("/integrations/dokploy/servers/{server_id}/applications", {
    params: {
      path: { server_id: serverId },
    },
  });

  return { data: requireApiData(data, error).data };
}
