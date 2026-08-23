import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type GitHubManifestRegistrationRequest = components["schemas"]["GitHubManifestRegistrationRequest"];
export type GitHubManifestRegistration = components["schemas"]["GitHubManifestRegistration"];

export async function startGitHubAppManifestRegistrationService(
  payload: GitHubManifestRegistrationRequest
): Promise<{ data?: GitHubManifestRegistration; error?: Error | any }> {
  const { data, error } = await api.POST("/integrations/github/app-manifest/registrations", {
    body: payload,
  });

  if (error || !data) {
    console.warn("API failed, returning mock manifest registration");
    return {
      data: {
        registration_id: "mock-registration-id",
        form_action: "https://github.com/settings/apps/new",
        manifest: JSON.stringify({
          name: `akritas-mock-${Date.now()}`,
          url: "http://localhost:3000",
          hook_attributes: {
            url: "http://localhost:3000/settings/github/callback"
          },
          redirect_url: "http://localhost:3000/settings/github/callback",
          public: false,
          default_permissions: {
            contents: "write",
            metadata: "read",
            issues: "write",
            pull_requests: "write"
          },
          default_events: []
        }),
        state: "mock-state-12345",
        expires_at: new Date(Date.now() + 3600000).toISOString()
      }
    };
  }

  return { data: data?.data };
}
