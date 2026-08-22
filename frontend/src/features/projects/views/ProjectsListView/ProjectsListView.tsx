
import React from "react";
import Link from "next/link";
import { Plus, Search } from "lucide-react";
import { Button } from "@/core/ui/primitives/Button";
import { APP_ROUTES } from "@/core/routes/routes.config";
import { getProjectsService } from "../../services/get-projects.service";
import { ProjectGrid } from "./components/ProjectGrid/ProjectGrid";
import styles from "./ProjectsListView.module.css";

export const ProjectsListView = async () => {
  let projectResponse;
  
  try {
    projectResponse = await getProjectsService();
  } catch (error) {
    throw error;
  }

  const hasProjects = projectResponse.data && projectResponse.data.length > 0;

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.titleGroup}>
          <h1 className={styles.pageTitle}>Manage monitored applications and their connection to infrastructure.</h1>
        </div>
      </div>

      <div className={styles.actionBar}>
        <div className={styles.searchWrapper}>
          <Search size={16} className={styles.searchIcon} />
          <input 
            type="text" 
            placeholder="Search projects..." 
            className={styles.searchInput}
          />
        </div>
        
        <Link href={APP_ROUTES.PROJECTS.NEW}>
          <Button variant="ghost" size="md" className={styles.newProjectBtn}>
            New Project
          </Button>
        </Link>
      </div>
      
      <div className={styles.divider} />

      {hasProjects ? (
        <ProjectGrid projects={projectResponse.data} />
      ) : (
        <ProjectGrid projects={[]} />
      )}
    </div>
  );
};
