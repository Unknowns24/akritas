import { env as getEnv } from "next-runtime-env";

/**
 * Runtime environment configuration for Akritas frontend.
 * Centralizes access to environment variables.
 */
export const env = {
  apiUrl: getEnv("NEXT_PUBLIC_API_URL") || "/api/v1",
  isProduction: getEnv("NODE_ENV") === "production",
  isDevelopment: getEnv("NODE_ENV") === "development",
} as const;
