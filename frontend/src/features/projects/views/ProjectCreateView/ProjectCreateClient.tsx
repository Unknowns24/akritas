import React from "react";

export function ProjectCreateClient() {
  return (
    <div style={{ padding: "24px", maxWidth: "800px", margin: "0 auto" }}>
      <header style={{ marginBottom: "32px" }}>
        <h1 style={{ fontSize: "24px", fontWeight: "600", marginBottom: "8px" }}>Create New Project</h1>
        <p style={{ color: "var(--text-secondary)" }}>
          Set up a new project to start monitoring and analyzing incidents.
        </p>
      </header>

      <div style={{ backgroundColor: "var(--surface-1)", padding: "24px", borderRadius: "8px", border: "1px solid var(--surface-2)" }}>
        <p style={{ color: "var(--text-secondary)", fontStyle: "italic" }}>
          The project creation flow will be implemented here soon, integrating the GitHub Repository Selector and Dokploy Application Selector.
        </p>
      </div>
    </div>
  );
}
