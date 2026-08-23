import { api } from "@/core/libs/api-client";

export async function testGitHubConnectionService(
  accountId: string
): Promise<{ success: boolean }> {
  const { error } = await api.POST("/integrations/github/accounts/{account_id}/connection-test", {
    params: {
      path: { account_id: accountId },
    },
  });

  if (error) throw error;

  return { success: true };
}
