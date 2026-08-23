import { IncidentDetailView } from "@/features/incidents/views/IncidentDetailView/IncidentDetailView";
import { Suspense } from "react";
import { Loader2 } from "lucide-react";

interface IncidentDetailPageProps {
  params: {
    id: string;
  };
}

export default function IncidentDetailPage({ params }: IncidentDetailPageProps) {
  return (
    <Suspense
      fallback={
        <div style={{ display: 'flex', justifyContent: 'center', padding: '40px' }}>
          <Loader2 className="animate-spin" size={32} />
        </div>
      }
    >
      <IncidentDetailView id={params.id} />
    </Suspense>
  );
}
