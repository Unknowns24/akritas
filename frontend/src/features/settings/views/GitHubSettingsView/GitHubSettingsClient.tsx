"use client";

import React, { useState, useEffect } from "react";
import styles from "./GitHubSettingsClient.module.css";
import { Button } from "@/core/ui/primitives/Button";
import { Modal } from "@/core/ui/primitives/Modal";
import { Badge } from "@/core/ui/primitives/Badge";
import { toast } from "sonner";
import { GitHubAccountForm } from "./components/GitHubAccountForm";
import { Plus, Trash2, Edit2, RefreshCw, CheckCircle, AlertCircle } from "lucide-react";
import { GitHubAccount, getGitHubAccountsService } from "../../services/github/get-github-accounts.service";
import { createGitHubPatService } from "../../services/github/create-github-pat.service";
import { updateGitHubPatService } from "../../services/github/update-github-pat.service";
import { deleteGitHubAccountService } from "../../services/github/delete-github-account.service";
import { testGitHubConnectionService } from "../../services/github/test-github-connection.service";
import { getGitHubRepositoriesService } from "../../services/github/get-github-repositories.service";
import { startGitHubAppManifestRegistrationService } from "../../services/github/start-github-app-manifest-registration.service";

interface GitHubSettingsClientProps {
  initialAccounts: GitHubAccount[];
}

const GithubIcon = ({ size = 24, className = "" }: { size?: number, className?: string }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4" />
    <path d="M9 18c-4.51 2-5-2-7-2" />
  </svg>
);

