"use client";

import React, { useState } from "react";
import { SettingsNav } from "./components/SettingsNav/SettingsNav";
import { GeneralSettings } from "./components/GeneralSettings/GeneralSettings";
import styles from "./SettingsView.module.css";

const tabs = ["General", "GitHub", "Dokploy", "QVAC", "Automation"];

export const SettingsView: React.FC = () => {
  const [activeTab, setActiveTab] = useState("General");

  return (
    <div className={styles.container}>
      <SettingsNav tabs={tabs} activeTab={activeTab} onTabChange={setActiveTab} />

      {activeTab === "General" && <GeneralSettings />}
      {activeTab !== "General" && (
        <div className={styles.specificSection}>
          <h2>{activeTab} Settings</h2>
          <p>Configuration options for {activeTab} will appear here.</p>
        </div>
      )}
    </div>
  );
};
