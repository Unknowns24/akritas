"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Play } from "lucide-react";
import { toast } from "sonner";
import { getErrorMessage } from "@/core/errors/format-error";
import { Button } from "@/core/ui/primitives/Button";
import { startIncidentInvestigationService } from "../../../services/start-incident-investigation.service";
import type { Incident } from "../../../services/get-incident.service";
import styles from "../IncidentDetailView.module.css";

type InvestigationActionButtonProps = {
  incidentId: string;
  phase: Incident["phase"];
};

export function InvestigationActionButton({ incidentId, phase }: InvestigationActionButtonProps) {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const canStart = phase === "detected" || phase === "failed";

  if (!canStart) {
    return null;
  }

  async function handleStart() {
    setIsLoading(true);
    try {
      await startIncidentInvestigationService(incidentId);
      toast.success("Investigation queued");
      router.refresh();
    } catch (error) {
      toast.error(getErrorMessage(error, "Failed to start investigation"));
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <Button
      type="button"
      variant="primary"
      className={styles.actionButton}
      isLoading={isLoading}
      onClick={handleStart}
    >
      <Play size={14} />
      {phase === "failed" ? "Retry investigation" : "Run investigation"}
    </Button>
  );
}
