import Link from "next/link";
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

const GithubIcon = ({ size = 20, className = "" }: { size?: number; className?: string }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4" />
    <path d="M9 18c-4.51 2-5-2-7-2" />
  </svg>
);

const DokployIcon = ({ size = 20, className = "" }: { size?: number; className?: string }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
    <polyline points="3.27 6.96 12 12.01 20.73 6.96" />
    <line x1="12" y1="22.08" x2="12" y2="12" />
  </svg>
);

function ProjectCard({ project }: { project: ProjectSummary }) {
  const isHealthy = project.health_status === "healthy";
  
  return (
    <Link href={`/projects/${project.id}`} className={`${styles.card} ${isHealthy ? styles.cardHealthy : styles.cardWarning}`}>
      <div className={styles.header}>
        <div>
          <h3 className={styles.title}>{project.name}</h3>
          {project.description && <p className={styles.description}>{project.description}</p>}
        </div>
        <Badge variant={isHealthy ? "success" : "warning"}>
          {project.health_status.toUpperCase()}
        </Badge>
      </div>

      <div className={styles.integrations}>
        {project.github_repository ? (
          <div className={styles.integrationItem}>
            <GithubIcon size={14} className={styles.integrationIcon} />
            <span className={styles.integrationText}>{project.github_repository.full_name}</span>
          </div>
        ) : (
          <div className={styles.integrationItemEmpty}>
            <GithubIcon size={14} className={styles.integrationIconEmpty} />
            <span className={styles.integrationText}>No repository</span>
          </div>
        )}
        
        {project.dokploy_application ? (
          <div className={styles.integrationItem}>
            <DokployIcon size={14} className={styles.integrationIcon} />
            <span className={styles.integrationText}>
              {project.dokploy_application.display_name} 
              <span className={styles.integrationEnv}>({project.dokploy_application.environment || "env"})</span>
            </span>
          </div>
        ) : (
          <div className={styles.integrationItemEmpty}>
            <DokployIcon size={14} className={styles.integrationIconEmpty} />
            <span className={styles.integrationText}>No application</span>
          </div>
        )}
      </div>

      <div className={styles.footer}>
        <div className={styles.stat}>
          <span className={styles.statLabel}>Monitoring</span>
          <span className={styles.statValue}>
            <span className={`${styles.statusDot} ${styles[`status_${project.monitoring_status}`] || styles.status_inactive}`} />
            {project.monitoring_status}
          </span>
        </div>
      </div>
    </Link>
  );
}
