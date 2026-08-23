"use client";

import React, { useEffect, useState, createContext, useContext } from "react";
import { useRouter, usePathname } from "next/navigation";
import { getAuthSetupStatusService, getCurrentSessionService, SessionResponse } from "../../services/auth.service";
import styles from "./AuthProvider.module.css";
import { APP_ROUTES } from "@/core/routes/routes.config";

interface AuthContextValue {
  session: SessionResponse["data"] | null;
}

const AuthContext = createContext<AuthContextValue>({ session: null });

export const useAuth = () => useContext(AuthContext);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const router = useRouter();
  const pathname = usePathname();
  const [loading, setLoading] = useState(true);
  const [authorized, setAuthorized] = useState(false);
  const [session, setSession] = useState<SessionResponse["data"] | null>(null);

  useEffect(() => {
    let mounted = true;
    
    const checkAuth = async () => {
      try {
        // 1. Check if setup is needed
        try {
          const setupStatus = await getAuthSetupStatusService();
          if (setupStatus.setup_required) {
            if (pathname !== APP_ROUTES.AUTH.SETUP) {
              if (mounted) router.replace(APP_ROUTES.AUTH.SETUP);
              return;
            }
            if (mounted) {
              setAuthorized(true);
              setLoading(false);
            }
            return;
          } else if (pathname === APP_ROUTES.AUTH.SETUP) {
            if (mounted) {
              setAuthorized(true);
              setLoading(false);
            }
            return;
          }
        } catch (e) {
          console.warn("Failed to check setup status", e);
        }

        // 2. Check session
        const isAuthRoute = Object.values(APP_ROUTES.AUTH).includes(pathname as any);
        try {
          const fetchedSession = await getCurrentSessionService();
          if (fetchedSession && fetchedSession.administrator) {
            if (mounted) setSession(fetchedSession);
            if (isAuthRoute) {
              if (mounted) router.replace(APP_ROUTES.OVERVIEW);
            } else {
              if (mounted) setAuthorized(true);
            }
          } else {
            if (!isAuthRoute) {
              if (mounted) router.replace(APP_ROUTES.AUTH.LOGIN);
            } else {
              if (mounted) setAuthorized(true);
            }
          }
        } catch (error) {
          if (!isAuthRoute) {
            if (mounted) router.replace(APP_ROUTES.AUTH.LOGIN);
          } else {
            if (mounted) setAuthorized(true);
          }
        }
      } finally {
        if (mounted) setLoading(false);
      }
    };

    checkAuth();
    
    return () => {
      mounted = false;
    };
  }, [pathname, router]);

  if (loading) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingCard}>
          <div style={{ display: "flex", flexDirection: "column", gap: "24px", alignItems: "center" }}>
            <div className={styles.logo}>AKRITAS</div>
            <div style={{ display: "flex", alignItems: "center", gap: "12px", justifyContent: "center" }}>
              <div className={styles.spinner}></div>
              <span className={styles.text}>Checking installation...</span>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <AuthContext.Provider value={{ session }}>
      {authorized ? children : null}
    </AuthContext.Provider>
  );
};
