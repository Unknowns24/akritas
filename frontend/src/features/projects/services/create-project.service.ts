import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type CreateProjectRequest = components["schemas"]["CreateProjectRequest"];
type Project = components["schemas"]["Project"];

export async function createProjectService(
  request: CreateProjectRequest
): Promise<{ data?: Project; error?: Error | any }> {
  const { data, error } = await api.POST("/projects", {
    body: request,
  });

  if (error || !data) {
    return { error: error || new Error("No data returned") };
  }
  

  return { data: data.data };
}
