import React, { useState } from "react";
import styles from "./GitHubAccountForm.module.css";
import { Button } from "@/core/ui/primitives/Button";
import { GitHubAccount } from "../../../../services/github/get-github-accounts.service";

interface GitHubAccountFormProps {
  initialData?: GitHubAccount;
  onSubmit: (data: { 
    auth_method: "personal_access_token" | "github_app";
    account_type: "personal" | "organization"; 
    display_name: string; 
    account_identifier: string; 
    personal_access_token?: string; 
  }) => Promise<void>;
  onCancel: () => void;
  isLoading: boolean;
}

export const GitHubAccountForm: React.FC<GitHubAccountFormProps> = ({
  initialData,
  onSubmit,
  onCancel,
  isLoading,
}) => {
  const [authMethod, setAuthMethod] = useState<"personal_access_token" | "github_app">(
    initialData?.authentication_method === "github_app" ? "github_app" : "personal_access_token"
  );
  const [accountType, setAccountType] = useState<"personal" | "organization">(
    initialData?.account_type || "personal"
  );
  const [displayName, setDisplayName] = useState(initialData?.display_name || "");
  const [accountIdentifier, setAccountIdentifier] = useState(initialData?.account_identifier || "");
  const [pat, setPat] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmit({ 
      auth_method: authMethod,
      account_type: accountType, 
      display_name: displayName, 
      account_identifier: accountIdentifier, 
      ...(authMethod === "personal_access_token" ? { personal_access_token: pat } : {})
    });
  };

  const isEdit = !!initialData;

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      {!isEdit && (
        <div className={styles.formGroup}>
          <label className={styles.label}>Authentication Method</label>
          <div className={styles.radioGroup}>
            <label className={styles.radioLabel}>
              <input
                type="radio"
                name="auth_method"
                value="github_app"
                checked={authMethod === "github_app"}
                onChange={() => setAuthMethod("github_app")}
              />
              GitHub App (Recommended)
            </label>
            <label className={styles.radioLabel}>
              <input
                type="radio"
                name="auth_method"
                value="personal_access_token"
                checked={authMethod === "personal_access_token"}
                onChange={() => setAuthMethod("personal_access_token")}
              />
              Personal Access Token
            </label>
          </div>
        </div>
      )}

      <div className={styles.formGroup}>
        <label className={styles.label}>Account Type</label>
        <div className={styles.radioGroup}>
          <label className={styles.radioLabel}>
            <input
              type="radio"
              name="account_type"
              value="personal"
              checked={accountType === "personal"}
              onChange={() => setAccountType("personal")}
              disabled={isEdit}
            />
            Personal
          </label>
          <label className={styles.radioLabel}>
            <input
              type="radio"
              name="account_type"
              value="organization"
              checked={accountType === "organization"}
              onChange={() => setAccountType("organization")}
              disabled={isEdit}
            />
            Organization
          </label>
        </div>
      </div>

      <div className={styles.formGroup}>
        <label className={styles.label} htmlFor="display_name">Display Name</label>
        <input
          id="display_name"
          className={styles.input}
          value={displayName}
          onChange={(e) => setDisplayName(e.target.value)}
          placeholder="e.g. My Personal Account"
          required
        />
        <p className={styles.helpText}>
          A friendly name to identify this connection.
        </p>
      </div>

      <div className={styles.formGroup}>
        <label className={styles.label} htmlFor="account_identifier">Account Identifier</label>
        <input
          id="account_identifier"
          className={styles.input}
          value={accountIdentifier}
          onChange={(e) => setAccountIdentifier(e.target.value)}
          placeholder="e.g. your-username or org-name"
          required
          disabled={isEdit}
        />
        <p className={styles.helpText}>
          The exact GitHub username or organization name.
        </p>
      </div>

      {authMethod === "personal_access_token" && (
        <div className={styles.formGroup}>
          <label className={styles.label} htmlFor="pat">Personal Access Token (PAT)</label>
          <input
            id="pat"
            type="password"
            className={styles.input}
            value={pat}
            onChange={(e) => setPat(e.target.value)}
            placeholder={isEdit ? "Leave blank to keep current PAT" : "ghp_..."}
            required={!isEdit}
          />
          <p className={styles.helpText}>
            Requires &apos;repo&apos; scope. Classic tokens recommended.
          </p>
        </div>
      )}

      <div className={styles.actions}>
        <Button variant="ghost" type="button" onClick={onCancel} disabled={isLoading}>
          Cancel
        </Button>
        <Button variant="primary" type="submit" isLoading={isLoading}>
          {isEdit ? "Update Account" : "Connect Account"}
        </Button>
      </div>
    </form>
  );
};
