import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type GitHubAccount = components["schemas"]["GitHubAccount"];

export async function getGitHubAccountsService(): Promise<{
  data?: GitHubAccount[];
  error?: Error | any;
}> {
  const { data, error } = await api.GET("/integrations/github/accounts", {});

  if (error) {
    // In Server Components (SSR), cross-domain cookies aren't sent to localhost.
    // Instead of redirecting or crashing, return empty data gracefully.
    // The client component's useEffect will re-fetch with the browser's cookie and populate the list.
    if (typeof window === "undefined") {
      return { data: [] };
    }
    return { error };
  }

  if (!data) {
    if (typeof window === "undefined") return { data: [] };
    return { error: new Error("No data returned") };
  }
  

  return { data: data?.data };
}
