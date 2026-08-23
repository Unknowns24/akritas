import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";
import { DokployServer } from "./list-dokploy-servers.service";

export type UpdateDokployServerRequest = components["schemas"]["UpdateDokployServerRequest"];

export async function updateDokployServerService(
  serverId: string,
  payload: UpdateDokployServerRequest
): Promise<{ data?: DokployServer; error?: Error | any }> {
  const { data, error } = await api.PATCH("/integrations/dokploy/servers/{server_id}", {
    params: {
      path: { server_id: serverId },
    },
    body: payload,
  });

  if (error || !data) throw error || new Error("No data returned");
  /* [MOCK DOCS]
  if (error || !data) {
    console.warn("API failed, returning mock updated Dokploy server");
    return {
      data: {
        id: serverId,
        name: payload.name || "Updated Mock Server",
        base_url: payload.base_url || "https://updated.example.com",
        server_identifier: "updated-mock",
        connection_status: "connected",
        credential_configured: true,
        application_count: 5,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
    };
  }
  */

  return { data: data.data };
}
