/**
 * Application Navigation Routes (Frontend UI Pages)
 * Aligned with docs/frontend-architecture.md (Section 28 & product routes)
 */
export const APP_ROUTES = {
  // Public Authentication Routes
  AUTH: {
    SETUP: "/setup",
    LOGIN: "/login",
    RECOVERY: "/recovery",
  },
  // Private Product Routes
  OVERVIEW: "/",
  PROJECTS: {
    LIST: "/projects",
    NEW: "/projects/new",
    DETAIL: (id: string) => `/projects/${id}`,
    MONITORING: (id: string) => `/projects/${id}/monitoring`,
  },
  INCIDENTS: {
    LIST: "/incidents",
    DETAIL: (id: string) => `/incidents/${id}`,
  },
  SETTINGS: {
    ROOT: "/settings",
    INTEGRATIONS: "/settings/integrations",
    AUTOMATION: "/settings/automation",
  },
} as const;

/**
 * Backend API Endpoints
 * 1:1 match with docs/openapi.yaml
 */
export const API_ENDPOINTS = {
  SYSTEM: {
    HEALTH: "/health",
    READINESS: "/readiness",
    STATUS: "/system/status",
    DIAGNOSTICS: "/system/diagnostics",
    OVERVIEW: "/overview",
    ACTIVITY: "/activity",
  },
  AUTH: {
    SETUP_STATUS: "/auth/setup-status",
    SETUP: "/auth/setup",
    SETUP_VERIFY: "/auth/setup/verify",
    LOGIN: "/auth/login",
    SESSION: "/auth/session",
    RECOVERY: "/auth/recovery",
    RECOVERY_VERIFY: "/auth/recovery/verify",
  },
  PROJECTS: {
    BASE: "/projects",
    BY_ID: (id: string) => `/projects/${id}`,
    MONITORING_CONFIG: (id: string) => `/projects/${id}/monitoring-configuration`,
  },
  INCIDENTS: {
    BASE: "/incidents",
    BY_ID: (id: string) => `/incidents/${id}`,
    LOG_EVENTS: (id: string) => `/incidents/${id}/log-events`,
    TIMELINE: (id: string) => `/incidents/${id}/timeline`,
    INVESTIGATIONS: (id: string) => `/incidents/${id}/investigations`,
    REMEDIATION: (id: string) => `/incidents/${id}/remediation`,
  },
  INVESTIGATIONS: {
    BY_ID: (id: string) => `/investigations/${id}`,
    EVIDENCE: (id: string) => `/investigations/${id}/evidence`,
  },
  REMEDIATIONS: {
    BY_ID: (id: string) => `/remediations/${id}`,
    VALIDATION_RESULTS: (id: string) => `/remediations/${id}/validation-results`,
    PULL_REQUEST: (id: string) => `/remediations/${id}/pull-request`,
  },
  OPERATIONS: {
    BY_ID: (id: string) => `/operations/${id}`,
  },
  INTEGRATIONS: {
    GITHUB: {
      ACCOUNTS: "/integrations/github/accounts",
      BY_ID: (id: string) => `/integrations/github/accounts/${id}`,
      TEST: (id: string) => `/integrations/github/accounts/${id}/connection-test`,
      REPOSITORIES: (id: string) => `/integrations/github/accounts/${id}/repositories`,
      REGISTRATIONS: "/integrations/github/app-manifest/registrations",
    },
    DOKPLOY: {
      SERVERS: "/integrations/dokploy/servers",
      BY_ID: (id: string) => `/integrations/dokploy/servers/${id}`,
      TEST: (id: string) => `/integrations/dokploy/servers/${id}/connection-test`,
      APPLICATIONS: (id: string) => `/integrations/dokploy/servers/${id}/applications`,
    },
    QVAC: {
      CONFIG: "/integrations/qvac/configuration",
      TEST: "/integrations/qvac/connection-test",
      STATUS: "/integrations/qvac/status",
    },
  },
  SETTINGS: {
    AUTOMATION: "/settings/automation",
  },
} as const;
