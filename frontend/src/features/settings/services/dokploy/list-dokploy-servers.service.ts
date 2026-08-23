import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type DokployServer = components["schemas"]["DokployServer"];

export async function listDokployServersService(): Promise<ServiceData<DokployServer[]>> {
  const { data, error } = await api.GET("/integrations/dokploy/servers", {
    params: {
      query: {}
    }
  });

  return { data: requireApiData(data, error).data };
}
