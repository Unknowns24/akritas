"use client";

import React, { useEffect, useState } from "react";
import { ProjectForm } from "@/features/projects/components/ProjectForm/ProjectForm";
import { Project, getProjectService } from "@/features/projects/services/get-project.service";
import { Loader2 } from "lucide-react";

interface ProjectSettingsClientProps {
  projectId: string;
}

export function ProjectSettingsClient({ projectId }: ProjectSettingsClientProps) {
  const [project, setProject] = useState<Project | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchProject = async () => {
      setIsLoading(true);
      const { data } = await getProjectService(projectId);
      if (data) {
        setProject(data);
      }
      setIsLoading(false);
    };

    fetchProject();
  }, [projectId]);

  if (isLoading) {
    return (
      <div style={{ display: "flex", justifyContent: "center", padding: "64px" }}>
        <Loader2 className="animate-spin" size={32} style={{ color: "var(--brand-500)" }} />
      </div>
    );
  }

  if (!project) {
    return (
      <div style={{ padding: "24px", color: "var(--red-500)" }}>
        Failed to load project settings.
      </div>
    );
  }

  return (
    <div style={{ padding: "24px", maxWidth: "800px", margin: "0 auto" }}>
      <header style={{ marginBottom: "32px" }}>
        <h1 style={{ fontSize: "24px", fontWeight: "600", marginBottom: "8px" }}>Project Settings</h1>
        <p style={{ color: "var(--text-secondary)" }}>
          Manage your project configuration and integrations.
        </p>
      </header>

      <ProjectForm initialData={project} />
    </div>
  );
}
