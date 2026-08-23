import React from "react";
import { ProjectDetailClient } from "@/features/projects/views/ProjectDetailView/ProjectDetailClient";
import { getProjectService } from "@/features/projects/services";

export default async function ProjectDetailPage({
  params,
}: {
  params: { id: string };
}) {
  // En Next.js 15, los params en app router deben tratarse de forma síncrona si los consumimos directamente, o usar page props
  // Para simplificar, accedemos directo a params.id
  const projectId = params.id;

  // Obtenemos el proyecto usando el mock o API
  const projectData = await getProjectService(projectId);

  return <ProjectDetailClient initialProject={projectData} />;
}
