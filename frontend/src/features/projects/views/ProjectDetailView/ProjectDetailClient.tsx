"use client";

import React from "react";
import Link from "next/link";
import styles from "./ProjectDetailClient.module.css";
import { Badge } from "@/core/ui/primitives/Badge";
import { Button } from "@/core/ui/primitives/Button";
import { useState, useEffect } from "react";
import { updateMonitoringConfigService, getProjectService } from "../../services";
import { Activity, Clock, Code2, Database, ShieldAlert, GitBranch, Box, Settings, ExternalLink, Loader2 } from "lucide-react";
import type { components } from "@/core/libs/api-client";
import { GithubIcon, DokployIcon } from "@/core/ui/icons";

type ProjectResponse = components["schemas"]["ProjectResponse"];

interface ProjectDetailClientProps {
  initialProject: ProjectResponse | undefined;
  projectId?: string;
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

export function ProjectDetailClient({ initialProject, projectId }: ProjectDetailClientProps) {
  const [project, setProject] = useState<ProjectResponse["data"] | undefined>(initialProject?.data);
  const [isLoading, setIsLoading] = useState(!initialProject?.data && !!projectId);
  const [isUpdating, setIsUpdating] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  useEffect(() => {
    if (!project && projectId) {
      setIsLoading(true);
      getProjectService(projectId)
        .then((res) => {
          if (res?.data) {
            setProject(res.data);
          }
        })
        .catch((err) => {
          console.error("Failed to load project client-side:", err);
        })
        .finally(() => {
          setIsLoading(false);
        });
    }
  }, [projectId]);

  if (isLoading) {
    return (
      <div style={{ display: "flex", justifyContent: "center", alignItems: "center", minHeight: "300px" }}>
        <Loader2 style={{ animation: "spin 1s linear infinite" }} size={32} />
      </div>
    );
  }

  if (!project) {
    return (
      <div className={styles.emptyState}>
        <ShieldAlert size={48} />
        <h2>Project not found</h2>
        <p>The project you are looking for does not exist or you don't have access.</p>
      </div>
    );
  }

  const isHealthy = project.health_status === "healthy";
  const repo = project.github_repository;
  const source = project.dokploy_source;
  const mon = project.monitoring_configuration;
  const rules = project.built_in_detection_rules || [];

  const handleToggleMonitoring = async () => {
    if (!mon) return;
    
    setIsUpdating(true);
    setErrorMsg(null);
    
    const newConfig = {
      ...mon,
      enabled: !mon.enabled
    };
    
    const { data, error } = await updateMonitoringConfigService(project.id, newConfig);
    
    if (error) {
      setErrorMsg("Failed to update monitoring configuration.");
      setIsUpdating(false);
      return;
    }
    
    if (data) {
      setProject({
        ...project,
        monitoring_configuration: data,
        monitoring_status: data.enabled ? "starting" : "disabled" // Optimistic UI update for status
      });
    }
    
    setIsUpdating(false);
  };

  const isIngestionWarning = project.monitoring_status === "error" || project.monitoring_status === "degraded" || (!project.last_observed_at && project.monitoring_status !== "disabled");

  return (
    <div className={styles.container}>
      {/* HEADER SECTION */}
      <div className={styles.headerCard}>
        <div className={styles.headerInfo}>
          <div className={styles.titleRow}>
            <h1 className={styles.title}>{project.name}</h1>
            <Link href={`/projects/${project.id}/settings`}>
              <Button variant="secondary" style={{ padding: "8px 12px" }}>
                <Settings size={16} />
                <span style={{ marginLeft: "8px" }}>Settings</span>
              </Button>
            </Link>
          </div>
          {project.description && (
            <p className={styles.description}>{project.description}</p>
          )}
          <div className={styles.headerMeta}>
            <div className={styles.metaItem}>
              <Clock size={14} />
              Created {new Date(project.created_at).toLocaleDateString()}
            </div>
            {project.last_observed_at ? (
              <div className={styles.metaItem}>
                <Activity size={14} />
                Last log: {formatRelativeTime(project.last_observed_at)} ({new Date(project.last_observed_at).toLocaleString()})
              </div>
            ) : (
              <div className={styles.metaItem} style={{ color: "var(--status-warning)" }}>
                <Activity size={14} />
                Awaiting first log
              </div>
            )}
          </div>
        </div>
        <div className={styles.badges}>
          <Badge variant={isHealthy ? "success" : "warning"}>
            HEALTH: {(project.health_status || "UNKNOWN").toUpperCase()}
          </Badge>
          <Badge variant={
            project.monitoring_status === "monitoring" ? "intel" :
            project.monitoring_status === "error" ? "error" :
            project.monitoring_status === "degraded" ? "warning" : "neutral"
          }>
            MONITORING: {(project.monitoring_status || "UNKNOWN").toUpperCase()}
          </Badge>
        </div>
      </div>

      {isIngestionWarning && (
        <div className={`${styles.ingestionBanner} ${project.monitoring_status === 'error' ? styles.bannerError : styles.bannerWarning}`}>
          <Activity size={18} />
          <div className={styles.bannerText}>
            <strong>Ingestion Status: </strong>
            {project.monitoring_status === "error" ? "Monitoring connection failed. No logs are being ingested." : 
             !project.last_observed_at ? "Waiting for logs to arrive from Dokploy. Ensure your container is running." : 
             "Monitoring is degraded. Some logs might be delayed or dropped."}
          </div>
        </div>
      )}

      <div className={styles.grid}>
        {/* GITHUB CONTEXT */}
        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <GithubIcon size={20} className={styles.cardIcon} />
            <h3 className={styles.cardTitle}>GitHub Integration</h3>
          </div>
          {repo ? (
            <div className={styles.dataList}>
              <div className={styles.dataRow}>
                <span className={styles.dataLabel}>Repository</span>
                <span className={styles.dataValue}>
                  {repo.full_name}
                  {repo.html_url && (
                    <a href={repo.html_url} target="_blank" rel="noreferrer" className={styles.link}>
                      <ExternalLink size={14} />
                    </a>
                  )}
                </span>
              </div>
              <div className={styles.dataRow}>
                <span className={styles.dataLabel}>Default Branch</span>
                <span className={styles.dataValue}>
                  <GitBranch size={14} className={styles.cardIcon} />
                  {repo.default_branch}
                </span>
              </div>
              <div className={styles.dataRow}>
                <span className={styles.dataLabel}>Visibility</span>
                <span className={styles.dataValue}>
                  <Badge variant="neutral">{repo.private ? "Private" : "Public"}</Badge>
                </span>
              </div>
            </div>
          ) : (
            <div className={styles.emptyState}>
              <Code2 size={24} />
              <p>No repository connected.</p>
            </div>
          )}
        </div>

        {/* DOKPLOY CONTEXT */}
        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <DokployIcon size={20} className={styles.cardIcon} />
            <h3 className={styles.cardTitle}>Dokploy Integration</h3>
          </div>
          {source ? (
            <div className={styles.dataList}>
              <div className={styles.dataRow}>
                <span className={styles.dataLabel}>{source.type === "application" ? "Application" : "Compose Service"}</span>
                <span className={styles.dataValue}>
                  <Box size={14} className={styles.cardIcon} />
                  {source.display_name}
                  {source.type === "compose_service" && <span style={{ marginLeft: 4, color: "var(--text-secondary)" }}>({source.service_name})</span>}
                </span>
              </div>
              <div className={styles.dataRow}>
                <span className={styles.dataLabel}>Environment</span>
                <span className={styles.dataValue}>
                  <Badge variant={source.type === "application" && source.environment === "production" ? "intel" : "neutral"}>
                    {(source as any).environment || "None"}
                  </Badge>
                </span>
              </div>
              <div className={styles.dataRow}>
                <span className={styles.dataLabel}>Status</span>
                <span className={styles.dataValue}>
                  {source.status && (
                    <Badge variant={source.status === "running" ? "success" : source.status === "stopped" ? "neutral" : "error"}>
                      {source.status}
                    </Badge>
                  )}
                </span>
              </div>
            </div>
          ) : (
            <div className={styles.emptyState}>
              <DokployIcon size={24} />
              <p>No source connected.</p>
            </div>
          )}
        </div>

        {/* MONITORING CONFIGURATION */}
        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <div className={styles.cardHeaderLeft}>
              <Activity size={20} className={styles.cardIcon} />
              <h3 className={styles.cardTitle}>Monitoring Configuration</h3>
            </div>
            {mon && (
              <Button 
                variant={mon.enabled ? "secondary" : "primary"} 
                size="sm" 
                onClick={handleToggleMonitoring}
                disabled={isUpdating}
              >
                {isUpdating ? "Updating..." : mon.enabled ? "Disable" : "Enable"}
              </Button>
            )}
          </div>
          {errorMsg && (
            <div className={styles.errorMessage}>
              {errorMsg}
            </div>
          )}
          {mon ? (
            <div className={styles.dataList}>
              <div className={styles.dataRow}>
                <span className={styles.dataLabel}>Status</span>
                <span className={styles.dataValue}>
                  {mon.enabled ? <Badge variant="success">Enabled</Badge> : <Badge variant="neutral">Disabled</Badge>}
                </span>
              </div>
              <div className={styles.dataRow}>
                <span className={styles.dataLabel}>Grouping Window</span>
                <span className={styles.dataValue}>{mon.grouping_window}</span>
              </div>
              <div className={styles.dataRow}>
                <span className={styles.dataLabel}>Context Window</span>
                <span className={styles.dataValue}>
                  B: {mon.context_before} / A: {mon.context_after}
                </span>
              </div>
              <div className={styles.dataRow}>
                <span className={styles.dataLabel}>Error Patterns</span>
                <span className={styles.dataValue}>{mon.error_patterns.length} configured</span>
              </div>
            </div>
          ) : (
            <div className={styles.emptyState}>
              <Settings size={24} />
              <p>Monitoring not configured.</p>
            </div>
          )}
        </div>

        {/* DETECTION RULES */}
        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <ShieldAlert size={20} className={styles.cardIcon} />
            <h3 className={styles.cardTitle}>Detection Rules</h3>
          </div>
          {rules.length > 0 ? (
            <div className={styles.ruleList}>
              {rules.map((rule, idx) => (
                <div key={idx} className={styles.ruleItem}>
                  <div className={styles.ruleInfo}>
                    <span className={styles.ruleName}>{rule.display_name}</span>
                    <span className={styles.ruleDesc}>Code: {rule.code}</span>
                  </div>
                  <Badge 
                    variant={rule.enabled ? "success" : "neutral"} 
                    className={!rule.enabled ? styles.badgeDisabled : undefined}
                  >
                    {rule.enabled ? "Active" : "Inactive"}
                  </Badge>
                </div>
              ))}
            </div>
          ) : (
            <div className={styles.emptyState}>
              <Database size={24} />
              <p>No detection rules configured.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
