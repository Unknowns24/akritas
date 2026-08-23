import { ProjectSettingsClient } from "@/features/projects/views/ProjectSettingsView/ProjectSettingsClient";

interface PageProps {
  params: Promise<{
    id: string;
  }>;
}

export default async function ProjectSettingsPage({ params }: PageProps) {
  const { id } = await params;
  return <ProjectSettingsClient projectId={id} />;
}
