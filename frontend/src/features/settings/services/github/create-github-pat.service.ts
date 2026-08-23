import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

type CreateRequest = components["schemas"]["CreateGitHubPatAccountRequest"];
export type GitHubAccount = components["schemas"]["GitHubAccount"];

export async function createGitHubPatService(body: CreateRequest): Promise<ServiceData<GitHubAccount>> {
  const { data, error } = await api.POST("/integrations/github/accounts", {
    body,
  });

  return { data: requireApiData(data, error).data };
}
