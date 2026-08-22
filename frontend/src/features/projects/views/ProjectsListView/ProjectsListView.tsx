
import React from "react";
import Link from "next/link";
import { FolderGit2, Plus } from "lucide-react";
import { Button } from "@/core/ui/primitives/Button";
import { APP_ROUTES } from "@/core/routes/routes.config";
import styles from "./ProjectsListView.module.css";

export const ProjectsListView: React.FC = () => {
  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.titleGroup}>
          <h1 className={styles.pageTitle}>Monitored Projects</h1>
          <p className={styles.pageSubtitle}>
            Dokploy application bindings, GitHub repositories, and log monitoring configuration.
          </p>
        </div>
        <Link href={APP_ROUTES.PROJECTS.NEW}>
          <Button variant="primary" size="md" leftIcon={<Plus size={16} />}>
            New Project
          </Button>
        </Link>
      </div>

      <div className={styles.emptyCard}>
        <FolderGit2 size={40} className={styles.emptyIcon} />
        <h2 className={styles.emptyTitle}>No Projects Configured</h2>
        <p className={styles.emptyText}>
          Create a project to bind a Dokploy application with its corresponding GitHub repository
          and begin real-time incident monitoring.
        </p>
        <Link href={APP_ROUTES.PROJECTS.NEW} style={{ marginTop: "8px" }}>
          <Button variant="secondary" size="md" leftIcon={<Plus size={16} />}>
            Create First Project
          </Button>
        </Link>
      </div>
    </div>
  );
};
