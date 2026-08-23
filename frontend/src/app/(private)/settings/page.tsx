import { Suspense } from "react";
import { GeneralSettingsView } from "@/features/settings/views";

export default function SettingsPage() {
  return (
    <Suspense fallback={<div style={{ color: "var(--text-dim)", fontSize: "14px" }}>Loading settings...</div>}>
      <GeneralSettingsView />
    </Suspense>
  );
}
