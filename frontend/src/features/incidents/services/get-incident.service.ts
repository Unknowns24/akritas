import { api } from "@/core/libs/api-client";
import { requireApiData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type Incident = components["schemas"]["Incident"];

export async function getIncidentService(id: string): Promise<Incident> {
  const { data, error } = await api.GET("/incidents/{incident_id}", {
    params: {
      path: { incident_id: id },
    },
  });

  return requireApiData(data, error).data;
}
