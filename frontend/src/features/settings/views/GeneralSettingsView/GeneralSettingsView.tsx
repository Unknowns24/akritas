import { getSystemStatusService } from "@/features/settings/services";
import { IntegrationStatusCard, SystemDiagnosticsCard } from "./components";

export async function GeneralSettingsView() {
  const status = await getSystemStatusService();

  return (
    <div>
      <IntegrationStatusCard status={status.data} />
      <SystemDiagnosticsCard status={status.data} />
    </div>
  );
}
