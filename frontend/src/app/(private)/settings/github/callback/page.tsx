"use client";

import { useEffect, useMemo } from "react";
import { useSearchParams } from "next/navigation";
import { env } from "@/core/config/env.config";
import { RefreshCw } from "lucide-react";
import styles from "./CallbackPage.module.css";

export default function GitHubCallbackPage() {
  const searchParams = useSearchParams();
  const hasValidCallback = useMemo(() => {
    return Boolean(
      (searchParams.get("code") && searchParams.get("state")) ||
        searchParams.get("installation_id"),
    );
  }, [searchParams]);

  useEffect(() => {
    // We are on the client, let's process the URL
    const code = searchParams.get("code");
    const state = searchParams.get("state");
    const installation_id = searchParams.get("installation_id");
    const setup_action = searchParams.get("setup_action");

    if (code && state) {
      // 1. Return from App Manifest Creation
      // Redirect to backend callback to exchange code and redirect to GitHub installation
      const url = new URL(`${env.apiUrl}/integrations/github/app-manifest/callback`);
      url.searchParams.set("code", code);
      url.searchParams.set("state", state);
      window.location.href = url.toString();
    } else if (installation_id) {
      // 2. Return from App Installation
      // Redirect to backend callback to save installation and redirect back to Settings
      const url = new URL(`${env.apiUrl}/integrations/github/app-installations/callback`);
      url.searchParams.set("installation_id", installation_id);
      if (setup_action) url.searchParams.set("setup_action", setup_action);
      if (state) url.searchParams.set("state", state);
      window.location.href = url.toString();
    }
  }, [searchParams]);

  if (!hasValidCallback) {
    return (
      <div className={styles.container}>
        <h2>Error in GitHub Integration</h2>
        <p className={styles.error}>No valid callback parameters found in URL.</p>
        <a href="/settings/github" className={styles.link}>Return to GitHub Settings</a>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <RefreshCw size={48} className={styles.spin} />
      <h2>Connecting to GitHub...</h2>
      <p>Please wait while we complete the setup. You will be redirected shortly.</p>
    </div>
  );
}
