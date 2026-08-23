import React from "react";
import styles from "./AuthLayout.module.css";
import { AkritasLogo } from "@/core/ui";

export default function AuthLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className={styles.container}>
      <div className={styles.logo}>
        <AkritasLogo variant={"logo-text-white"} size={180} />
      </div>
      {children}
    </div>
  );
}
