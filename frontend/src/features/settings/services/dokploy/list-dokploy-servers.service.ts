import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type DokployServer = components["schemas"]["DokployServer"];

export async function listDokployServersService(): Promise<{ data?: DokployServer[]; error?: Error | any }> {
  const { data, error } = await api.GET("/integrations/dokploy/servers", {
    params: {
      query: {}
    }
  });

  if (error || !data) throw error || new Error("No data returned");
  

  return { data: data.data };
}
