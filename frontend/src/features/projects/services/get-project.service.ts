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


  return data;
}
