import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type GitHubAccount = components["schemas"]["GitHubAccount"];

export async function getGitHubAccountsService(): Promise<ServiceData<GitHubAccount[]>> {
  const { data, error } = await api.GET("/integrations/github/accounts", {});

  return { data: requireApiData(data, error).data };
}
