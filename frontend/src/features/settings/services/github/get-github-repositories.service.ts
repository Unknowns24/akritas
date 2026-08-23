import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type GitHubRepository = components["schemas"]["GitHubRepository"];

export async function getGitHubRepositoriesService(
  accountId: string
): Promise<{ data?: GitHubRepository[]; error?: Error | any }> {
  const { data, error } = await api.GET("/integrations/github/accounts/{account_id}/repositories", {
    params: {
      path: { account_id: accountId },
    },
  });

  if (error || !data) throw error || new Error("No data returned");
  

  return { data: data?.data };
}
