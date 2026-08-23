import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type ConnectionTestResponse = components["schemas"]["ConnectionTestResponse"];

export async function testDokployConnectionService(
  serverId: string
): Promise<{ data?: ConnectionTestResponse; error?: Error | any }> {
  const { data, error } = await api.POST("/integrations/dokploy/servers/{server_id}/connection-test", {
    params: {
      path: { server_id: serverId },
    },
  });

  if (error || !data) throw error || new Error("No data returned");
  

  return { data };
}