export const GitHubSettingsClient: React.FC<GitHubSettingsClientProps> = ({ initialAccounts }) => {
  const [accounts, setAccounts] = useState<GitHubAccount[]>(initialAccounts);
  const [repoCounts, setRepoCounts] = useState<Record<string, number>>({});
  const [connectionStatuses, setConnectionStatuses] = useState<Record<string, "testing" | "success" | "error">>({});
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingAccount, setEditingAccount] = useState<GitHubAccount | undefined>();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Fetch accounts on mount in case initial is empty or we want to ensure freshness
  useEffect(() => {
    refreshAccounts();
  }, []);

  // Fetch repo counts whenever accounts change
  useEffect(() => {
    accounts.forEach(async (acc) => {
      if (!repoCounts[acc.id] && acc.id) {
        const { data } = await getGitHubRepositoriesService(acc.id);
        if (data) {
          setRepoCounts((prev) => ({ ...prev, [acc.id]: data.length }));
        }
      }
    });
  }, [accounts]);

  const refreshAccounts = async () => {
    const { data } = await getGitHubAccountsService();
    if (data) setAccounts(data);
  };

  const handleOpenCreate = () => {
    setEditingAccount(undefined);
    setError(null);
    setIsModalOpen(true);
  };

  const handleOpenEdit = (acc: GitHubAccount) => {
    setEditingAccount(acc);
    setError(null);
    setIsModalOpen(true);
  };

  const handleSubmit = async (formData: { 
    auth_method: "personal_access_token" | "github_app";
    account_type: "personal" | "organization"; 
    display_name: string; 
    account_identifier: string; 
    personal_access_token?: string; 
  }) => {
    setIsLoading(true);
    setError(null);

    if (formData.auth_method === "github_app") {
      const res = await startGitHubAppManifestRegistrationService({
        display_name: formData.display_name,
        owner_type: formData.account_type,
        organization: formData.account_type === "organization" ? formData.account_identifier : undefined,
      });

      if (res.error || !res.data) {
        setError(res.error?.message || "Failed to start GitHub App registration");
        setIsLoading(false);
        return;
      }

      // Dynamically submit the form to GitHub
      const form = document.createElement("form");
      form.method = "POST";
      form.action = res.data.form_action;
      
      const input = document.createElement("input");
      input.type = "hidden";
      input.name = "manifest";
      input.value = res.data.manifest;
      
      form.appendChild(input);
      document.body.appendChild(form);
      form.submit();
      
      // We don't set isLoading(false) here because the browser will navigate away
      return;
    }

    let res;
    if (editingAccount) {
      // Update
      const body: any = { display_name: formData.display_name };
      if (formData.personal_access_token) body.personal_access_token = formData.personal_access_token;
      res = await updateGitHubPatService(editingAccount.id!, body);
    } else {
      // Create PAT
      res = await createGitHubPatService({
        account_type: formData.account_type,
        display_name: formData.display_name,
        account_identifier: formData.account_identifier,
        personal_access_token: formData.personal_access_token!
      });
    }

    if (res.error) {
      setError(res.error.message || "An error occurred");
      toast.error(res.error.message || "An error occurred");
      setIsLoading(false);
      return;
    }

    await refreshAccounts();
    toast.success(`GitHub account ${editingAccount ? "updated" : "connected"} successfully`);
    setIsLoading(false);
    setIsModalOpen(false);
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure you want to remove this GitHub account?")) return;
    await deleteGitHubAccountService(id);
    await refreshAccounts();
    toast.success("GitHub account removed");
  };

  const handleTestConnection = async (id: string) => {
    setConnectionStatuses((prev) => ({ ...prev, [id]: "testing" }));
    const { success } = await testGitHubConnectionService(id);
    setConnectionStatuses((prev) => ({ ...prev, [id]: success ? "success" : "error" }));
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div>
          <h2 className={styles.title}>GitHub Accounts</h2>
          <p className={styles.subtitle}>
            Connect GitHub personal or organization accounts to allow Akritas to read code and create fixes.
          </p>
        </div>
        <Button variant="primary" onClick={handleOpenCreate} className={styles.addBtn}>
          <Plus size={16} /> Connect Account
        </Button>
      </div>

      <div className={styles.list}>
        {accounts.length === 0 ? (
          <div className={styles.emptyState}>
            <GithubIcon size={48} className={styles.emptyIcon} />
            <p>No GitHub accounts connected yet.</p>
          </div>
        ) : (
          accounts.map((acc) => (
            <div key={acc.id} className={styles.accountCard}>
              <div className={styles.accountInfo}>
                <div className={styles.avatar}>
                  <GithubIcon size={24} />
                </div>
                <div className={styles.details}>
                  <div className={styles.nameRow}>
                    <span className={styles.name}>{acc.display_name} ({acc.account_identifier})</span>
                    <Badge variant="neutral">{acc.account_type}</Badge>
                    {acc.authentication_method === "github_app" && (
                      <Badge variant="intel">GitHub App</Badge>
                    )}
                  </div>
                  <div className={styles.metaRow}>
                    <span className={styles.metaText}>
                      {repoCounts[acc.id!] !== undefined 
                        ? `${repoCounts[acc.id!]} repositories` 
                        : "Loading repositories..."}
                    </span>
                    <span className={styles.dot}>•</span>
                    <span className={styles.metaText}>
                      Added {new Date(acc.created_at!).toLocaleDateString()}
                    </span>
                  </div>
                </div>
              </div>

              <div className={styles.accountStatus}>
                {connectionStatuses[acc.id!] === "testing" && (
                  <span className={styles.statusTesting}><RefreshCw size={14} className={styles.spin} /> Testing...</span>
                )}
                {connectionStatuses[acc.id!] === "success" && (
                  <span className={styles.statusSuccess}><CheckCircle size={14} /> Connected</span>
                )}
                {connectionStatuses[acc.id!] === "error" && (
                  <span className={styles.statusError}><AlertCircle size={14} /> Error</span>
                )}
              </div>

              <div className={styles.accountActions}>
                <Button variant="secondary" size="sm" onClick={() => handleTestConnection(acc.id!)}>
                  Test
                </Button>
                {acc.authentication_method === "personal_access_token" && (
                  <Button variant="ghost" size="sm" onClick={() => handleOpenEdit(acc)}>
                    <Edit2 size={16} />
                  </Button>
                )}
                <Button variant="ghost" size="sm" onClick={() => handleDelete(acc.id!)} className={styles.deleteBtn}>
                  <Trash2 size={16} />
                </Button>
              </div>
            </div>
          ))
        )}
      </div>

      <Modal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)}
        title={editingAccount ? "Edit GitHub Account" : "Connect GitHub Account"}
      >
        {error && <div className={styles.errorAlert}>{error}</div>}
        <GitHubAccountForm 
          initialData={editingAccount}
          onSubmit={handleSubmit}
          onCancel={() => setIsModalOpen(false)}
          isLoading={isLoading}
        />
      </Modal>
    </div>
  );
};
