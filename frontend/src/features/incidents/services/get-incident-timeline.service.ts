import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type TimelineEventListResponse =
  components["schemas"]["TimelineEventListResponse"];
export type TimelineEvent = components["schemas"]["TimelineEvent"];

export async function getIncidentTimelineService(
  incidentId: string,
): Promise<TimelineEventListResponse> {
  const { data, error } = await api.GET("/incidents/{incident_id}/timeline", {
    params: {
      path: { incident_id: incidentId },
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
      } as unknown as TimelineEventListResponse;
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
      } as unknown as TimelineEventListResponse;
    }
    throw new Error("No data returned");
  }

  return data;
}
