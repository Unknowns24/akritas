import { api } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type GitHubManifestRegistrationRequest =
  components["schemas"]["GitHubManifestRegistrationRequest"];
export type GitHubManifestRegistration =
  components["schemas"]["GitHubManifestRegistration"];

export async function startGitHubAppManifestRegistrationService(
  payload: GitHubManifestRegistrationRequest,
): Promise<{ data?: GitHubManifestRegistration; error?: Error | any }> {
  const { data, error } = await api.POST(
    "/integrations/github/app-manifest/registrations",
    {
      body: payload,
    },
  );

  if (error || !data) throw error || new Error("No data returned");
  /* [MOCK DOCS]
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
  */
  if (data?.data?.manifest) {
    try {
      const manifestObj = JSON.parse(data.data.manifest);

      const baseUrl =
        typeof window !== "undefined"
          ? window.location.origin
          : "http://localhost:3000";

      // Aggressively inject all required URLs for GitHub App Manifest
      manifestObj.url = manifestObj.url || baseUrl;

      if (!manifestObj.hook_attributes) {
        manifestObj.hook_attributes = {};
      }

      // GitHub strictly rejects 'localhost' for webhook URLs.
      // We use a dummy public URL for local development so the app creation succeeds.
      // The user can update it later in GitHub settings if they set up ngrok.
      manifestObj.hook_attributes.url =
        "https://example.com/webhook-placeholder";

      // GitHub DOES allow localhost for the OAuth redirect callback
      manifestObj.redirect_url =
        manifestObj.redirect_url || `${baseUrl}/settings/github/callback`;
      manifestObj.public = false;

      if (!manifestObj.name) {
        manifestObj.name = `akritas-dev-${Math.floor(Math.random() * 100000)}`;
      }

      const finalManifestStr = JSON.stringify(manifestObj);

      // Return a completely new object to avoid read-only mutation issues
      return {
        data: {
          ...data.data,
          manifest: finalManifestStr,
        },
      };
    } catch (e) {
      console.warn("Could not parse or mutate manifest JSON", e);
    }
  }

  return { data: data?.data };
}
