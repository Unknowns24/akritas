"use client";

import React, { useState, useEffect, useMemo } from "react";
import styles from "./GitHubRepositorySelector.module.css";
import { Button } from "@/core/ui/primitives/Button";
import { Badge } from "@/core/ui/primitives/Badge";
import { Search, RefreshCw, CheckCircle, AlertCircle, Book } from "lucide-react";
import { GitHubAccount, getGitHubAccountsService } from "@/features/settings/services/github/get-github-accounts.service";
import { testGitHubConnectionService } from "@/features/settings/services/github/test-github-connection.service";
import { GitHubRepository, getGitHubRepositoriesService } from "@/features/settings/services/github/get-github-repositories.service";

interface GitHubRepositorySelectorProps {
  onSelect: (repo: GitHubRepository) => void;
  selectedRepoId?: string;
}

const GithubIcon = ({ size = 20, className = "" }: { size?: number; className?: string }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4" />
    <path d="M9 18c-4.51 2-5-2-7-2" />
  </svg>
);

export const GitHubRepositorySelector: React.FC<GitHubRepositorySelectorProps> = ({ onSelect, selectedRepoId }) => {
  const [accounts, setAccounts] = useState<GitHubAccount[]>([]);
  const [selectedAccountId, setSelectedAccountId] = useState<string>("");
  const [isLoadingAccounts, setIsLoadingAccounts] = useState(true);

  const [connectionStatus, setConnectionStatus] = useState<"idle" | "testing" | "success" | "error">("idle");
  const [repositories, setRepositories] = useState<GitHubRepository[]>([]);
  const [isLoadingRepos, setIsLoadingRepos] = useState(false);
  
  const [searchQuery, setSearchQuery] = useState("");

  // 1. Fetch connected accounts on mount
  useEffect(() => {
    const fetchAccounts = async () => {
      setIsLoadingAccounts(true);
      const { data } = await getGitHubAccountsService();
      if (data && data.length > 0) {
        setAccounts(data);
        setSelectedAccountId(data[0].id!);
      }
      setIsLoadingAccounts(false);
    };
    fetchAccounts();
  }, []);

  // 2. Fetch repos when account changes or connection succeeds
  useEffect(() => {
    if (!selectedAccountId) {
      setRepositories([]);
      return;
    }
    
    // reset status when changing account
    setConnectionStatus("idle");
    
    const fetchRepos = async () => {
      setIsLoadingRepos(true);
      const { data } = await getGitHubRepositoriesService(selectedAccountId);
      if (data) {
        setRepositories(data);
      }
      setIsLoadingRepos(false);
    };
    
    fetchRepos();
  }, [selectedAccountId]);

  const handleTestConnection = async () => {
    if (!selectedAccountId) return;
    setConnectionStatus("testing");
    const { success } = await testGitHubConnectionService(selectedAccountId);
    setConnectionStatus(success ? "success" : "error");
    
    // Optionally refetch repos if connection was successful
    if (success) {
      setIsLoadingRepos(true);
      const { data } = await getGitHubRepositoriesService(selectedAccountId);
      if (data) {
        setRepositories(data);
      }
      setIsLoadingRepos(false);
    }
  };

  const filteredRepositories = useMemo(() => {
    if (!searchQuery.trim()) return repositories;
    const lowerQuery = searchQuery.toLowerCase();
    return repositories.filter(
      repo => repo.full_name?.toLowerCase().includes(lowerQuery) || repo.name?.toLowerCase().includes(lowerQuery)
    );
  }, [repositories, searchQuery]);

  return (
    <div className={styles.container}>
      {/* Account Selection and Test Connection */}
      <div className={styles.accountSection}>
        <select 
          className={styles.select}
          value={selectedAccountId}
          onChange={(e) => setSelectedAccountId(e.target.value)}
          disabled={isLoadingAccounts || accounts.length === 0}
        >
          {isLoadingAccounts && <option value="">Loading accounts...</option>}
          {!isLoadingAccounts && accounts.length === 0 && <option value="">No GitHub accounts found</option>}
          {accounts.map(acc => (
            <option key={acc.id} value={acc.id}>{acc.display_name} ({acc.account_identifier})</option>
          ))}
        </select>
        
        <Button 
          type="button"
          variant="secondary" 
          onClick={handleTestConnection}
          disabled={!selectedAccountId || connectionStatus === "testing"}
        >
          {connectionStatus === "testing" ? "Testing..." : "Test Connection"}
        </Button>

        {connectionStatus === "testing" && (
          <span className={styles.statusTesting}><RefreshCw size={16} className={styles.spin} /></span>
        )}
        {connectionStatus === "success" && (
          <span className={styles.statusSuccess}><CheckCircle size={16} /> OK</span>
        )}
        {connectionStatus === "error" && (
          <span className={styles.statusError}><AlertCircle size={16} /> Failed</span>
        )}
      </div>

      {/* Search Bar */}
      <div className={styles.searchBar}>
        <Search size={16} className={styles.searchIcon} />
        <input 
          type="text" 
          placeholder="Search repositories..." 
          className={styles.searchInput}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          disabled={!selectedAccountId || isLoadingRepos}
        />
      </div>

      {/* Repository List */}
      <div className={styles.listContainer}>
        {isLoadingRepos ? (
          <div className={styles.emptyState}>
            <RefreshCw size={24} className={styles.spin} />
            <p>Loading repositories...</p>
          </div>
        ) : !selectedAccountId ? (
          <div className={styles.emptyState}>
            <GithubIcon size={48} className={styles.emptyIcon} />
            <p>Select a GitHub account to view repositories</p>
          </div>
        ) : filteredRepositories.length === 0 ? (
          <div className={styles.emptyState}>
            <Book size={48} className={styles.emptyIcon} />
            <p>{searchQuery ? "No repositories match your search" : "No repositories found for this account"}</p>
          </div>
        ) : (
          <div className={styles.listScroll}>
            {filteredRepositories.map(repo => {
              const isSelected = selectedRepoId === repo.repository_identifier;
              return (
                <div 
                  key={repo.repository_identifier} 
                  className={`${styles.repoItem} ${isSelected ? styles.selected : ""}`}
                  onClick={() => onSelect(repo)}
                >
                  <Book size={18} color="var(--text-secondary)" />
                  <div className={styles.repoInfo}>
                    <div className={styles.repoNameRow}>
                      <span className={styles.repoName}>{repo.full_name}</span>
                      {repo.private && <Badge variant="neutral">Private</Badge>}
                    </div>
                    {repo.html_url && (
                      <span className={styles.repoUrl}>{repo.html_url}</span>
                    )}
                  </div>
                  {isSelected && <CheckCircle size={18} color="var(--success)" />}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
};
