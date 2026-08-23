"use client";

import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Activity, AlertTriangle, FolderGit2, Settings } from "lucide-react";
import { APP_ROUTES } from "@/core/routes/routes.config";
import { StatusPip } from "@/core/ui/primitives/StatusPip";
import { AkritasLogo } from "@/core/ui/primitives/AkritasLogo";
import styles from "./Sidebar.module.css";

export interface SidebarProps {
  isOpen?: boolean;
  onClose?: () => void;
}

interface NavItemConfig {
  label: string;
  href: string;
  icon: React.ReactNode;
  matchPrefix?: boolean;
}

const NAV_ITEMS: NavItemConfig[] = [
  {
    label: "Overview",
    href: APP_ROUTES.OVERVIEW,
    icon: <Activity size={18} />,
  },
  {
    label: "Incidents",
    href: APP_ROUTES.INCIDENTS.LIST,
    icon: <AlertTriangle size={18} />,
    matchPrefix: true,
  },
  {
    label: "Projects",
    href: APP_ROUTES.PROJECTS.LIST,
    icon: <FolderGit2 size={18} />,
    matchPrefix: true,
  },
  {
    label: "Settings",
    href: APP_ROUTES.SETTINGS.ROOT,
    icon: <Settings size={18} />,
    matchPrefix: true,
  },
];

export const Sidebar: React.FC<SidebarProps> = ({
  isOpen = false,
  onClose,
}) => {
  const pathname = usePathname();

  const isItemActive = (item: NavItemConfig) => {
    if (item.href === APP_ROUTES.OVERVIEW) {
      return pathname === APP_ROUTES.OVERVIEW;
    }
    if (item.matchPrefix) {
      return pathname.startsWith(item.href);
    }
    return pathname === item.href;
  };

  return (
    <aside
      className={`${styles.sidebar} ${isOpen ? styles.sidebarOpen : ""}`.trim()}
    >
      <div className={styles.brand}>
        <AkritasLogo size={128} />
      </div>

      <nav className={styles.nav} aria-label="Main Navigation">
        <span className={styles.sectionLabel}>Operations</span>
        {NAV_ITEMS.map((item) => {
          const active = isItemActive(item);

          return (
            <Link
              key={item.href}
              href={item.href}
              className={`${styles.navItem} ${active ? styles.navItemActive : ""}`.trim()}
              onClick={onClose}
              aria-current={active ? "page" : undefined}
            >
              <span className={styles.navIcon}>{item.icon}</span>
              <span>{item.label}</span>
            </Link>
          );
        })}
      </nav>

      <div className={styles.footer}>
        <StatusPip status="healthy" pulse label="System Status" />
      </div>
    </aside>
  );
};
