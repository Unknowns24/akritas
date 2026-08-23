import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type SystemStatusResponse =
  components["schemas"]["SystemStatusResponse"];

export async function getSystemStatusService(): Promise<SystemStatusResponse> {
  try {
    const { data, error } = await api.GET("/system/status", {});
    if (!error && data) return data;
  } catch (e) {
    // Ignorar si falla, hacemos fallback a un estado vacío para no romper la pantalla
  }

  console.warn(
    "System status endpoint unavailable or returned an error, falling back to empty state...",
  );

  return {
    data: {
      github_account_count: 0,
      dokploy_server_count: 0,
      components: [],
    },
  };
}
