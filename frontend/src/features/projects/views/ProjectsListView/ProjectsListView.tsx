import React from "react";
import { getProjectsService } from "../../services/get-projects.service";
import { ProjectsListClient } from "./ProjectsListClient";
import styles from "./ProjectsListView.module.css";
import { FolderGit2 } from "lucide-react";

export const ProjectsListView = async () => {
  let projectResponse;
  
  try {
    projectResponse = await getProjectsService();
  } catch (error) {
    throw error;
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
        <div className={styles.emptyCard}>
          <FolderGit2 size={48} className={styles.emptyIcon} />
          <h3 className={styles.emptyTitle}>No projects configured</h3>
          <p className={styles.emptyText}>
            You haven't connected any projects yet. Create your first project to start monitoring your infrastructure and deployments.
          </p>
        </div>
      ) : (
        <ProjectsListClient initialProjects={projects} />
      )}
    </div>
  );
};
