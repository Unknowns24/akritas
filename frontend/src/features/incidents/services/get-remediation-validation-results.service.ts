import { api } from "@/core/libs/api-client";
import type { ValidationResult } from "../types/remediation.types";

export interface GetValidationResultsParams {
  remediationId: string;
  limit?: number;
  cursor?: string;
  sort?: string;
}

export async function getRemediationValidationResultsService({
  remediationId,
  limit,
  cursor,
  sort,
}: GetValidationResultsParams): Promise<ValidationResult[]> {
  const { data, error } = await api.GET("/remediations/{remediation_id}/validation-results", {
    params: {
      path: { remediation_id: remediationId },
      query: { limit, cursor, sort },
    },
  });

  if (error) {
    if (typeof window === "undefined") {
      return [];
    }
    throw error;
  }

  if (!data) {
    return [];
  }

  return data.data;
}
