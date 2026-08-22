import { Code2, Server } from "lucide-react";
import { Badge } from "@/core/ui/primitives/Badge";
import styles from "./ProjectGrid.module.css";
import type { components } from "@/core/libs/api-client";

type ProjectSummary = components["schemas"]["ProjectSummary"];

interface ProjectGridProps {
  projects: ProjectSummary[];
}

export function ProjectGrid({ projects }: ProjectGridProps) {
  return (
    <div className={styles.grid}>
      {projects.map((project) => (
        <ProjectCard key={project.id} project={project} />
      ))}
    </div>
  );
}

function ProjectCard({ project }: { project: ProjectSummary }) {
  const isHealthy = project.health_status === "healthy";
  
  return (
    <div className={`${styles.card} ${isHealthy ? styles.cardHealthy : styles.cardWarning}`}>
      <div className={styles.header}>
        <h3 className={styles.title}>{project.name}</h3>
        <Badge variant={isHealthy ? "success" : "warning"}>
          {project.health_status.toUpperCase()}
        </Badge>
      </div>

      <div className={styles.integrations}>
        {project.github_repository && (
          <div className={styles.integrationIcon}>
            <Code2 size={16} />
          </div>
        )}
        {project.dokploy_application && (
          <div className={styles.integrationIcon}>
            <Server size={16} />
          </div>
        )}
      </div>

      <div className={styles.footer}>
        <div className={styles.stat}>
          <span className={styles.statLabel}>Status</span>
          <span className={styles.statValue}>{project.monitoring_status}</span>
        </div>
      </div>
    </div>
  );
}
