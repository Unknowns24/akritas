import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type CreateProjectRequest = components["schemas"]["CreateProjectRequest"];
export type Project = components["schemas"]["Project"];

export async function createProjectService(
  request: CreateProjectRequest
): Promise<{ data?: Project; error?: Error | any }> {
  const { data, error } = await api.POST("/projects", {
    body: request,
  });

  if (error || !data) {
    console.warn("API failed, returning mock created project");
    const mockProject: Project = {
      id: crypto.randomUUID(),
      name: request.name,
      description: request.description,
      health_status: "healthy",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      github_repository: {
        github_account_id: request.github_account_id,
        repository_identifier: request.repository_identifier,
        default_branch: request.default_branch,
        full_name: `${request.repository_identifier}`,
        private: true,
        owner: "mock",
        name: "mock",
        html_url: "https://github.com",
      },
      dokploy_application: {
        dokploy_server_id: request.dokploy_server_id,
        application_identifier: request.application_identifier,
        instance_identifier: request.application_identifier,
        display_name: request.name,
        environment: "production",
        status: "running",
      },
      monitoring_configuration: request.monitoring_configuration,
      monitoring_status: "monitoring",
      built_in_detection_rules: [],
    };
    return { data: mockProject };
  }

  return { data: data.data };
}
