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

  if (error || !data) throw error || new Error("No data returned");
  

  return { data: data.data };
}
