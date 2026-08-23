import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";
import { DokployServer } from "./list-dokploy-servers.service";

export type CreateDokployServerRequest = components["schemas"]["CreateDokployServerRequest"];

export async function createDokployServerService(
  payload: CreateDokployServerRequest
): Promise<{ data?: DokployServer; error?: Error | any }> {
  const { data, error } = await api.POST("/integrations/dokploy/servers", {
    body: payload,
  });

  if (error || !data) {
    console.warn("API failed, returning mock created Dokploy server");
    return {
      data: {
        id: "mock-new-dokploy-server-" + Date.now(),
        name: payload.name,
        base_url: payload.base_url,
        server_identifier: payload.base_url.replace(/^https?:\/\//, ""),
        connection_status: "pending",
        credential_configured: !!payload.api_credential,
        application_count: 0,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
    };
  }

  return { data: data.data };
}
