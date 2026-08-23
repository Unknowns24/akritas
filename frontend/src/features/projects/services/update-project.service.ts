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

  if (error || !data) throw error || new Error("No data returned");
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
      dokploy_application: request.dokploy_server_id ? {
        dokploy_server_id: request.dokploy_server_id,
        application_identifier: request.application_identifier || "app",
        instance_identifier: request.application_identifier || "app",
        display_name: "Mock App",
        environment: "production",
        status: "running",
      } : {
        dokploy_server_id: "mock",
        application_identifier: "app",
        instance_identifier: "app",
        display_name: "Mock App",
        environment: "production",
        status: "running",
      },
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
