import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type ProjectListResponse = components["schemas"]["ProjectListResponse"];

export async function getProjectsService(): Promise<ProjectListResponse> {
  const { data, error } = await api.GET("/projects");

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
          dokploy_application: {
            dokploy_server_id: "server-1",
            application_identifier: "app-1",
            instance_identifier: "inst-1",
            display_name: "ecommerce-api",
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
          dokploy_application: {
            dokploy_server_id: "server-1",
            application_identifier: "app-2",
            instance_identifier: "inst-2",
            display_name: "payments-api",
            status: "running",
          },
          created_at: new Date().toISOString(),
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

  return data;
}
