import { ProjectSettingsClient } from "@/features/projects/views/ProjectSettingsView/ProjectSettingsClient";

interface PageProps {
  params: {
    id: string;
  };
}

export default function ProjectSettingsPage({ params }: PageProps) {
  return <ProjectSettingsClient projectId={params.id} />;
}
