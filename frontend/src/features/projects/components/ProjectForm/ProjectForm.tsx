"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import styles from "./ProjectForm.module.css";
import { Button } from "@/core/ui/primitives/Button";
import { AlertTriangle } from "lucide-react";
import { GitHubRepositorySelector } from "@/features/projects/components/GitHubRepositorySelector/GitHubRepositorySelector";
import { DokployApplicationSelector } from "@/features/projects/components/DokployApplicationSelector/DokployApplicationSelector";
import type { Project } from "@/features/projects/services/get-project.service";
import type { GitHubRepository } from "@/features/settings/services/github/get-github-repositories.service";
import type { DokployApplication } from "@/features/settings/services/dokploy/get-dokploy-applications.service";
import { createProjectService } from "@/features/projects/services/create-project.service";
import { updateProjectService } from "@/features/projects/services/update-project.service";
import { toast } from "sonner";

interface ProjectFormProps {
  initialData?: Project;
}

export function ProjectForm({ initialData }: ProjectFormProps) {
  const router = useRouter();
  const isEditing = !!initialData;

  const [name, setName] = useState(initialData?.name || "");
  const [description, setDescription] = useState(initialData?.description || "");

  // GitHub State
  const [githubAccountId, setGithubAccountId] = useState(initialData?.github_repository?.github_account_id || "");
  const [repositoryId, setRepositoryId] = useState(initialData?.github_repository?.repository_identifier || "");
  const [defaultBranch, setDefaultBranch] = useState(initialData?.github_repository?.default_branch || "main");

  // Dokploy State
  const [dokployServerId, setDokployServerId] = useState(initialData?.dokploy_application?.dokploy_server_id || "");
  const [applicationId, setApplicationId] = useState(initialData?.dokploy_application?.application_identifier || "");

  // Monitoring Config State
  const [monitoringEnabled, setMonitoringEnabled] = useState(initialData?.monitoring_configuration?.enabled ?? true);
  const [groupingWindow, setGroupingWindow] = useState(initialData?.monitoring_configuration?.grouping_window || "PT30M");
  const [errorPatterns, setErrorPatterns] = useState<string[]>(initialData?.monitoring_configuration?.error_patterns || []);
  const [ignoredPatterns, setIgnoredPatterns] = useState<string[]>(initialData?.monitoring_configuration?.ignored_patterns || []);
  const [contextBefore, setContextBefore] = useState<number>(initialData?.monitoring_configuration?.context_before ?? 20);
  const [contextAfter, setContextAfter] = useState<number>(initialData?.monitoring_configuration?.context_after ?? 20);

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const isMissingIntegrations = !githubAccountId || !repositoryId || !dokployServerId || !applicationId;
  const isFormValid = name.trim().length > 0 && !isMissingIntegrations;

  const handleGitHubSelect = (repo: GitHubRepository) => {
    setGithubAccountId(repo.github_account_id || "");
    setRepositoryId(repo.repository_identifier || "");
    setDefaultBranch(repo.default_branch || "main");
  };

  const handleDokploySelect = (app: DokployApplication) => {
    setDokployServerId(app.dokploy_server_id || "");
    setApplicationId(app.application_identifier || "");
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isFormValid) return;

    setIsSubmitting(true);
    setError(null);

    try {
      if (isEditing) {
        const { data, error: updateError } = await updateProjectService(initialData.id, {
          name,
          description,
          github_account_id: githubAccountId,
          repository_identifier: repositoryId,
          default_branch: defaultBranch,
          dokploy_server_id: dokployServerId,
          application_identifier: applicationId,
        });

        if (updateError) throw updateError;
        if (data) {
          toast.success("Project updated successfully");
          router.push(`/projects/${data.id}`);
        }
      } else {
        const { data, error: createError } = await createProjectService({
          name,
          description,
          github_account_id: githubAccountId,
          repository_identifier: repositoryId,
          default_branch: defaultBranch,
          dokploy_server_id: dokployServerId,
          application_identifier: applicationId,
          monitoring_configuration: {
            enabled: monitoringEnabled,
            grouping_window: groupingWindow,
            error_patterns: errorPatterns,
            ignored_patterns: ignoredPatterns,
            context_before: contextBefore,
            context_after: contextAfter,
          },
        });

        if (createError) throw createError;
        if (data) {
          toast.success("Project created successfully");
          router.push(`/projects`);
        }
      }
    } catch (err: any) {
      setError(err.message || "An error occurred while saving the project.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      {isMissingIntegrations && (
        <div className={styles.warningBanner}>
          <AlertTriangle className={styles.warningIcon} size={20} />
          <p className={styles.warningText}>
            You must connect and select both a GitHub Repository and a Dokploy Application to {isEditing ? "save" : "create"} a project. Please complete the setup below.
          </p>
        </div>
      )}

      {error && (
        <div className={styles.warningBanner} style={{ borderColor: 'var(--red-500)', backgroundColor: 'rgba(239, 68, 68, 0.1)' }}>
          <AlertTriangle className={styles.warningIcon} style={{ color: 'var(--red-500)' }} size={20} />
          <p className={styles.warningText} style={{ color: 'var(--red-100)' }}>{error}</p>
        </div>
      )}

      <div className={styles.section}>
        <h3 className={styles.sectionTitle}>General Information</h3>
        <div className={styles.field}>
          <label htmlFor="projectName" className={styles.label}>Project Name</label>
          <input
            id="projectName"
            className={styles.input}
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="e.g. E-Commerce Platform"
            required
          />
        </div>
        <div className={styles.field}>
          <label htmlFor="projectDesc" className={styles.label}>Description (Optional)</label>
          <input
            id="projectDesc"
            className={styles.input}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Brief description of this project"
          />
        </div>
      </div>

      <div className={styles.section}>
        <h3 className={styles.sectionTitle}>GitHub Integration</h3>
        <GitHubRepositorySelector 
          onSelect={handleGitHubSelect} 
          selectedRepoId={repositoryId}
        />
      </div>

      <div className={styles.section}>
        <h3 className={styles.sectionTitle}>Dokploy Integration</h3>
        <DokployApplicationSelector 
          onSelect={handleDokploySelect} 
          selectedAppId={applicationId}
        />
      </div>

      <div className={styles.section}>
        <h3 className={styles.sectionTitle}>Monitoring Configuration</h3>
        <div className={styles.checkboxRow}>
          <input
            id="monitoringEnabled"
            type="checkbox"
            checked={monitoringEnabled}
            onChange={(e) => setMonitoringEnabled(e.target.checked)}
          />
          <label htmlFor="monitoringEnabled" className={styles.label} style={{ margin: 0 }}>Enable automatic monitoring and alerting</label>
        </div>
        
        <div className={styles.fieldRow}>
          <div className={styles.field}>
            <label htmlFor="groupingWindow" className={styles.label}>Grouping Window</label>
            <input
              id="groupingWindow"
              className={styles.input}
              value={groupingWindow}
              onChange={(e) => setGroupingWindow(e.target.value)}
              placeholder="e.g. PT30M"
              required
            />
          </div>
          <div className={styles.field}>
            <label htmlFor="contextBefore" className={styles.label}>Context Before (lines)</label>
            <input
              id="contextBefore"
              type="number"
              min="0"
              className={styles.input}
              value={contextBefore}
              onChange={(e) => setContextBefore(Number(e.target.value))}
              required
            />
          </div>
          <div className={styles.field}>
            <label htmlFor="contextAfter" className={styles.label}>Context After (lines)</label>
            <input
              id="contextAfter"
              type="number"
              min="0"
              className={styles.input}
              value={contextAfter}
              onChange={(e) => setContextAfter(Number(e.target.value))}
              required
            />
          </div>
        </div>

        <div className={styles.field}>
          <label htmlFor="errorPatterns" className={styles.label}>Error Patterns (Regex, one per line)</label>
          <textarea
            id="errorPatterns"
            className={styles.textarea}
            value={errorPatterns.join('\n')}
            onChange={(e) => setErrorPatterns(e.target.value.split('\n').filter(p => p.trim() !== ''))}
            placeholder="e.g. ^Error:.*"
            rows={3}
          />
        </div>

        <div className={styles.field}>
          <label htmlFor="ignoredPatterns" className={styles.label}>Ignored Patterns (Regex, one per line)</label>
          <textarea
            id="ignoredPatterns"
            className={styles.textarea}
            value={ignoredPatterns.join('\n')}
            onChange={(e) => setIgnoredPatterns(e.target.value.split('\n').filter(p => p.trim() !== ''))}
            placeholder="e.g. ^Warning:.*"
            rows={3}
          />
        </div>
      </div>

      <div className={styles.actions}>
        <Button 
          type="button" 
          variant="secondary" 
          onClick={() => router.back()}
          disabled={isSubmitting}
        >
          Cancel
        </Button>
        <Button 
          type="submit" 
          variant="primary"
          disabled={!isFormValid || isSubmitting}
        >
          {isSubmitting ? "Saving..." : (isEditing ? "Save Changes" : "Create Project")}
        </Button>
      </div>
    </form>
  );
}
