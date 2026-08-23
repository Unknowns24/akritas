import { SettingsNav } from "@/features/settings/components";
import styles from "./SettingsLayout.module.css";

export default function SettingsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>Settings</h1>
        <p className={styles.description}>
          Configure Akritas integrations, infrastructure and autonomous investigation behavior.
        </p>
      </div>

      <SettingsNav />

      <main className={styles.main}>{children}</main>
    </div>
  );
}
