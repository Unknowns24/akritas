"use client";

import React, { useState } from "react";
import { Sidebar } from "@/core/ui/layout/Sidebar";
import { Header } from "@/core/ui/layout/Header";
import styles from "./AppShell.module.css";

export interface AppShellProps {
  children: React.ReactNode;
}

export const AppShell: React.FC<AppShellProps> = ({ children }) => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const toggleMobileMenu = () => {
    setMobileMenuOpen((prev) => !prev);
  };

  const closeMobileMenu = () => {
    setMobileMenuOpen(false);
  };

  return (
    <div className={styles.wrapper}>
      {/* Mobile Backdrop */}
      <div
        className={`${styles.backdrop} ${mobileMenuOpen ? styles.backdropOpen : ""}`.trim()}
        onClick={closeMobileMenu}
        aria-hidden="true"
      />

      {/* Sidebar */}
      <Sidebar isOpen={mobileMenuOpen} onClose={closeMobileMenu} />

      {/* Main Content Flow */}
      <div className={styles.contentWrapper}>
        <Header onToggleMobileMenu={toggleMobileMenu} />
        <main className={styles.mainContent}>{children}</main>
      </div>
    </div>
  );
};
