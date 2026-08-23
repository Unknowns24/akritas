import { GitHubSettingsClient } from "@/features/settings/views/GitHubSettingsView";
import { getGitHubAccountsService } from "@/features/settings/services/github/get-github-accounts.service";

export default async function GitHubSettingsPage() {
  const { data: accounts } = await getGitHubAccountsService();
  
  return (
    <GitHubSettingsClient initialAccounts={accounts || []} />
  );
}
