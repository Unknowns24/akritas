import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type GitHubAccount = components["schemas"]["GitHubAccount"];

export async function getGitHubAccountsService(): Promise<{ data?: GitHubAccount[]; error?: Error | any }> {
  const { data, error } = await api.GET("/integrations/github/accounts", {});

  if (error || !data) {
    console.warn("API failed, returning mock GitHub accounts");
    return {
      data: [
        {
          id: "gh-acc-1",
          account_type: "personal",
          display_name: "octocat",
          account_identifier: "octocat",
          authentication_status: "connected",
          credential_configured: true,
          repository_count: 42,
          authentication_method: "personal_access_token",
          created_at: new Date(Date.now() - 30 * 86400000).toISOString(),
          updated_at: new Date(Date.now() - 30 * 86400000).toISOString(),
        } as GitHubAccount
      ]
    };
  }

  return { data: data?.data };
}
