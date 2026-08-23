"use client";

import React, { useState, useEffect, useMemo } from "react";
import styles from "./DokploySourceSelector.module.css";
import { Badge } from "@/core/ui/primitives/Badge";
import { Search, RefreshCw, CheckCircle, AlertCircle, Server, Box, Layers, ChevronDown, ChevronRight } from "lucide-react";
import { DokployServer, listDokployServersService } from "@/features/settings/services/dokploy/list-dokploy-servers.service";
import { testDokployConnectionService } from "@/features/settings/services/dokploy/test-dokploy-connection.service";
import { DokployApplication, getDokployApplicationsService } from "@/features/settings/services/dokploy/get-dokploy-applications.service";
import { DokployCompose, getDokployComposesService } from "@/features/settings/services/dokploy/get-dokploy-composes.service";
import { DokployComposeService, getDokployComposeServicesService } from "@/features/settings/services/dokploy/get-dokploy-compose-services.service";
import type { components } from "@/core/libs/api-client";

type DokploySourceSelectorType = components["schemas"]["DokploySourceSelector"];

interface DokploySourceSelectorProps {
  onSelect: (source: DokploySourceSelectorType) => void;
  selectedSource?: DokploySourceSelectorType;
}

export const DokploySourceSelector: React.FC<DokploySourceSelectorProps> = ({ onSelect, selectedSource }) => {
  const [servers, setServers] = useState<DokployServer[]>([]);
  const [selectedServerId, setSelectedServerId] = useState<string>("");
  const [isLoadingServers, setIsLoadingServers] = useState(true);

  const [connectionStatus, setConnectionStatus] = useState<"idle" | "testing" | "success" | "error">("idle");
  const [applications, setApplications] = useState<DokployApplication[]>([]);
  const [composes, setComposes] = useState<DokployCompose[]>([]);
  const [expandedComposes, setExpandedComposes] = useState<Record<string, boolean>>({});
  const [composeServices, setComposeServices] = useState<Record<string, DokployComposeService[]>>({});
  const [isLoadingSources, setIsLoadingSources] = useState(false);
  const [loadingServices, setLoadingServices] = useState<Record<string, boolean>>({});
  
  const [searchQuery, setSearchQuery] = useState("");

  useEffect(() => {
    const fetchServers = async () => {
      setIsLoadingServers(true);
      const { data } = await listDokployServersService();
      if (data && data.length > 0) {
        setServers(data);
        setSelectedServerId(data[0].id!);
      }
      setIsLoadingServers(false);
    };
    fetchServers();
  }, []);

  useEffect(() => {
    if (!selectedServerId) {
      setApplications([]);
      setComposes([]);
      setConnectionStatus("idle");
      return;
    }
    
    const testAndFetch = async () => {
      setConnectionStatus("testing");
      setApplications([]);
      setComposes([]);
      setIsLoadingSources(true);

      try {
        const { data: testResult } = await testDokployConnectionService(selectedServerId);
        
        if (testResult?.data?.status !== "connected") {
          setConnectionStatus("error");
          setIsLoadingSources(false);
          return;
        }
        
        setConnectionStatus("success");

        const [appsRes, composesRes] = await Promise.all([
          getDokployApplicationsService(selectedServerId).catch(() => ({ data: [] })),
          getDokployComposesService(selectedServerId).catch(() => ({ data: [] }))
        ]);

        if (appsRes.data) setApplications(appsRes.data);
        if (composesRes.data) setComposes(composesRes.data);
      } catch (err) {
        setConnectionStatus("error");
      } finally {
        setIsLoadingSources(false);
      }
    };
    
    testAndFetch();
  }, [selectedServerId]);

  const handleToggleCompose = async (compose: DokployCompose) => {
    const id = compose.compose_identifier;
    const isExpanded = !!expandedComposes[id];
    
    setExpandedComposes(prev => ({ ...prev, [id]: !isExpanded }));

    if (!isExpanded && !composeServices[id]) {
      setLoadingServices(prev => ({ ...prev, [id]: true }));
      try {
        const { data } = await getDokployComposeServicesService(selectedServerId, id);
        if (data) {
          setComposeServices(prev => ({ ...prev, [id]: data }));
        }
      } catch (err) {
        console.error("Failed to load compose services", err);
      } finally {
        setLoadingServices(prev => ({ ...prev, [id]: false }));
      }
    }
  };

  const filteredApplications = useMemo(() => {
    if (!searchQuery.trim()) return applications;
    const lowerQuery = searchQuery.toLowerCase();
    return applications.filter(
      app => app.display_name?.toLowerCase().includes(lowerQuery) || 
             app.application_identifier?.toLowerCase().includes(lowerQuery) ||
             app.instance_identifier?.toLowerCase().includes(lowerQuery)
    );
  }, [applications, searchQuery]);

  const filteredComposes = useMemo(() => {
    if (!searchQuery.trim()) return composes;
    const lowerQuery = searchQuery.toLowerCase();
    return composes.filter(
      comp => comp.display_name?.toLowerCase().includes(lowerQuery) || 
             comp.compose_identifier?.toLowerCase().includes(lowerQuery) ||
             comp.instance_identifier?.toLowerCase().includes(lowerQuery)
    );
  }, [composes, searchQuery]);

  const isAppSelected = (appId: string) => {
    return selectedSource?.type === "application" && selectedSource.resource_identifier === appId;
  };

  const isServiceSelected = (composeId: string, serviceName: string) => {
    return selectedSource?.type === "compose_service" && 
           selectedSource.resource_identifier === composeId && 
           selectedSource.service_name === serviceName;
  };

  const hasResults = filteredApplications.length > 0 || filteredComposes.length > 0;

  return (
    <div className={styles.container}>
      <div className={styles.serverSection}>
        <select 
          className={styles.select}
          value={selectedServerId}
          onChange={(e) => setSelectedServerId(e.target.value)}
          disabled={isLoadingServers || servers.length === 0 || connectionStatus === "testing"}
        >
          {isLoadingServers && <option value="">Loading servers...</option>}
          {!isLoadingServers && servers.length === 0 && <option value="">No Dokploy servers found</option>}
          {servers.map(srv => (
            <option key={srv.id} value={srv.id}>{srv.name} ({srv.base_url})</option>
          ))}
        </select>
        
        {connectionStatus === "testing" && (
          <span className={styles.statusTesting}><RefreshCw size={16} className={styles.spin} /> Testing...</span>
        )}
        {connectionStatus === "success" && (
          <span className={styles.statusSuccess}><CheckCircle size={16} /> Connected</span>
        )}
        {connectionStatus === "error" && (
          <span className={styles.statusError}><AlertCircle size={16} /> Failed to connect</span>
        )}
      </div>

      <div className={styles.searchBar}>
        <Search size={16} className={styles.searchIcon} />
        <input 
          type="text" 
          placeholder="Search applications and composes..." 
          className={styles.searchInput}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          disabled={!selectedServerId || isLoadingSources || connectionStatus === "error"}
        />
      </div>

      <div className={styles.listContainer}>
        {isLoadingSources ? (
          <div className={styles.emptyState}>
            <RefreshCw size={24} className={styles.spin} />
            <p>Loading sources...</p>
          </div>
        ) : connectionStatus === "error" ? (
           <div className={styles.emptyState}>
            <AlertCircle size={48} className={styles.emptyIcon} />
            <p>Could not connect to the selected server.</p>
          </div>
        ) : !selectedServerId ? (
          <div className={styles.emptyState}>
            <Server size={48} className={styles.emptyIcon} />
            <p>Select a Dokploy server to view its sources.</p>
          </div>
        ) : !hasResults ? (
          <div className={styles.emptyState}>
            <Box size={48} className={styles.emptyIcon} />
            <p>{searchQuery ? "No sources match your search" : "No sources found on this server"}</p>
          </div>
        ) : (
          <div className={styles.listScroll}>
            
            {/* Applications */}
            {filteredApplications.length > 0 && (
              <div className={styles.sourceGroup}>
                <div className={styles.groupHeader}>Applications</div>
                {filteredApplications.map(app => (
                  <div 
                    key={app.application_identifier} 
                    className={`${styles.appItem} ${isAppSelected(app.application_identifier) ? styles.selected : ""}`}
                    onClick={() => onSelect({
                      type: "application",
                      dokploy_server_id: selectedServerId,
                      resource_identifier: app.application_identifier
                    })}
                  >
                    <Box size={18} color="var(--text-secondary)" />
                    <div className={styles.appInfo}>
                      <div className={styles.appNameRow}>
                        <span className={styles.appName}>{app.display_name}</span>
                        {app.environment && (
                          <Badge variant={app.environment === "production" ? "intel" : "neutral"}>
                            {app.environment}
                          </Badge>
                        )}
                        {app.status && (
                          <Badge variant={app.status === "running" ? "success" : app.status === "stopped" ? "neutral" : "error"}>
                            {app.status}
                          </Badge>
                        )}
                      </div>
                      <span className={styles.appId}>{app.instance_identifier} ({app.application_identifier})</span>
                    </div>
                    {isAppSelected(app.application_identifier) && <CheckCircle size={18} color="var(--status-success)" />}
                  </div>
                ))}
              </div>
            )}

            {/* Composes */}
            {filteredComposes.length > 0 && (
              <div className={styles.sourceGroup}>
                <div className={styles.groupHeader}>Composes</div>
                {filteredComposes.map(comp => {
                  const isExpanded = !!expandedComposes[comp.compose_identifier];
                  const services = composeServices[comp.compose_identifier] || [];
                  const isLoadingSvc = loadingServices[comp.compose_identifier];

                  const isSelectedCompose = selectedSource?.type === "compose_service" && selectedSource.resource_identifier === comp.compose_identifier;

                  return (
                    <div key={comp.compose_identifier} className={`${styles.composeContainer} ${isSelectedCompose ? styles.selected : ""}`}>
                      <div 
                        className={styles.composeHeader}
                        onClick={() => handleToggleCompose(comp)}
                      >
                        {isExpanded ? <ChevronDown size={18} /> : <ChevronRight size={18} />}
                        <Layers size={18} color="var(--text-secondary)" />
                        <div className={styles.appInfo}>
                          <div className={styles.appNameRow}>
                            <span className={styles.appName}>{comp.display_name}</span>
                            {comp.environment_identifier && (
                              <Badge variant="neutral">{comp.environment_identifier}</Badge>
                            )}
                            {comp.status && (
                              <Badge variant={comp.status === "running" ? "success" : comp.status === "stopped" ? "neutral" : "error"}>
                                {comp.status}
                              </Badge>
                            )}
                          </div>
                          <span className={styles.appId}>{comp.instance_identifier} ({comp.compose_identifier})</span>
                        </div>
                      </div>

                      {isExpanded && (
                        <div className={styles.composeServices}>
                          {isLoadingSvc ? (
                            <div className={styles.serviceLoading}>
                              <RefreshCw size={14} className={styles.spin} /> Loading services...
                            </div>
                          ) : services.length === 0 ? (
                            <div className={styles.serviceEmpty}>No services found in this compose.</div>
                          ) : (
                            services.map(svc => (
                              <div
                                key={svc.service_name}
                                className={`${styles.serviceItem} ${isServiceSelected(comp.compose_identifier, svc.service_name) ? styles.selected : ""}`}
                                onClick={() => onSelect({
                                  type: "compose_service",
                                  dokploy_server_id: selectedServerId,
                                  resource_identifier: comp.compose_identifier,
                                  service_name: svc.service_name
                                })}
                              >
                                <div className={styles.serviceName}>{svc.service_name}</div>
                                {isServiceSelected(comp.compose_identifier, svc.service_name) && <CheckCircle size={16} color="var(--status-success)" />}
                              </div>
                            ))
                          )}
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            )}
            
          </div>
        )}
      </div>
    </div>
  );
};
