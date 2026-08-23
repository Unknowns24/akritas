"use client";

import React, { useEffect, useState, useCallback } from "react";
import styles from "./QvacSettingsView.module.css";
import { Button } from "@/core/ui/primitives/Button";
import { toast } from "sonner";
import { getErrorMessage } from "@/core/errors";
import { getQvacConfigurationService } from "../../services/qvac/get-qvac-configuration.service";
import { putQvacConfigurationService, PutQvacConfigurationRequest } from "../../services/qvac/put-qvac-configuration.service";
import { Badge } from "@/core/ui/primitives/Badge";

type AuthType = "none" | "bearer" | "basic";

export const QvacSettingsView: React.FC = () => {
  const [endpointUrl, setEndpointUrl] = useState("");
  const [connectionTimeout, setConnectionTimeout] = useState(30);
  const [authType, setAuthType] = useState<AuthType>("none");
  
  const [bearerToken, setBearerToken] = useState("");
  const [basicUsername, setBasicUsername] = useState("");
  const [basicPassword, setBasicPassword] = useState("");
  
  const [credentialConfigured, setCredentialConfigured] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadConfig = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const config = await getQvacConfigurationService();
      setEndpointUrl(config.endpoint_url || "");
      setConnectionTimeout(config.connection_timeout_seconds || 30);
      setAuthType(config.authentication_type || "none");
      setCredentialConfigured(config.credential_configured || false);
      setBearerToken("");
      setBasicUsername("");
      setBasicPassword("");
    } catch (err: unknown) {
      const message = getErrorMessage(err, "Failed to load QVAC configuration");
      setError(message);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    loadConfig();
  }, [loadConfig]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    setError(null);

    try {
      const request: PutQvacConfigurationRequest = {
        endpoint_url: endpointUrl,
        connection_timeout_seconds: connectionTimeout,
        authentication: {
          type: authType,
        },
      };

      if (authType === "bearer" && bearerToken) {
        request.authentication.bearer_token = bearerToken;
      } else if (authType === "basic") {
        if (basicUsername) request.authentication.basic_username = basicUsername;
        if (basicPassword) request.authentication.basic_password = basicPassword;
      }

      const config = await putQvacConfigurationService(request);
      
      setEndpointUrl(config.endpoint_url);
      setConnectionTimeout(config.connection_timeout_seconds);
      setAuthType(config.authentication_type);
      setCredentialConfigured(config.credential_configured);
      
      setBearerToken("");
      setBasicUsername("");
      setBasicPassword("");
      
      toast.success("QVAC configuration updated successfully");
    } catch (err: unknown) {
      const message = getErrorMessage(err, "Failed to update QVAC configuration");
      setError(message);
      toast.error("Failed to save changes");
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return <div className={styles.container}>Loading QVAC configuration...</div>;
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h2 className={styles.title}>QVAC Inference Engine</h2>
        <p className={styles.subtitle}>
          Configure your local QVAC instance. QVAC processes sensitive data like logs and stack traces without sending them to external APIs.
        </p>
      </div>

      <div className={styles.card}>
        {error && <div className={styles.errorAlert}>{error}</div>}

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.formGroup}>
            <label htmlFor="endpoint" className={styles.label}>Endpoint URL</label>
            <input
              id="endpoint"
              type="url"
              value={endpointUrl}
              onChange={(e) => setEndpointUrl(e.target.value)}
              className={styles.input}
              placeholder="http://localhost:8081"
              required
            />
            <span className={styles.hint}>Must be a loopback or private network address.</span>
          </div>

          <div className={styles.formGroup}>
            <label htmlFor="authType" className={styles.label}>Authentication Type</label>
            <select
              id="authType"
              value={authType}
              onChange={(e) => setAuthType(e.target.value as AuthType)}
              className={styles.select}
            >
              <option value="none">None</option>
              <option value="bearer">Bearer Token</option>
              <option value="basic">Basic Auth</option>
            </select>
          </div>

          {authType === "bearer" && (
            <div className={styles.formGroup}>
              <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                <label htmlFor="bearerToken" className={styles.label}>Bearer Token</label>
                {credentialConfigured && (
                  <Badge variant="intel">Configured</Badge>
                )}
              </div>
              <input
                id="bearerToken"
                type="password"
                value={bearerToken}
                onChange={(e) => setBearerToken(e.target.value)}
                className={styles.input}
                placeholder={credentialConfigured ? "Enter to replace existing token" : "Enter bearer token"}
                required={!credentialConfigured}
              />
              <span className={styles.hint}>
                Enter the token value without the 'Bearer ' prefix.
              </span>
            </div>
          )}

          {authType === "basic" && (
            <>
              <div className={styles.formGroup}>
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                  <label htmlFor="basicUsername" className={styles.label}>Username</label>
                  {credentialConfigured && (
                    <Badge variant="intel">Configured</Badge>
                  )}
                </div>
                <input
                  id="basicUsername"
                  type="text"
                  value={basicUsername}
                  onChange={(e) => setBasicUsername(e.target.value)}
                  className={styles.input}
                  placeholder={credentialConfigured ? "Enter to replace existing credentials" : "Enter username"}
                  required={!credentialConfigured}
                />
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="basicPassword" className={styles.label}>Password</label>
                <input
                  id="basicPassword"
                  type="password"
                  value={basicPassword}
                  onChange={(e) => setBasicPassword(e.target.value)}
                  className={styles.input}
                  placeholder={credentialConfigured ? "Enter to replace existing credentials" : "Enter password"}
                  required={!credentialConfigured}
                />
              </div>
            </>
          )}

          <div className={styles.actions}>
            <Button type="submit" variant="primary" isLoading={isSaving} disabled={isSaving}>
              Save Configuration
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};
