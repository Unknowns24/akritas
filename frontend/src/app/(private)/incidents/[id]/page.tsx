import { IncidentDetailView } from "@/features/incidents/views/IncidentDetailView/IncidentDetailView";
import { Suspense } from "react";
import { Loader2 } from "lucide-react";

interface IncidentDetailPageProps {
  params: Promise<{
    id: string;
  }>;
}

export default async function IncidentDetailPage({
  params,
}: IncidentDetailPageProps) {
  const { id } = await params;
  
  return (
    <Suspense
      fallback={
        <div
          style={{ display: "flex", justifyContent: "center", padding: "40px" }}
        >
          <Loader2 className="animate-spin" size={32} />
        </div>
      }
    >
      <IncidentDetailView id={id} />
    </Suspense>
  );
}
