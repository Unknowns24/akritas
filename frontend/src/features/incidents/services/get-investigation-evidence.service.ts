import { api } from "@/core/libs/api-client";
import { requireApiData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type EvidenceListResponse =
  components["schemas"]["EvidenceListResponse"];
export type Evidence = components["schemas"]["Evidence"];

export async function getInvestigationEvidenceService(
  investigationId: string,
): Promise<EvidenceListResponse> {
  const { data, error } = await api.GET(
    "/investigations/{investigation_id}/evidence",
    {
      params: {
        path: { investigation_id: investigationId },
      },
    },
  );

  return requireApiData(data, error);
}
