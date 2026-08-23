import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type DokployComposeService = components["schemas"]["DokployComposeService"];

export async function getDokployComposeServicesService(
  serverId: string,
  composeId: string
): Promise<{ data?: DokployComposeService[]; error?: Error | any }> {
  const { data, error } = await api.GET("/integrations/dokploy/servers/{server_id}/composes/{compose_id}/services", {
    params: {
      path: { server_id: serverId, compose_id: composeId },
    },
  });

  if (error || !data) throw error || new Error("No data returned");

  return { data: data.data };
}
