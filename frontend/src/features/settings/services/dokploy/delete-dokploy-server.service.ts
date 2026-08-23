import { api } from "@/core/libs/api-client";

export async function deleteDokployServerService(
  serverId: string
): Promise<{ success: boolean }> {
  const { error } = await api.DELETE("/integrations/dokploy/servers/{server_id}", {
    params: {
      path: { server_id: serverId },
    },
  });

  if (error) throw error;

  return { success: true };
}
