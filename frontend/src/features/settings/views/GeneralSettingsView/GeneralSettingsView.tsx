import { AlertTriangle } from "lucide-react";
import { getErrorMessage, isApiNotFoundError } from "@/core/errors";
import { Card, CardBody, CardHeader } from "@/core/ui/primitives/Card";
import { getSystemStatusService } from "@/features/settings/services";
import { IntegrationStatusCard, SystemDiagnosticsCard } from "./components";
import styles from "./GeneralSettingsView.module.css";

export async function GeneralSettingsView() {
  let status;

  try {
    status = await getSystemStatusService();
  } catch (error) {
    if (!isApiNotFoundError(error)) {
      throw error;
    }

    return (
      <Card>
        <CardHeader>
          <div className={styles.unavailableHeader}>
            <AlertTriangle size={16} />
            <span>System status endpoint unavailable</span>
          </div>
        </CardHeader>
        <CardBody>
          <p className={styles.unavailableText}>
            {getErrorMessage(error, "System status endpoint unavailable.")}
          </p>
        </CardBody>
      </Card>
    );
  }

  return (
    <div>
      <IntegrationStatusCard status={status.data} />
      <SystemDiagnosticsCard status={status.data} />
    </div>
  );
}
