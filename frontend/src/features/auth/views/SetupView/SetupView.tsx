"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/core/ui/primitives/Button";
import { Info, AlertTriangle, ShieldCheck, Copy } from "lucide-react";
import { APP_ROUTES } from "@/core/routes/routes.config";
import { OtpInput } from "@/core/ui/primitives/OtpInput";
import { QRCodeSVG } from "qrcode.react";
import {
  getAuthSetupStatusService,
  startAdministratorSetupService,
  verifyAdministratorSetupService,
  SetupRequest,
  TotpEnrollment,
} from "../../services/auth.service";
import styles from "./SetupView.module.css";

export const SetupView = () => {
  const router = useRouter();
  
  const [loading, setLoading] = useState(true);
  const [alreadyInitialized, setAlreadyInitialized] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Step 1 State
  const [step, setStep] = useState<1 | 2>(1);
  const [formData, setFormData] = useState<SetupRequest>({
    display_name: "",
    email: "",
    password: "",
    bootstrap_token: "",
  });
  const [confirmPassword, setConfirmPassword] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Step 2 State
  const [enrollment, setEnrollment] = useState<TotpEnrollment | null>(null);
  const [totpCode, setTotpCode] = useState("");

  useEffect(() => {
    // Check if we are actually allowed to be here
    getAuthSetupStatusService()
      .then((status) => {
        if (!status.setup_required) {
          setAlreadyInitialized(true);
        }
      })
      .catch(() => {
        // If it fails, assume it's initialized for safety
        setAlreadyInitialized(true);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  const handleSetupSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (formData.password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    if (formData.password.length < 12) {
      setError("Password must be at least 12 characters");
      return;
    }

    setIsSubmitting(true);
    try {
      const result = await startAdministratorSetupService(formData);
      setEnrollment(result);
      setStep(2);
    } catch (err: any) {
      setError(err.message || "Failed to initialize Akritas");
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
      await verifyAdministratorSetupService(totpCode, enrollment.enrollment_id);
      // Success! Redirect to the authenticated app
      router.replace(APP_ROUTES.OVERVIEW);
    } catch (err: any) {
      setError(err.message || "Invalid authenticator code");
    } finally {
      setIsSubmitting(false);
    }
  };

  const copyToClipboard = () => {
    if (enrollment?.manual_entry_key) {
      navigator.clipboard.writeText(enrollment.manual_entry_key);
    }
  };

  if (loading) {
    return <div className={styles.loadingSpinner} />;
  }

  if (alreadyInitialized) {
    return (
      <div className={styles.card}>
        <div className={styles.initializedIcon}>
          <ShieldCheck size={32} />
        </div>
        <h1 className={styles.title} style={{ textAlign: "center" }}>Akritas is already initialized</h1>
        <p className={styles.subtitle} style={{ textAlign: "center", marginBottom: "32px" }}>
          Administrator setup is no longer available for this installation.
        </p>
        <Button variant="secondary" size="lg" className={styles.fullWidthBtn} onClick={() => router.push(APP_ROUTES.AUTH.LOGIN)}>
          Go to sign in
        </Button>
      </div>
    );
  }

  if (step === 2 && enrollment) {
    return (
      <div className={styles.cardWide}>
        <div className={styles.header}>
          <div className={styles.shieldIcon}><ShieldCheck size={16} /></div>
          <h1 className={styles.title}>Secure your administrator account.</h1>
        </div>
        <p className={styles.subtitle}>
          Add Akritas to your authenticator app to complete setup.
        </p>
        
        <div className={styles.stepsBar}>
          <span className={styles.stepPast}>1 Account</span>
          <span className={styles.stepSep}>›</span>
          <span className={styles.stepCurrent}><span className={styles.circleActive}>2</span> Authenticator</span>
          <span className={styles.stepSep}>›</span>
          <span className={styles.stepFuture}>3 Complete</span>
        </div>

        <div className={styles.step2Grid}>
          <div className={styles.qrCol}>
            <div className={styles.qrContainer}>
              <QRCodeSVG value={enrollment.otpauth_uri} size={180} level="M" />
            </div>
            <h3 className={styles.qrTitle}>Scan with your<br/>authenticator app</h3>
            <p className={styles.qrSub}>
              Works with Google Authenticator, Authy, 1Password, or similar TOTP apps.
            </p>
          </div>
          
          <div className={styles.manualCol}>
            <div className={styles.manualGroup}>
              <label>MANUAL CONFIGURATION</label>
              <div className={styles.manualInputWrapper}>
                <input type="text" readOnly value={enrollment.manual_entry_key} />
                <button type="button" onClick={copyToClipboard}><Copy size={16} /></button>
              </div>
            </div>

            <div className={styles.warningBox}>
              <AlertTriangle size={16} className={styles.warningIcon} />
              <p>This enrollment information is shown only during setup. Do not store it in browser storage.</p>
            </div>

            <form onSubmit={handleVerifySubmit} className={styles.totpForm}>
              <label>Enter 6-digit code</label>
              <OtpInput value={totpCode} onChange={setTotpCode} disabled={isSubmitting} />
              {error && <div className={styles.errorText}>{error}</div>}
              
              <Button type="submit" variant="secondary" size="lg" disabled={totpCode.length !== 6 || isSubmitting} className={styles.fullWidthBtn} style={{ marginTop: "24px" }}>
                Complete setup →
              </Button>
            </form>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.card}>
      <h1 className={styles.title}>Initialize Akritas</h1>
      <p className={styles.subtitle}>
        Create the administrator account for this Akritas installation.
      </p>

      <div className={styles.infoBox}>
        <Info size={18} className={styles.infoIcon} />
        <p>This setup is available only until the first administrator is activated.</p>
      </div>

      <form onSubmit={handleSetupSubmit} className={styles.form}>
        <div className={styles.inputGroup}>
          <label>Display Name</label>
          <input
            type="text"
            required
            placeholder="e.g. System Admin"
            value={formData.display_name}
            onChange={(e) => setFormData({ ...formData, display_name: e.target.value })}
            disabled={isSubmitting}
          />
        </div>

        <div className={styles.inputGroup}>
          <label>Email address</label>
          <input
            type="email"
            required
            placeholder="admin@example.com"
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
          <span className={styles.hint}>Minimum 12 characters required.</span>
        </div>

        <div className={styles.inputGroup}>
          <label>Confirm password</label>
          <input
            type="password"
            required
            placeholder="••••••••••••"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            disabled={isSubmitting}
          />
        </div>

        <div className={styles.divider} />

        <div className={styles.sectionTitle}>Deployment Authorization</div>
        <p className={styles.sectionDesc}>
          The bootstrap token is configured when Akritas is deployed and is used only to authorize initial setup.
        </p>

        <div className={styles.inputGroup}>
          <label>Bootstrap token</label>
          <input
            type="password"
            required
            placeholder="ak_bt_XXXXXXXXXXXXXXXXX"
            value={formData.bootstrap_token}
            onChange={(e) => setFormData({ ...formData, bootstrap_token: e.target.value })}
            disabled={isSubmitting}
          />
        </div>

        {error && <div className={styles.errorBox}><AlertTriangle size={16}/> {error}</div>}

        <Button type="submit" variant="secondary" size="lg" className={styles.fullWidthBtn} disabled={isSubmitting}>
          Continue to authenticator setup →
        </Button>
        
        <div className={styles.footerLinks}>
          <span>Already initialized?</span>
          <button type="button" onClick={() => router.push(APP_ROUTES.AUTH.LOGIN)}>Sign in</button>
        </div>
      </form>
    </div>
  );
};
