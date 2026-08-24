"use client";

import { Suspense, useEffect } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { RefreshCw, AlertCircle } from "lucide-react";
import { env } from "@/core/config/env.config";
import styles from "./CallbackPage.module.css";

function buildApiUrl(path: string): URL {
  return new URL(`${env.apiUrl.replace(/\/$/, "")}${path}`, window.location.origin);
}

function GitHubCallbackContent() {
  const searchParams = useSearchParams();

  const errorParam =
    searchParams.get("error_description") ?? searchParams.get("error");

  const code = searchParams.get("code");
  const state = searchParams.get("state");
  const installationId = searchParams.get("installation_id");
  const setupAction = searchParams.get("setup_action");

  const isManifest = Boolean(code && state);
  const isInstallation = Boolean(installationId);
  const hasValidCallback = isManifest || isInstallation;

  const error = errorParam
    ? errorParam
    : !hasValidCallback
      ? "No valid callback parameters found in URL."
      : null;

  useEffect(() => {
    if (error) return;

    if (isManifest && code && state) {
      const url = buildApiUrl("/integrations/github/app-manifest/callback");

      url.searchParams.set("code", code);
      url.searchParams.set("state", state);

      window.location.assign(url.toString());
      return;
    }

    if (isInstallation && installationId) {
      const url = buildApiUrl("/integrations/github/app-installations/callback");

      url.searchParams.set("installation_id", installationId);

      if (setupAction) {
        url.searchParams.set("setup_action", setupAction);
      }

      if (state) {
        url.searchParams.set("state", state);
      }

      window.location.assign(url.toString());
    }
  }, [
    error,
    isManifest,
    isInstallation,
    code,
    state,
    installationId,
    setupAction,
  ]);

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
      <p>
        Please wait while we complete the setup. You will be redirected shortly.
      </p>
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
