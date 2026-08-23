import { api } from "@/core/libs/api-client";
import { requireApiData, type ServiceData } from "@/core/libs/api-client";
import type { components } from "@/core/libs/api-client";

export type GitHubManifestRegistrationRequest =
  components["schemas"]["GitHubManifestRegistrationRequest"];
export type GitHubManifestRegistration =
  components["schemas"]["GitHubManifestRegistration"];

export async function startGitHubAppManifestRegistrationService(
  payload: GitHubManifestRegistrationRequest,
): Promise<ServiceData<GitHubManifestRegistration>> {
  const response = await api.POST(
    "/integrations/github/app-manifest/registrations",
    {
      body: payload,
    },
  );
  const envelope = requireApiData(response.data, response.error);

  if (envelope.data.manifest) {
    try {
      const manifestObj: Record<string, unknown> = JSON.parse(envelope.data.manifest);

      const baseUrl =
        typeof window !== "undefined"
          ? window.location.origin
          : "http://localhost:3000";

      // Aggressively inject all required URLs for GitHub App Manifest
      manifestObj.url = manifestObj.url || baseUrl;

      if (
        typeof manifestObj.hook_attributes !== "object" ||
        manifestObj.hook_attributes === null
      ) {
        manifestObj.hook_attributes = {};
      }
      const hookAttributes = manifestObj.hook_attributes as Record<string, unknown>;

      // GitHub strictly rejects 'localhost' for webhook URLs.
      // We use a dummy public URL for local development so the app creation succeeds.
      // The user can update it later in GitHub settings if they set up ngrok.
      hookAttributes.url = "https://example.com/webhook-placeholder";

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
          ...envelope.data,
          manifest: finalManifestStr,
        },
      };
    } catch (error) {
      console.warn("Could not parse or mutate manifest JSON", error);
    }
  }

  return { data: envelope.data };
}
