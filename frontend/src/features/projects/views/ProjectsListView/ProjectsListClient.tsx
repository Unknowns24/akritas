"use client";

import React, { useState, useMemo } from "react";
import Link from "next/link";
import { Search } from "lucide-react";
import { Button } from "@/core/ui/primitives/Button";
import { APP_ROUTES } from "@/core/routes/routes.config";
import { ProjectGrid } from "./components/ProjectGrid/ProjectGrid";
import { EmptyState } from "@/core/ui/feedback";
import type { components } from "@/core/libs/api-client";
import styles from "./ProjectsListView.module.css";

type ProjectSummary = components["schemas"]["ProjectSummary"];

interface ProjectsListClientProps {
  initialProjects: ProjectSummary[];
}

export function ProjectsListClient({ initialProjects }: ProjectsListClientProps) {
  const [searchQuery, setSearchQuery] = useState("");

  const filteredProjects = useMemo(() => {
    if (!searchQuery.trim()) {
      return initialProjects;
    }
    const query = searchQuery.toLowerCase().trim();
    return initialProjects.filter(project => 
      project.name.toLowerCase().includes(query) || 
      (project.health_status && project.health_status.toLowerCase().includes(query)) ||
      (project.monitoring_status && project.monitoring_status.toLowerCase().includes(query))
    );
  }, [initialProjects, searchQuery]);

  return (
    <>
      <div className={styles.actionBar}>
        <div className={styles.searchWrapper}>
          <Search size={16} className={styles.searchIcon} />
          <input 
            type="text" 
            placeholder="Search projects..." 
            className={styles.searchInput}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        
        <Link href={APP_ROUTES.PROJECTS.NEW}>
          <Button variant="ghost" size="md" className={styles.newProjectBtn}>
            New Project
          </Button>
        </Link>
      </div>
      
      <div className={styles.divider} />

      {filteredProjects.length > 0 ? (
        <ProjectGrid projects={filteredProjects} />
      ) : (
        <EmptyState 
          icon={<Search size={32} />}
          title="No projects found"
          description={`We couldn't find any projects matching "${searchQuery}". Try adjusting your search query.`}
        />
      )}
    </>
  );
}
