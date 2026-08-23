"use client";

import { useEffect, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { env } from "@/core/config/env.config";
import { RefreshCw, AlertCircle } from "lucide-react";
import styles from "./CallbackPage.module.css";

function GitHubCallbackContent() {
  const searchParams = useSearchParams();

  // Compute params and derived error synchronously during render
  const errorParam = searchParams.get("error_description") || searchParams.get("error");
  const code = searchParams.get("code");
  const state = searchParams.get("state");
  const installationId = searchParams.get("installation_id");
  const setupAction = searchParams.get("setup_action");

  const isManifest = Boolean(code && state);
  const isInstallation = Boolean(installationId);

  const error = errorParam
    ? errorParam
    : !isManifest && !isInstallation
    ? "No valid callback parameters found in URL."
    : null;

  useEffect(() => {
    // Only perform the side effect (external redirect) if parameters are valid and error-free
    if (error) return;

    if (isManifest && code && state) {
      // 1. Return from App Manifest Creation
      const url = new URL(`${env.apiUrl}/integrations/github/app-manifest/callback`, window.location.origin);
      url.searchParams.set("code", code);
      url.searchParams.set("state", state);
      window.location.href = url.toString();
    } else if (isInstallation && installationId) {
      // 2. Return from App Installation
      const url = new URL(`${env.apiUrl}/integrations/github/app-installations/callback`, window.location.origin);
      url.searchParams.set("installation_id", installationId);
      if (setupAction) url.searchParams.set("setup_action", setupAction);
      if (state) url.searchParams.set("state", state);
      window.location.href = url.toString();
    }
  }, [error, isManifest, isInstallation, code, state, installationId, setupAction]);

  if (error) {
    return (
      <div className={styles.container}>
        <AlertCircle size={48} className={styles.error} />
        <h2>Error in GitHub Integration</h2>
        <p className={styles.error}>{error}</p>
        <Link href="/settings/github" className={styles.link}>
          Return to GitHub Settings
        </Link>
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

export default function GitHubCallbackPage() {
  return (
    <Suspense
      fallback={
        <div className={styles.container}>
          <RefreshCw size={48} className={styles.spin} />
          <h2>Loading...</h2>
        </div>
      }
    >
      <GitHubCallbackContent />
    </Suspense>
  );
}
