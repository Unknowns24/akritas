import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

type UpdateRequest = components["schemas"]["UpdateGitHubAccountRequest"];
export type GitHubAccount = components["schemas"]["GitHubAccount"];

export async function updateGitHubPatService(
  accountId: string,
  body: UpdateRequest
): Promise<{ data?: GitHubAccount; error?: Error | any }> {
  const { data, error } = await api.PATCH("/integrations/github/accounts/{account_id}", {
    params: {
      path: { account_id: accountId },
    },
    body,
  });

  if (error || !data) throw error || new Error("No data returned");
  

  return { data: data?.data };
}
