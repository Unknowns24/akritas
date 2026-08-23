"use client";

import React, { useState, useEffect, useMemo } from "react";
import styles from "./DokployApplicationSelector.module.css";
import { Badge } from "@/core/ui/primitives/Badge";
import { Search, RefreshCw, CheckCircle, AlertCircle, Server, Box } from "lucide-react";
import { DokployServer, listDokployServersService } from "@/features/settings/services/dokploy/list-dokploy-servers.service";
import { testDokployConnectionService } from "@/features/settings/services/dokploy/test-dokploy-connection.service";
import { DokployApplication, getDokployApplicationsService } from "@/features/settings/services/dokploy/get-dokploy-applications.service";

interface DokployApplicationSelectorProps {
  onSelect: (app: DokployApplication) => void;
  selectedAppId?: string;
}

export const DokployApplicationSelector: React.FC<DokployApplicationSelectorProps> = ({ onSelect, selectedAppId }) => {
  const [servers, setServers] = useState<DokployServer[]>([]);
  const [selectedServerId, setSelectedServerId] = useState<string>("");
  const [isLoadingServers, setIsLoadingServers] = useState(true);

  const [connectionStatus, setConnectionStatus] = useState<"idle" | "testing" | "success" | "error">("idle");
  const [applications, setApplications] = useState<DokployApplication[]>([]);
  const [isLoadingApps, setIsLoadingApps] = useState(false);
  
  const [searchQuery, setSearchQuery] = useState("");

  // 1. Fetch connected servers on mount
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

  // 2. Fetch apps automatically when server changes
  useEffect(() => {
    if (!selectedServerId) {
      setApplications([]);
      setConnectionStatus("idle");
      return;
    }
    
    const testAndFetch = async () => {
      // Step A: Test connection
      setConnectionStatus("testing");
      setApplications([]);
      setIsLoadingApps(true);

      const { data: testResult } = await testDokployConnectionService(selectedServerId);
      
      if (testResult?.data?.status !== "connected") {
        setConnectionStatus("error");
        setIsLoadingApps(false);
        return;
      }
      
      setConnectionStatus("success");

      // Step B: Fetch apps
      const { data: appsData } = await getDokployApplicationsService(selectedServerId);
      if (appsData) {
        setApplications(appsData);
      }
      setIsLoadingApps(false);
    };
    
    testAndFetch();
  }, [selectedServerId]);

  const filteredApplications = useMemo(() => {
    if (!searchQuery.trim()) return applications;
    const lowerQuery = searchQuery.toLowerCase();
    return applications.filter(
      app => app.display_name?.toLowerCase().includes(lowerQuery) || 
             app.application_identifier?.toLowerCase().includes(lowerQuery) ||
             app.instance_identifier?.toLowerCase().includes(lowerQuery)
    );
  }, [applications, searchQuery]);

  return (
    <div className={styles.container}>
      {/* Server Selection and Connection Status */}
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

      {/* Search Bar */}
      <div className={styles.searchBar}>
        <Search size={16} className={styles.searchIcon} />
        <input 
          type="text" 
          placeholder="Search applications..." 
          className={styles.searchInput}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          disabled={!selectedServerId || isLoadingApps || connectionStatus === "error"}
        />
      </div>

      {/* Application List */}
      <div className={styles.listContainer}>
        {isLoadingApps ? (
          <div className={styles.emptyState}>
            <RefreshCw size={24} className={styles.spin} />
            <p>Loading applications...</p>
          </div>
        ) : connectionStatus === "error" ? (
           <div className={styles.emptyState}>
            <AlertCircle size={48} className={styles.emptyIcon} />
            <p>Could not connect to the selected server.</p>
          </div>
        ) : !selectedServerId ? (
          <div className={styles.emptyState}>
            <Server size={48} className={styles.emptyIcon} />
            <p>Select a Dokploy server to view its applications.</p>
          </div>
        ) : filteredApplications.length === 0 ? (
          <div className={styles.emptyState}>
            <Box size={48} className={styles.emptyIcon} />
            <p>{searchQuery ? "No applications match your search" : "No applications found on this server"}</p>
          </div>
        ) : (
          <div className={styles.listScroll}>
            {filteredApplications.map(app => {
              const isSelected = selectedAppId === app.application_identifier;
              return (
                <div 
                  key={app.application_identifier} 
                  className={`${styles.appItem} ${isSelected ? styles.selected : ""}`}
                  onClick={() => onSelect(app)}
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
                  {isSelected && <CheckCircle size={18} color="var(--status-success)" />}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
};
