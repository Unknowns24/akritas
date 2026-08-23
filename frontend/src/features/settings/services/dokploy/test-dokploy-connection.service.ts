import { api } from "@/core/libs/api-client";
import { requireApiData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type ConnectionTestResult = components["schemas"]["ConnectionTestResult"];

export async function testDokployConnectionService(
  serverId: string
): Promise<{ data: ConnectionTestResult }> {
  const { data, error } = await api.POST("/integrations/dokploy/servers/{server_id}/connection-test", {
    params: {
      path: { server_id: serverId },
    },
  });

  return { data: requireApiData(data, error).data };
}
