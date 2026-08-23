import { api } from "@/core/libs/api-client";
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
      } as unknown as EvidenceListResponse;
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
      } as unknown as EvidenceListResponse;
    }
    throw new Error("No data returned");
  }

  return data;
}
