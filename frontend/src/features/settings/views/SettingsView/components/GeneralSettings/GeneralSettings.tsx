import React from "react";
import { Bot, Code, Database, Play, Search } from "lucide-react";
import styles from "./GeneralSettings.module.css";

export const GeneralSettings: React.FC = () => {
  return (
    <div className={styles.generalContainer}>
      <div className={styles.header}>
        <h1 className={styles.pageTitle}>Settings</h1>
        <p className={styles.pageSubtitle}>
          Configure Akritas integrations, infrastructure and autonomous investigation behavior.
        </p>
      </div>

      <div className={styles.section}>
        <h2 className={styles.sectionTitle}>Integration Status</h2>
        <div className={styles.integrationList}>
          <div className={styles.integrationItem}>
            <div className={styles.integrationIconWrapper}>
              <Code size={20} />
            </div>
            <div className={styles.integrationInfo}>
              <span className={styles.integrationName}>GitHub</span>
              <div className={styles.integrationMeta}>
                <span>1 account</span>
                <span className={styles.dot}></span>
                <span className={styles.statusConnected}>Connected</span>
              </div>
            </div>
          </div>

          <div className={styles.integrationItem}>
            <div className={styles.integrationIconWrapper}>
              <Database size={20} />
            </div>
            <div className={styles.integrationInfo}>
              <span className={styles.integrationName}>Dokploy</span>
              <div className={styles.integrationMeta}>
                <span>2 servers</span>
                <span className={styles.dot}></span>
                <span className={styles.statusConnected}>Connected</span>
              </div>
            </div>
          </div>

          <div className={styles.integrationItem}>
            <div className={styles.integrationIconWrapper}>
              <Bot size={20} />
            </div>
            <div className={styles.integrationInfo}>
              <span className={styles.integrationName}>QVAC</span>
              <div className={styles.integrationMeta}>
                <span>http://localhost:11434</span>
                <span className={styles.dot}></span>
                <span className={styles.statusConnected}>Connected</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className={styles.section}>
        <div className={styles.sectionHeader}>
          <h2 className={styles.sectionTitle}>System Diagnostics</h2>
          <span className={styles.lastCheck}>LAST SYSTEM CHECK 46S AGO</span>
        </div>
        <div className={styles.diagnosticsGrid}>
          <div className={styles.diagnosticCard}>
            <div className={styles.diagnosticHeader}>
              <Code size={14} />
              <span>GITHUB API</span>
            </div>
            <div className={styles.diagnosticStatus}>
              <span className={`${styles.statusIndicator} ${styles.statusHealthy}`}></span>
              <span>Healthy</span>
            </div>
          </div>

          <div className={styles.diagnosticCard}>
            <div className={styles.diagnosticHeader}>
              <Bot size={14} />
              <span>QVAC</span>
            </div>
            <div className={styles.diagnosticStatus}>
              <span className={`${styles.statusIndicator} ${styles.statusHealthy}`}></span>
              <span>Healthy</span>
            </div>
          </div>

          <div className={styles.diagnosticCard}>
            <div className={styles.diagnosticHeader}>
              <Database size={14} />
              <span>DOKPLOY</span>
            </div>
            <div className={styles.diagnosticStatus}>
              <span className={`${styles.statusIndicator} ${styles.statusHealthy}`}></span>
              <span>Healthy</span>
            </div>
          </div>

          <div className={`${styles.diagnosticCard} ${styles.cardHighlight}`}>
            <div className={styles.diagnosticHeader}>
              <Search size={14} />
              <span>INVESTIGATOR</span>
            </div>
            <div className={styles.diagnosticStatus}>
              <span className={`${styles.statusIndicator} ${styles.statusRunning}`}></span>
              <span>Running</span>
            </div>
          </div>
        </div>
        <div className={styles.diagnosticsActions}>
          <button className={styles.runDiagnosticsBtn}>
            <Play size={14} />
            Run Diagnostics
          </button>
        </div>
      </div>
    </div>
  );
};
