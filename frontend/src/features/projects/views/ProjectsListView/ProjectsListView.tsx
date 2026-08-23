import React from "react";
import { getProjectsService } from "../../services/get-projects.service";
import { ProjectsListClient } from "./ProjectsListClient";
import styles from "./ProjectsListView.module.css";
import { EmptyState, ErrorState } from "@/core/ui/feedback";
import { FolderGit2 } from "lucide-react";

export const ProjectsListView = async () => {
  let projectResponse;
  
  try {
    projectResponse = await getProjectsService();
  } catch (error) {
    return (
      <div className={styles.container}>
        <div className={styles.header}>
          <div className={styles.titleGroup}>
            <h1 className={styles.pageTitle}>Projects</h1>
          </div>
        </div>
        <ErrorState error={error as Error} />
      </div>
    );
  }

  const projects = projectResponse.data || [];
  const hasProjects = projects.length > 0;

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.titleGroup}>
          <h1 className={styles.pageTitle}>Manage monitored applications and their connection to infrastructure.</h1>
        </div>
      </div>

      {!hasProjects ? (
        <EmptyState 
          icon={<FolderGit2 size={48} />}
          title="No projects configured"
          description="You haven't connected any projects yet. Create your first project to start monitoring your infrastructure and deployments."
        />
      ) : (
        <ProjectsListClient initialProjects={projects} />
      )}
    </div>
  );
};
