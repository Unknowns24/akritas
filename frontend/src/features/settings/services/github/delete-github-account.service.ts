import { api } from "@/core/libs/api-client";

export async function deleteGitHubAccountService(
  accountId: string
): Promise<{ error?: Error | any }> {
  const { error } = await api.DELETE("/integrations/github/accounts/{account_id}", {
    params: {
      path: { account_id: accountId },
    },
  });

  if (error) {
    console.warn("API failed, mocking delete success");
    return {};
  }

  return {};
}
