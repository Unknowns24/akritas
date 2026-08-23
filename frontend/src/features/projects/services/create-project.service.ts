import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type CreateProjectRequest = components["schemas"]["CreateProjectRequest"];
type Project = components["schemas"]["Project"];

export async function createProjectService(
  request: CreateProjectRequest
): Promise<ServiceData<Project>> {
  const { data, error } = await api.POST("/projects", {
    body: request,
  });

  return { data: requireApiData(data, error).data };
}
