import React from "react";
import { ProjectForm } from "@/features/projects/components/ProjectForm/ProjectForm";

export function ProjectCreateClient() {
  return (
    <div style={{ padding: "24px", maxWidth: "800px", margin: "0 auto" }}>
      <header style={{ marginBottom: "32px" }}>
        <h1 style={{ fontSize: "24px", fontWeight: "600", marginBottom: "8px" }}>Create New Project</h1>
        <p style={{ color: "var(--text-secondary)" }}>
          Set up a new project to start monitoring and analyzing incidents.
        </p>
      </header>

      <ProjectForm />
    </div>
  );
}
