"use client";

import React from "react";
import { usePathname } from "next/navigation";
import { Menu, Sparkles } from "lucide-react";
import { Breadcrumbs, BreadcrumbItem } from "@/core/ui/layout/Breadcrumbs";
import { APP_ROUTES } from "@/core/routes/routes.config";
import styles from "./Header.module.css";

export interface HeaderProps {
  onToggleMobileMenu?: () => void;
  breadcrumbs?: BreadcrumbItem[];
}

export const Header: React.FC<HeaderProps> = ({ onToggleMobileMenu, breadcrumbs }) => {
  const pathname = usePathname();

  const getAutoBreadcrumbs = (): BreadcrumbItem[] => {
    if (breadcrumbs && breadcrumbs.length > 0) {
      return breadcrumbs;
    }

    if (pathname === APP_ROUTES.OVERVIEW) {
      return [{ label: "Dashboard", href: APP_ROUTES.OVERVIEW }];
    }

    const segments = pathname.split("/").filter(Boolean);
    const items: BreadcrumbItem[] = [{ label: "Home", href: APP_ROUTES.OVERVIEW }];

    let currentHref = "";
    segments.forEach((segment) => {
      currentHref += `/${segment}`;
      const capitalized = segment.charAt(0).toUpperCase() + segment.slice(1);
      items.push({ label: capitalized, href: currentHref });
    });

    return items;
  };

  const currentBreadcrumbs = getAutoBreadcrumbs();

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
      </div>
    </header>
  );
};
