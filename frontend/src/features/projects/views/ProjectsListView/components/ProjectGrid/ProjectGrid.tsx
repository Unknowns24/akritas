import Link from "next/link";
import { GithubIcon, DokployIcon } from "@/core/ui/icons";
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

function formatRelativeTime(dateString: string) {
  const diff = Date.now() - new Date(dateString).getTime();
  const minutes = Math.floor(diff / 60000);
  if (minutes < 1) return "Just now";
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  return `${Math.floor(hours / 24)}d ago`;
}

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
        
        {project.dokploy_source ? (
          <div className={styles.integrationItem}>
            <DokployIcon size={14} className={styles.integrationIcon} />
            <span className={styles.integrationText}>
              {project.dokploy_source.display_name}
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
            <span style={{ textTransform: "capitalize" }}>{project.monitoring_status}</span>
          </span>
        </div>
        <div className={styles.stat}>
          <span className={styles.statLabel}>Ingestion</span>
          <span className={`${styles.statValue} ${!project.last_observed_at ? styles.textWarning : ""}`}>
            {project.last_observed_at ? formatRelativeTime(project.last_observed_at) : "Awaiting first log"}
          </span>
        </div>
      </div>
    </Link>
  );
}
