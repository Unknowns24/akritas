import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type LogEventListResponse =
  components["schemas"]["LogEventListResponse"];
export type LogEvent = components["schemas"]["LogEvent"];

export interface GetIncidentLogEventsParams {
  limit?: number;
  cursor?: string;
}

export async function getIncidentLogEventsService(
  incidentId: string,
  params?: GetIncidentLogEventsParams,
): Promise<LogEventListResponse> {
  const { data, error } = await api.GET("/incidents/{incident_id}/log-events", {
    params: {
      path: { incident_id: incidentId },
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
      } as unknown as LogEventListResponse;
    }
    throw error;
  }

  if (!data) {
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
      } as unknown as LogEventListResponse;
    }
    throw new Error("No data returned");
  }

  return data;
}
