import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

type Operation = components["schemas"]["Operation"];

function randomUUID(): string {
  if (globalThis.crypto?.randomUUID) {
    return globalThis.crypto.randomUUID();
  }

  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (char) => {
    const value = Math.floor(Math.random() * 16);
    const replacement = char === "x" ? value : (value & 0x3) | 0x8;
    return replacement.toString(16);
  });
}

export async function startIncidentInvestigationService(
  incidentId: string
): Promise<ServiceData<Operation>> {
  const { data, error } = await api.POST("/incidents/{incident_id}/investigations", {
    params: {
      header: { "Idempotency-Key": randomUUID() },
      path: { incident_id: incidentId },
    },
  });

  return { data: requireApiData(data, error).data };
}
