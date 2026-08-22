import React from "react";
import Link from "next/link";
import { Plus } from "lucide-react";
import { Button } from "@/core/ui/primitives/Button";
import { APP_ROUTES } from "@/core/routes/routes.config";
import { MetricsGrid } from "./components/MetricsGrid/MetricsGrid";
import { ActiveIncidentsCard } from "./components/ActiveIncidentsCard/ActiveIncidentsCard";
import { InvestigationEngineBanner } from "./components/InvestigationEngineBanner/InvestigationEngineBanner";
import { PipelineLifecycleCard } from "./components/PipelineLifecycleCard/PipelineLifecycleCard";
import styles from "./OverviewView.module.css";

export const OverviewView: React.FC = () => {
  return (
    <div className={styles.container}>
      {/* Header section */}
      <div className={styles.header}>
        <div className={styles.titleGroup}>
          <h1 className={styles.pageTitle}>Operational Overview</h1>
          <p className={styles.pageSubtitle}>
            Autonomous monitoring, local QVAC investigation, and validated remediation.
          </p>
        </div>
        <Link href={APP_ROUTES.PROJECTS.NEW}>
          <Button variant="primary" size="md" leftIcon={<Plus size={16} />}>
            New Project
          </Button>
        </Link>
      </div>

      <MetricsGrid />

      {/* Main Content Grid */}
      <div className={styles.contentGrid}>
        {/* Left Column: Recent Incidents / Activity */}
        <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
          <ActiveIncidentsCard />
        </div>

        {/* Right Column: AI Intel & Pipeline Status */}
        <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
          <InvestigationEngineBanner />
          <PipelineLifecycleCard />
        </div>
      </div>
    </div>
  );
};
