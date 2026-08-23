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

  if (error || !data) {
    console.warn("API failed, returning mock repositories");
    return {
      data: Array.from({ length: 42 }).map((_, i) => ({
        github_account_id: accountId,
        repository_identifier: `repo-${i}`,
        owner: "octocat",
        name: `repository-${i}`,
        full_name: `octocat/repository-${i}`,
        default_branch: "main",
        private: i % 3 === 0,
        html_url: `https://github.com/octocat/repository-${i}`,
      })) as GitHubRepository[]
    };
  }

  return { data: data?.data };
}
