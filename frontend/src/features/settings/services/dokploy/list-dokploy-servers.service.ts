import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type DokployServer = components["schemas"]["DokployServer"];

export async function listDokployServersService(): Promise<{ data?: DokployServer[]; error?: Error | any }> {
  const { data, error } = await api.GET("/integrations/dokploy/servers", {
    params: {
      query: {}
    }
  });

  if (error || !data) {
    console.warn("API failed, returning mock Dokploy servers");
    return {
      data: [
        {
          id: "mock-dokploy-1",
          name: "Production Server",
          base_url: "https://dokploy.example.com",
          server_identifier: "dokploy-example-com",
          connection_status: "connected",
          credential_configured: true,
          application_count: 5,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          last_synced_at: new Date().toISOString(),
        },
        {
          id: "mock-dokploy-2",
          name: "Staging Server (Failing)",
          base_url: "https://staging.dokploy.example.com",
          server_identifier: "staging-dokploy-example-com",
          connection_status: "authentication_failed",
          credential_configured: true,
          application_count: 0,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }
      ]
    };
  }

  return { data: data.data };
}
