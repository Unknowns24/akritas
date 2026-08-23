"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/core/ui/primitives/Button";
import { AlertCircle } from "lucide-react";
import { APP_ROUTES } from "@/core/routes/routes.config";
import { getErrorMessage } from "@/core/errors";
import { loginAdministratorService, LoginRequest } from "../../services/auth.service";
import styles from "./LoginView.module.css";

export const LoginView = () => {
  const router = useRouter();
  const [formData, setFormData] = useState<LoginRequest>({
    email: "",
    password: "",
    totp_code: "",
  });
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      await loginAdministratorService(formData);
      router.replace(APP_ROUTES.OVERVIEW);
    } catch (error: unknown) {
      setError(getErrorMessage(error, "Email, password, or authenticator code is incorrect"));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className={styles.card}>
      <h1 className={styles.title}>Sign in to Akritas</h1>
      <p className={styles.subtitle}>
        Authenticate to access production monitoring and incident response.
      </p>

      {error && (
        <div className={styles.errorBox}>
          <AlertCircle size={16} className={styles.errorIcon} />
          <p>{error}</p>
        </div>
      )}

      <form onSubmit={handleSubmit} className={styles.form}>
        <div className={styles.inputGroup}>
          <label>Email address</label>
          <input
            type="email"
            required
            placeholder="operator@akritas.local"
            value={formData.email}
            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
            disabled={isSubmitting}
          />
        </div>

        <div className={styles.inputGroup}>
          <label>Password</label>
          <input
            type="password"
            required
            placeholder="••••••••••••"
            value={formData.password}
            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
            disabled={isSubmitting}
          />
        </div>

        <div className={styles.inputGroup}>
          <div className={styles.labelRow}>
            <label>Authenticator code</label>
            <span className={styles.hint}>6-DIGIT TOTP</span>
          </div>
          <input
            type="text"
            required
            placeholder="000000"
            inputMode="numeric"
            maxLength={6}
            value={formData.totp_code}
            onChange={(e) => setFormData({ ...formData, totp_code: e.target.value.replace(/[^0-9]/g, "") })}
            disabled={isSubmitting}
            className={styles.totpInput}
          />
        </div>

        <Button type="submit" variant="secondary" size="lg" className={styles.fullWidthBtn} disabled={isSubmitting}>
          Sign in →
        </Button>

        <div className={styles.footerLinks}>
          <button type="button" onClick={() => router.push(APP_ROUTES.AUTH.RECOVERY)}>
            Lost authenticator?
          </button>
        </div>
      </form>
    </div>
  );
};
