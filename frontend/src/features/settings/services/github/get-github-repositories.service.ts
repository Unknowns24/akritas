import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type GitHubRepository = components["schemas"]["GitHubRepository"];

export async function getGitHubRepositoriesService(
  accountId: string
): Promise<ServiceData<GitHubRepository[]>> {
  const { data, error } = await api.GET("/integrations/github/accounts/{account_id}/repositories", {
    params: {
      path: { account_id: accountId },
    },
  });

  return { data: requireApiData(data, error).data };
}
