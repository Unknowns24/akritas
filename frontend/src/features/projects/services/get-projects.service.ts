import { api } from "@/core/libs/api-client";
import { requireApiData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type ProjectListResponse =
  components["schemas"]["ProjectListResponse"];

export interface GetProjectsParams {
  limit?: number;
  cursor?: string;
  name_like?: string;
}

export async function getProjectsService(
  params?: GetProjectsParams,
): Promise<ProjectListResponse> {
  const queryParams = params ? {
    ...(params.limit !== undefined ? { limit: params.limit } : {}),
    ...(params.cursor ? { cursor: params.cursor } : {}),
    ...(params.name_like ? { name_like: params.name_like } : {}),
  } : undefined;

  const { data, error } = await api.GET("/projects", {
    ...(queryParams && Object.keys(queryParams).length > 0 ? { params: { query: queryParams } } : {}),
  });

  return requireApiData(data, error);
}

