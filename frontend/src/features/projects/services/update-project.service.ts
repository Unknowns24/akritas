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
  

  return { data: data.data };
}
