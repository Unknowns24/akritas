"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/core/ui/primitives/Button";
import { AlertTriangle, Copy } from "lucide-react";
import { APP_ROUTES } from "@/core/routes/routes.config";
import { getErrorMessage } from "@/core/errors";
import { OtpInput } from "@/core/ui/primitives/OtpInput";
import { QRCodeSVG } from "qrcode.react";
import {
  startAdministratorRecoveryService,
  verifyAdministratorRecoveryService,
  RecoveryRequest,
  TotpEnrollment,
} from "../../services/auth.service";
import { ApiError } from "@/core/errors/api-error";
import styles from "../SetupView/SetupView.module.css"; // Reuse setup view styling since they are extremely similar

export const RecoveryView = () => {
  const router = useRouter();

  // Step 1 State
  const [step, setStep] = useState<1 | 2>(1);
  const [formData, setFormData] = useState<RecoveryRequest>({
    email: "",
    new_password: "",
    bootstrap_token: "",
  });
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Step 2 State
  const [enrollment, setEnrollment] = useState<TotpEnrollment | null>(null);
  const [totpCode, setTotpCode] = useState("");

  const handleRecoverySubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (formData.new_password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    if (formData.new_password.length < 12) {
      setError("Password must be at least 12 characters");
      return;
    }

    setIsSubmitting(true);
    try {
      const result = await startAdministratorRecoveryService(formData);
      setEnrollment(result);
      setStep(2);
    } catch (error: unknown) {
      if (error instanceof ApiError) {
        if (error.status === 401) {
          setError("Invalid bootstrap token or administrator email.");
        } else if (error.status === 429) {
          setError("Too many requests. Please try again later.");
        } else {
          setError(getErrorMessage(error, "Failed to start recovery process."));
        }
      } else {
        setError(getErrorMessage(error, "Failed to start recovery process."));
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleVerifySubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!enrollment || totpCode.length !== 6) return;
    
    setError(null);
    setIsSubmitting(true);
    try {
      await verifyAdministratorRecoveryService(totpCode, enrollment.enrollment_id);
      router.replace(APP_ROUTES.AUTH.LOGIN);
    } catch (error: unknown) {
      if (error instanceof ApiError) {
        if (error.status === 401) {
          setError("Invalid authenticator code or enrollment expired.");
        } else if (error.status === 429) {
          setError("Too many requests. Please try again later.");
        } else {
          setError(getErrorMessage(error, "Invalid authenticator code."));
        }
      } else {
        setError(getErrorMessage(error, "Invalid authenticator code."));
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const copyToClipboard = () => {
    if (enrollment?.manual_entry_key) {
      navigator.clipboard.writeText(enrollment.manual_entry_key);
    }
  };

  if (step === 2 && enrollment) {
    return (
      <div className={styles.card}>
        <h1 className={styles.title} style={{ textAlign: "center" }}>Configure a new authenticator.</h1>
        <p className={styles.subtitle} style={{ textAlign: "center" }}>
          Scan the new TOTP enrollment and verify a code to complete account recovery.
        </p>

        <div className={styles.qrCol} style={{ marginBottom: "24px", paddingTop: "0" }}>
          <div className={styles.qrContainer}>
            <QRCodeSVG value={enrollment.otpauth_uri} size={180} level="M" />
          </div>
        </div>
          
        <div className={styles.manualCol}>
          <div className={styles.manualGroup}>
            <label>MANUAL ENTRY KEY</label>
            <div className={styles.manualInputWrapper}>
              <input type="text" readOnly value={enrollment.manual_entry_key} />
              <button type="button" onClick={copyToClipboard}><Copy size={16} /></button>
            </div>
          </div>

          <form onSubmit={handleVerifySubmit} className={styles.totpForm}>
            <div style={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
              <label style={{ textAlign: "center" }}>VERIFY 6-DIGIT CODE</label>
              <OtpInput value={totpCode} onChange={setTotpCode} disabled={isSubmitting} />
            </div>
            
            <div className={styles.warningBox} style={{ marginTop: "24px" }}>
              <AlertTriangle size={16} className={styles.warningIcon} />
              <p>Completing recovery revokes all previous Akritas sessions.</p>
            </div>

            {error && <div className={styles.errorBox}>{error}</div>}
            
            <Button type="submit" variant="secondary" size="lg" disabled={totpCode.length !== 6 || isSubmitting} className={styles.fullWidthBtn} style={{ marginTop: "24px" }}>
              Complete recovery
            </Button>
          </form>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.card}>
      <h1 className={styles.title}>Recover administrator access</h1>
      <p className={styles.subtitle}>
        Reset the administrator password and authenticator using the deployment bootstrap token.
      </p>

      <div className={styles.warningBox} style={{ backgroundColor: "rgba(var(--status-error-rgb), 0.2)", borderColor: "var(--status-error)" }}>
        <AlertTriangle size={24} className={styles.warningIcon} />
        <p>Completing recovery will replace the current authenticator enrollment and revoke previous sessions.</p>
      </div>

      <form onSubmit={handleRecoverySubmit} className={styles.form}>
        <div className={styles.inputGroup}>
          <label>Administrator email</label>
          <input
            type="email"
            required
            placeholder="admin@akritas.internal"
            value={formData.email}
            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
            disabled={isSubmitting}
          />
        </div>

        <div className={styles.inputGroup}>
          <label>New password</label>
          <input
            type="password"
            required
            placeholder="••••••••••••"
            value={formData.new_password}
            onChange={(e) => setFormData({ ...formData, new_password: e.target.value })}
            disabled={isSubmitting}
          />
        </div>

        <div className={styles.inputGroup}>
          <label>Confirm new password</label>
          <input
            type="password"
            required
            placeholder="••••••••••••"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            disabled={isSubmitting}
          />
        </div>

        <div className={styles.inputGroup} style={{ marginTop: "12px" }}>
          <label>Bootstrap token</label>
          <input
            type="password"
            required
            placeholder="akr_boot_..."
            value={formData.bootstrap_token}
            onChange={(e) => setFormData({ ...formData, bootstrap_token: e.target.value })}
            disabled={isSubmitting}
          />
        </div>

        {error && <div className={styles.errorBox}><AlertTriangle size={16}/> {error}</div>}

        <Button type="submit" variant="secondary" size="lg" className={styles.fullWidthBtn} disabled={isSubmitting} style={{ marginTop: "12px" }}>
          Continue to authenticator setup →
        </Button>
        
        <div className={styles.footerLinks}>
          <button type="button" onClick={() => router.push(APP_ROUTES.AUTH.LOGIN)}>
            Return to sign in
          </button>
        </div>
      </form>
    </div>
  );
};
