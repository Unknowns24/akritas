"use client";

import React, { useCallback, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { Search, FolderGit2, RefreshCw } from "lucide-react";
import { Button } from "@/core/ui/primitives/Button";
import { APP_ROUTES } from "@/core/routes/routes.config";
import { ProjectGrid } from "./components/ProjectGrid/ProjectGrid";
import { EmptyState } from "@/core/ui/feedback";
import type { components } from "@/core/libs/api-client";
import { getProjectsService } from "../../services/get-projects.service";
import { getProjectService } from "../../services/get-project.service";
import { updateMonitoringConfigService } from "../../services/update-monitoring-config.service";
import { getErrorMessage } from "@/core/errors";
import { toast } from "sonner";
import styles from "./ProjectsListView.module.css";

type ProjectSummary = components["schemas"]["ProjectSummary"];

interface ProjectsListClientProps {
  initialProjects: ProjectSummary[];
}

export function ProjectsListClient({ initialProjects }: ProjectsListClientProps) {
  const [projects, setProjects] = useState<ProjectSummary[]>(initialProjects);
  const [searchQuery, setSearchQuery] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [updatingProjectId, setUpdatingProjectId] = useState<string | null>(null);

  const fetchProjects = useCallback(async () => {
    try {
      setIsLoading(true);
      const res = await getProjectsService();
      if (res?.data) {
        setProjects(res.data);
      }
    } catch (err) {
      console.error("Failed to load projects client-side:", err);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const handleToggleProjectMonitoring = useCallback(async (project: ProjectSummary) => {
    const shouldEnable = project.monitoring_status === "disabled";

    try {
      setUpdatingProjectId(project.id);
      const projectResponse = await getProjectService(project.id);
      const config = projectResponse.data.monitoring_configuration;

      await updateMonitoringConfigService(project.id, {
        ...config,
        enabled: shouldEnable,
      });

      toast.success(`Project ${shouldEnable ? "activated" : "deactivated"} successfully`);
      await fetchProjects();
    } catch (error: unknown) {
      toast.error(getErrorMessage(error, "Failed to update project monitoring."));
    } finally {
      setUpdatingProjectId(null);
    }
  }, [fetchProjects]);

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      void fetchProjects();
    }, 0);

    return () => window.clearTimeout(timeoutId);
  }, [fetchProjects]);

  const filteredProjects = useMemo(() => {
    if (!searchQuery.trim()) {
      return projects;
    }
    const query = searchQuery.toLowerCase().trim();
    return projects.filter(project => 
      project.name.toLowerCase().includes(query) || 
      (project.health_status && project.health_status.toLowerCase().includes(query)) ||
      (project.monitoring_status && project.monitoring_status.toLowerCase().includes(query))
    );
  }, [projects, searchQuery]);

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
        
        <div style={{ display: "flex", gap: "8px", alignItems: "center" }}>
          <Button 
            variant="secondary" 
            size="md" 
            onClick={fetchProjects}
            disabled={isLoading}
            title="Refresh projects"
          >
            <RefreshCw size={16} />
          </Button>

          <Link href={APP_ROUTES.PROJECTS.NEW}>
            <Button variant="primary" size="md" className={styles.newProjectBtn}>
              New Project
            </Button>
          </Link>
        </div>
      </div>
      
      <div className={styles.divider} />

      {filteredProjects.length > 0 ? (
        <ProjectGrid
          projects={filteredProjects}
          updatingProjectId={updatingProjectId}
          onToggleProjectMonitoring={handleToggleProjectMonitoring}
        />
      ) : projects.length === 0 ? (
        <EmptyState 
          icon={<FolderGit2 size={48} />}
          title="No projects configured"
          description="You haven't connected any projects yet. Create your first project to start monitoring your infrastructure and deployments."
          action={
            <Link href={APP_ROUTES.PROJECTS.NEW}>
              <Button variant="primary" size="md">
                Create Project
              </Button>
            </Link>
          }
        />
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
