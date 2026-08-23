"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import styles from "./SettingsNav.module.css";
import { APP_ROUTES } from "@/core/routes/routes.config";

const SETTINGS_LINKS = [
  { name: "General", href: APP_ROUTES.SETTINGS.ROOT },
  { name: "GitHub", href: "/settings/github" },
  { name: "Dokploy", href: "/settings/dokploy" },
  { name: "QVAC", href: "/settings/qvac" },
];

export function SettingsNav() {
  const pathname = usePathname();

  return (
    <nav className={styles.nav}>
      <ul className={styles.list}>
        {SETTINGS_LINKS.map((link) => {
          const isActive = pathname === link.href;

          return (
            <li key={link.name} className={styles.item}>
              <Link
                href={link.href}
                className={`${styles.link} ${isActive ? styles.active : ""}`}
              >
                {link.name}
              </Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}
