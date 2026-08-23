import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type UpdateProjectRequest = components["schemas"]["UpdateProjectRequest"];
type Project = components["schemas"]["Project"];

export async function updateProjectService(
  projectId: string,
  request: UpdateProjectRequest
): Promise<{ data?: Project; error?: Error | any }> {
  const { data, error } = await api.PATCH("/projects/{project_id}", {
    params: {
      path: { project_id: projectId },
    },
    body: request,
  });

  if (error || !data) {
    return { error: error || new Error("No data returned") };
  }
  /* [MOCK DOCS]
  if (error || !data) {
    console.warn("API failed, returning mock updated project");
    const mockProject: Project = {
      id: projectId,
      name: request.name || "Updated Mock Project",
      description: request.description || "Updated description",
      health_status: "healthy",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      github_repository: request.github_account_id ? {
        github_account_id: request.github_account_id,
        repository_identifier: request.repository_identifier || "repo",
        default_branch: request.default_branch || "main",
        full_name: `${request.repository_identifier || "repo"}`,
        private: true,
        owner: "mock",
        name: "mock",
        html_url: "https://github.com",
      } : {
        github_account_id: "mock",
        repository_identifier: "repo",
        default_branch: "main",
        full_name: "repo",
        private: true,
        owner: "mock",
        name: "mock",
        html_url: "https://github.com",
      },
      dokploy_source: request.dokploy_source ? {
        type: request.dokploy_source.type,
        dokploy_server_id: request.dokploy_source.dokploy_server_id,
        resource_identifier: request.dokploy_source.resource_identifier,
        instance_identifier: request.dokploy_source.resource_identifier,
        display_name: "Mock Source",
        environment: "production",
        status: "running",
        ...(request.dokploy_source.type === "compose_service" ? {
          service_name: (request.dokploy_source as any).service_name,
          runtime_type: "docker-compose"
        } : {})
      } as any : {
        type: "application",
        dokploy_server_id: "mock",
        resource_identifier: "app",
        instance_identifier: "app",
        display_name: "Mock Source",
        environment: "production",
        status: "running",
      } as any,
      monitoring_configuration: {
        enabled: true,
        error_patterns: [],
        ignored_patterns: [],
        grouping_window: "PT30M",
        context_before: 20,
        context_after: 20,
      },
      monitoring_status: "monitoring",
      built_in_detection_rules: [],
    };
    return { data: mockProject };
  }
  */

  return { data: data.data };
}
