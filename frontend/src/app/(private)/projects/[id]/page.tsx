import React from "react";
import { ProjectDetailClient } from "@/features/projects/views/ProjectDetailView/ProjectDetailClient";
import { getProjectService } from "@/features/projects/services";

export default async function ProjectDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  let projectData;
  try {
    projectData = await getProjectService(id);
  } catch (error) {
    console.error(`[ProjectDetailPage] Failed to fetch project ${id}:`, error);
  }

  return <ProjectDetailClient initialProject={projectData} projectId={id} />;
}
