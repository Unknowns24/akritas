import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

type CreateRequest = components["schemas"]["CreateGitHubPatAccountRequest"];
export type GitHubAccount = components["schemas"]["GitHubAccount"];

export async function createGitHubPatService(body: CreateRequest): Promise<{ data?: GitHubAccount; error?: Error | any }> {
  const { data, error } = await api.POST("/integrations/github/accounts", {
    body,
  });

  if (error || !data) throw error || new Error("No data returned");
  

  return { data: data?.data };
}
