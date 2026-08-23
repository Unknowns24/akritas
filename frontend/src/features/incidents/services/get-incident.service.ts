import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type Incident = components["schemas"]["Incident"];

export async function getIncidentService(id: string): Promise<Incident> {
  const { data, error } = await api.GET("/incidents/{incident_id}", {
    params: {
      path: { incident_id: id },
    },
  });

  if (error) {
    if (typeof window === "undefined") {
      return {} as unknown as Incident;
    }
    throw error;
  }

  if (!data) {
    if (typeof window === "undefined") {
      return {} as unknown as Incident;
    }
    throw new Error("No data returned");
  }

  return data.data;
}
