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
  /* [MOCK DOCS]
  if (error || !data) {
    console.warn(`API failed testing connection for ${serverId}, returning mock success`);
    
    // Simulate network delay
    await new Promise(resolve => setTimeout(resolve, 800));

    return {
      data: {
        data: {
          status: "connected",
          checked_at: new Date().toISOString(),
          latency_ms: 120,
          user_message: "Mock: Successfully connected to Dokploy API and verified credentials."
        }
      }
    };
  }
  */

  return { data };
}
