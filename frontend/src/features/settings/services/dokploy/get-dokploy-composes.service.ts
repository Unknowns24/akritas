import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type DokployCompose = components["schemas"]["DokployCompose"];

export async function getDokployComposesService(
  serverId: string
): Promise<{ data?: DokployCompose[]; error?: Error | any }> {
  const { data, error } = await api.GET("/integrations/dokploy/servers/{server_id}/composes", {
    params: {
      path: { server_id: serverId },
    },
  });

  if (error || !data) throw error || new Error("No data returned");

  return { data: data.data };
}
