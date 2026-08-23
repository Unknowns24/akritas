import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type DokployApplication = components["schemas"]["DokployApplication"];

export async function getDokployApplicationsService(
  serverId: string
): Promise<{ data?: DokployApplication[]; error?: Error | any }> {
  const { data, error } = await api.GET("/integrations/dokploy/servers/{server_id}/applications", {
    params: {
      path: { server_id: serverId },
    },
  });

  if (error || !data) throw error || new Error("No data returned");
  

  return { data: data.data };
}
