import { api } from "@/core/libs/api-client";
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

  if (error) {
    if (typeof window === "undefined") {
      return {
        data: [],
        paging: {
          limit: 10,
          total: 0,
          has_more: false,
          next_cursor: "",
          prev_cursor: "",
        },
      } as unknown as ProjectListResponse;
    }
    throw error;
  }

  if (!data) {
    if (typeof window === "undefined") {
      return {
        data: [],
        paging: {
          limit: 10,
          total: 0,
          has_more: false,
          next_cursor: "",
          prev_cursor: "",
        },
      } as unknown as ProjectListResponse;
    }
    throw new Error("No data returned");
  }

  return data;
}

