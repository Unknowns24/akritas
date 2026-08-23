import React from "react";
import { getProjectsService } from "../../services/get-projects.service";
import { ProjectsListClient } from "./ProjectsListClient";
import styles from "./ProjectsListView.module.css";
import { EmptyState, ErrorState } from "@/core/ui/feedback";
import { FolderGit2 } from "lucide-react";

export const ProjectsListView = async () => {
  let initialProjects = [];
  try {
    const projectResponse = await getProjectsService();
    initialProjects = projectResponse.data || [];
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

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.titleGroup}>
          <h1 className={styles.pageTitle}>Projects</h1>
          <p className={styles.pageSubtitle}>
            Manage monitored applications and their connection to infrastructure.
          </p>
        </div>
      </div>

      <ProjectsListClient initialProjects={initialProjects} />
    </div>
  );
};
