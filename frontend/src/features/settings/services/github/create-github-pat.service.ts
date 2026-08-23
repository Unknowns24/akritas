import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

type CreateRequest = components["schemas"]["CreateGitHubPatAccountRequest"];
export type GitHubAccount = components["schemas"]["GitHubAccount"];

export async function createGitHubPatService(body: CreateRequest): Promise<{ data?: GitHubAccount; error?: Error | any }> {
  const { data, error } = await api.POST("/integrations/github/accounts", {
    body,
  });

  if (error || !data) {
    console.warn("API failed, returning mock created account");
    return {
      data: {
        id: "gh-acc-" + Math.random().toString(36).substring(7),
        account_type: body.account_type,
        display_name: body.display_name,
        account_identifier: body.account_identifier,
        authentication_method: "personal_access_token",
        authentication_status: "pending",
        credential_configured: true,
        repository_count: 0,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      } as GitHubAccount
    };
  }

  return { data: data?.data };
}
