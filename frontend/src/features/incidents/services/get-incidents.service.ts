import { api } from "@/core/libs/api-client";
import { requireApiData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type IncidentListResponse =
  components["schemas"]["IncidentListResponse"];

export interface GetIncidentsParams {
  limit?: number;
  cursor?: string;
}

export async function getIncidentsService(
  params?: GetIncidentsParams,
): Promise<IncidentListResponse> {
  const { data, error } = await api.GET("/incidents", {
    params: {
      query: {
        limit: params?.limit,
        cursor: params?.cursor,
      },
    },
  });

  return requireApiData(data, error);
}
