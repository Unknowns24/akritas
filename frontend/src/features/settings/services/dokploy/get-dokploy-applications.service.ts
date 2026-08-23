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
  /* [MOCK DOCS]
  if (error || !data) {
    console.warn(`API failed fetching applications for ${serverId}, returning mock applications`);
    
    return {
      data: [
        {
          dokploy_server_id: serverId,
          application_identifier: "app-12345",
          instance_identifier: "prod-cluster-1",
          display_name: "Frontend Application",
          environment: "production",
          status: "running"
        },
        {
          dokploy_server_id: serverId,
          application_identifier: "app-67890",
          instance_identifier: "prod-cluster-2",
          display_name: "Backend API",
          environment: "production",
          status: "degraded"
        },
        {
          dokploy_server_id: serverId,
          application_identifier: "app-00000",
          instance_identifier: "staging-cluster-1",
          display_name: "Staging Testing Environment",
          environment: "staging",
          status: "stopped"
        }
      ]
    };
  }
  */

  return { data: data.data };
}
