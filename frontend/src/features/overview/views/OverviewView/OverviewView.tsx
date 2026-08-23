import React from "react";
import Link from "next/link";
import { AlertTriangle, Plus } from "lucide-react";
import { getErrorMessage, isApiNotFoundError } from "@/core/errors";
import { Button } from "@/core/ui/primitives/Button";
import { Card, CardBody, CardHeader } from "@/core/ui/primitives/Card";
import { APP_ROUTES } from "@/core/routes/routes.config";
import { MetricsGrid } from "./components/MetricsGrid/MetricsGrid";
import { ActiveIncidentsCard } from "./components/ActiveIncidentsCard/ActiveIncidentsCard";
import { InvestigationEngineBanner } from "./components/InvestigationEngineBanner/InvestigationEngineBanner";
import { PipelineLifecycleCard } from "./components/PipelineLifecycleCard/PipelineLifecycleCard";
import { getOverviewService } from "../../services/get-overview.service";
import styles from "./OverviewView.module.css";

export const OverviewView = async () => {
  let overview;
  let unavailableMessage: string | null = null;

  try {
    const res = await getOverviewService();
    overview = res.data;
  } catch (error) {
    if (!isApiNotFoundError(error)) {
      throw error;
    }

    unavailableMessage = getErrorMessage(error, "Overview endpoint unavailable.");
  }

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

      {overview ? (
        <>
          <MetricsGrid
            monitored_projects={overview.monitored_projects}
            active_incidents={overview.active_incidents}
            workflow_completed_incidents={overview.workflow_completed_incidents}
            pull_requests_created={overview.pull_requests_created}
          />

          <div className={styles.contentGrid}>
            <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
              <ActiveIncidentsCard investigations={overview.active_investigations} />
            </div>

            <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
              <InvestigationEngineBanner />
              <PipelineLifecycleCard />
            </div>
          </div>
        </>
      ) : (
        <Card>
          <CardHeader>
            <div className={styles.unavailableHeader}>
              <AlertTriangle size={16} />
              <span>Overview endpoint unavailable</span>
            </div>
          </CardHeader>
          <CardBody>
            <p className={styles.unavailableText}>{unavailableMessage}</p>
          </CardBody>
        </Card>
      )}
    </div>
  );
};
