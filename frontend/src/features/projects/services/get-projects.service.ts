import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type ProjectListResponse =
  components["schemas"]["ProjectListResponse"];

export interface GetProjectsParams {
  limit?: number;
  cursor?: string;
  name_like?: string;
}

export async function getProjectsService(
  params?: GetProjectsParams,
): Promise<ProjectListResponse> {
  const queryParams = params ? {
    ...(params.limit !== undefined ? { limit: params.limit } : {}),
    ...(params.cursor ? { cursor: params.cursor } : {}),
    ...(params.name_like ? { name_like: params.name_like } : {}),
  } : undefined;

  const { data, error } = await api.GET("/projects", {
    ...(queryParams && Object.keys(queryParams).length > 0 ? { params: { query: queryParams } } : {}),
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
      } as unknown as ProjectListResponse;
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
      } as unknown as ProjectListResponse;
    }
    throw new Error("No data returned");
  }

  return data;
}
/* [MOCK DOCS]
  if (error || !data) {
    console.warn("Failed to fetch projects, returning mock data:", error);
    return {
      data: [
        {
          id: "1",
          name: "E-Commerce Platform",
          description: "Main production e-commerce backend",
          monitoring_status: "monitoring",
          health_status: "healthy",
          github_repository: {
            github_account_id: "acc-1",
            repository_identifier: "repo-1",
            owner: "akritas",
            name: "ecommerce",
            full_name: "akritas/ecommerce",
            default_branch: "main",
            private: true,
            html_url: "https://github.com/akritas/ecommerce",
          },
            dokploy_source: {
              type: "application",
              dokploy_server_id: "server-1",
              resource_identifier: "app-1",
              instance_identifier: "app-1",
              display_name: "Frontend Service",
              environment: "production",
              status: "running",
            },
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        {
          id: "2",
          name: "Payment Gateway",
          description: "Stripe payment integration service",
          monitoring_status: "degraded",
          health_status: "warning",
          github_repository: {
            github_account_id: "acc-1",
            repository_identifier: "repo-2",
            owner: "akritas",
            name: "payments",
            full_name: "akritas/payments",
            default_branch: "main",
            private: true,
            html_url: "https://github.com/akritas/payments",
          },
            dokploy_source: {
              type: "application",
              dokploy_server_id: "server-1",
              resource_identifier: "app-2",
              instance_identifier: "app-2",
              display_name: "Backend API",
              environment: "production",
              status: "degraded",
            },created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
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
