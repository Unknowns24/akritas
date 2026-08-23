import { api } from "@/core/libs/api-client";
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

  if (error) {
    if (typeof window === "undefined") {
      return {
        data: [],
        paging: {
          limit: 10,
          total: 0,
          has_more: false,
          next_cursor: "",
          prev_cursor: "",
        },
      } as unknown as IncidentListResponse;
    }
    throw error;
  }

  if (!data) {
    if (typeof window === "undefined") {
      return { data: [], paging: { limit: 10, total: 0, has_more: false, next_cursor: "", prev_cursor: "" } } as unknown as IncidentListResponse;
    }
    throw new Error("No data returned");
  }


  return data;
}
