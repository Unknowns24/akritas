import React from "react";
import { Bell } from "lucide-react";
import styles from "./SettingsNav.module.css";

export interface SettingsNavProps {
  tabs: string[];
  activeTab: string;
  onTabChange: (tab: string) => void;
}

export const SettingsNav: React.FC<SettingsNavProps> = ({ tabs, activeTab, onTabChange }) => {
  return (
    <div className={styles.topNav}>
      <div className={styles.tabs}>
        {tabs.map((tab) => (
          <button
            key={tab}
            className={`${styles.tab} ${activeTab === tab ? styles.activeTab : ""}`.trim()}
            onClick={() => onTabChange(tab)}
          >
            {tab}
          </button>
        ))}
      </div>
      <div className={styles.topNavActions}>
        <button className={styles.iconButton} aria-label="Notifications">
          <Bell size={18} />
        </button>
        <div className={styles.avatar}>
          <img src="https://api.dicebear.com/7.x/avataaars/svg?seed=Akritas" alt="User Avatar" />
        </div>
      </div>
    </div>
  );
};
