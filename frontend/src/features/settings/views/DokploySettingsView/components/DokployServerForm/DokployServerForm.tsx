import React, { useState } from "react";
import styles from "./DokployServerForm.module.css";
import { Button } from "@/core/ui/primitives/Button";
import { DokployServer } from "../../../../services/dokploy/list-dokploy-servers.service";

interface DokployServerFormProps {
  initialData?: DokployServer;
  onSubmit: (data: { name: string; base_url: string; api_credential?: string }) => Promise<void>;
  onCancel: () => void;
  isLoading: boolean;
}

export const DokployServerForm: React.FC<DokployServerFormProps> = ({
  initialData,
  onSubmit,
  onCancel,
  isLoading,
}) => {
  const [name, setName] = useState(initialData?.name || "");
  const [baseUrl, setBaseUrl] = useState(initialData?.base_url || "");
  const [apiCredential, setApiCredential] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmit({
      name,
      base_url: baseUrl,
      ...(apiCredential ? { api_credential: apiCredential } : {})
    });
  };

  const isEdit = !!initialData;

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      <div className={styles.formGroup}>
        <label className={styles.label} htmlFor="name">Display Name</label>
        <input
          id="name"
          className={styles.input}
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="e.g. Production Cluster"
          required
        />
        <p className={styles.helpText}>
          A friendly name to identify this Dokploy instance.
        </p>
      </div>

      <div className={styles.formGroup}>
        <label className={styles.label} htmlFor="base_url">Server URL</label>
        <input
          id="base_url"
          type="url"
          className={styles.input}
          value={baseUrl}
          onChange={(e) => setBaseUrl(e.target.value)}
          placeholder="https://dokploy.mycompany.com"
          required
        />
        <p className={styles.helpText}>
          The public HTTPS URL where your Dokploy instance is reachable.
        </p>
      </div>

      <div className={styles.formGroup}>
        <label className={styles.label} htmlFor="api_credential">API Credential</label>
        <input
          id="api_credential"
          type="password"
          className={styles.input}
          value={apiCredential}
          onChange={(e) => setApiCredential(e.target.value)}
          placeholder={isEdit ? "Leave blank to keep current credential" : "Enter Dokploy API token"}
          required={!isEdit}
        />
        <p className={styles.helpText}>
          Generated from your Dokploy dashboard.
        </p>
      </div>

      <div className={styles.actions}>
        <Button variant="ghost" type="button" onClick={onCancel} disabled={isLoading}>
          Cancel
        </Button>
        <Button variant="primary" type="submit" isLoading={isLoading}>
          {isEdit ? "Update Server" : "Connect Server"}
        </Button>
      </div>
    </form>
  );
};
