import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type ProjectResponse = components["schemas"]["ProjectResponse"];
export type Project = components["schemas"]["Project"];

export async function getProjectService(id: string): Promise<ProjectResponse> {
  const { data, error } = await api.GET("/projects/{project_id}", {
    params: {
      path: { project_id: id },
    },
  });

  if (error || !data) throw error || new Error("No data returned");
  /* [MOCK DOCS]
  if (error || !data) {
    console.warn(`Failed to fetch project ${id}, returning mock data:`, error);
    return {
      data: {
        id: id,
        name: id === "1" ? "E-Commerce Platform" : "Payment Gateway",
        description: id === "1" ? "Main production e-commerce backend" : "Stripe payment integration service",
        monitoring_status: id === "1" ? "monitoring" : "degraded",
        health_status: id === "1" ? "healthy" : "warning",
        github_repository: {
          github_account_id: "acc-1",
          repository_identifier: "repo-1",
          owner: "akritas",
          name: id === "1" ? "ecommerce" : "payments",
          full_name: id === "1" ? "akritas/ecommerce" : "akritas/payments",
          default_branch: "main",
          private: true,
          html_url: `https://github.com/akritas/${id === "1" ? "ecommerce" : "payments"}`,
        },
        dokploy_application: {
          dokploy_server_id: "server-1",
          application_identifier: "app-1",
          instance_identifier: "inst-1",
          display_name: id === "1" ? "ecommerce-api" : "payments-api",
          status: "running",
          environment: "production"
        },
        monitoring_configuration: {
          enabled: true,
          error_patterns: [],
          ignored_patterns: [],
          grouping_window: "PT30M",
          context_before: 20,
          context_after: 20,
        },
        built_in_detection_rules: [
          {
            code: "http_5xx",
            display_name: "High HTTP 5xx Error Rate",
            enabled: true,
          },
          {
            code: "panic",
            display_name: "Application Panic / Crash",
            enabled: true,
          },
        ],
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        last_observed_at: new Date().toISOString(),
      }
    };
  }
  */

  return data;
}
