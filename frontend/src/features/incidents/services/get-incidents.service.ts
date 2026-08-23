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
      return { data: [], paging: { limit: 10, total: 0, has_more: false, next_cursor: "", prev_cursor: "" } } as unknown as IncidentListResponse;
    }
    throw error;
  }

  if (!data) {
    if (typeof window === "undefined") {
      return { data: [], paging: { limit: 10, total: 0, has_more: false, next_cursor: "", prev_cursor: "" } } as unknown as IncidentListResponse;
    }
    throw new Error("No data returned");
  }
  /* [MOCK DOCS]
  if (error || !data) {
    console.warn("Failed to fetch incidents, returning mock data:", error);
    return {
      data: [
        {
          id: "inc-1",
          key: "AKR-184",
          project: { id: "1", name: "E-Commerce Platform" },
          fingerprint: "db_conn_timeout",
          severity: "critical",
          title: "Nil pointer panic on GET /users/:id",
          summary:
            "Multiple instances reporting timeouts when connecting to the primary DB.",
          phase: "detected",
          occurrence_count: 37,
          first_seen_at: new Date(Date.now() - 3600000).toISOString(),
          last_seen_at: new Date().toISOString(),
        },
        {
          id: "inc-2",
          key: "AKR-2",
          project: { id: "2", name: "Payment Gateway" },
          fingerprint: "stripe_api_rate_limit",
          severity: "warning",
          title: "Stripe API Rate Limit Exceeded",
          summary: "Payment processing degraded due to rate limiting.",
          phase: "failed",
          occurrence_count: 5,
          first_seen_at: new Date(Date.now() - 86400000).toISOString(),
          last_seen_at: new Date(Date.now() - 82800000).toISOString(),
        },
      ],
      paging: {
        limit: 10,
        total: 2,
        has_more: false,
        next_cursor: "",
        prev_cursor: "",
      },
    };
  }
  */

  return data;
}
