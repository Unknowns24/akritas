import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type UpdateProjectRequest = components["schemas"]["UpdateProjectRequest"];
type Project = components["schemas"]["Project"];

export async function updateProjectService(
  projectId: string,
  request: UpdateProjectRequest
): Promise<ServiceData<Project>> {
  const { data, error } = await api.PATCH("/projects/{project_id}", {
    params: {
      path: { project_id: projectId },
    },
    body: request,
  });

  return { data: requireApiData(data, error).data };
}
