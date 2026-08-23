import { api } from "@/core/libs/api-client";
import { requireApiData } from "@/core/libs/api-client";
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

  return requireApiData(data, error);
}
