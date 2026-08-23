"use client";

import React, { useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { Menu, Sparkles, LogOut, Loader2 } from "lucide-react";
import { Breadcrumbs, BreadcrumbItem } from "@/core/ui/layout/Breadcrumbs";
import { APP_ROUTES } from "@/core/routes/routes.config";
import { logoutAdministratorService } from "@/features/auth/services/auth.service";
import { useAuth } from "@/features/auth/components/AuthProvider/AuthProvider";
import styles from "./Header.module.css";

export interface HeaderProps {
  onToggleMobileMenu?: () => void;
  breadcrumbs?: BreadcrumbItem[];
}

export const Header: React.FC<HeaderProps> = ({
  onToggleMobileMenu,
  breadcrumbs,
}) => {
  const pathname = usePathname();
  const router = useRouter();
  const { session } = useAuth();
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  const getAutoBreadcrumbs = (): BreadcrumbItem[] => {
    if (breadcrumbs && breadcrumbs.length > 0) {
      return breadcrumbs;
    }

    if (pathname === APP_ROUTES.OVERVIEW) {
      return [{ label: "Dashboard", href: APP_ROUTES.OVERVIEW }];
    }

    const segments = pathname.split("/").filter(Boolean);
    const items: BreadcrumbItem[] = [
      { label: "Home", href: APP_ROUTES.OVERVIEW },
    ];

    let currentHref = "";
    segments.forEach((segment) => {
      currentHref += `/${segment}`;
      const capitalized = segment.charAt(0).toUpperCase() + segment.slice(1);
      items.push({ label: capitalized, href: currentHref });
    });

    return items;
  };

  const handleLogout = async () => {
    setIsLoggingOut(true);
    try {
      await logoutAdministratorService();
    } catch (e) {
      console.error("Logout failed", e);
    } finally {
      router.push(APP_ROUTES.AUTH.LOGIN);
    }
  };

  const currentBreadcrumbs = getAutoBreadcrumbs();

  const displayName = session?.administrator?.display_name || "Administrator";
  const initial = displayName.charAt(0).toUpperCase();

  return (
    <header className={styles.header}>
      <div className={styles.leftSection}>
        {onToggleMobileMenu && (
          <button
            type="button"
            className={styles.mobileMenuButton}
            onClick={onToggleMobileMenu}
            aria-label="Toggle navigation menu"
          >
            <Menu size={18} />
          </button>
        )}
        <Breadcrumbs items={currentBreadcrumbs} />
      </div>

      <div className={styles.rightSection}>
        <div className={styles.envTag}>
          <Sparkles size={12} aria-hidden="true" />
          <span>QVAC Local AI Engine</span>
        </div>

        <div className={styles.userProfilePill}>
          <div className={styles.avatar}>{initial}</div>
          <div className={styles.userInfo}>
            <span className={styles.userName}>{displayName}</span>
          </div>
          <button
            className={styles.logoutAction}
            onClick={handleLogout}
            disabled={isLoggingOut}
            title="Log out"
            aria-label="Log out"
          >
            {isLoggingOut ? (
              <Loader2 size={14} className="animate-spin" />
            ) : (
              <LogOut size={14} />
            )}
          </button>
        </div>
      </div>
    </header>
  );
};
