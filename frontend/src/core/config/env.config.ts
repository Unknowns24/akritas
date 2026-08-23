import { env as getEnv } from "next-runtime-env";

/**
 * Runtime environment configuration for Akritas frontend.
 * Centralizes access to environment variables.
 */
const isServer = typeof window === "undefined";
const configuredUrl = getEnv("NEXT_PUBLIC_API_URL") || "/api/v1";
const rawApiUrl = configuredUrl.replace(/^NEXT_PUBLIC_API_URL=/, "").trim();
const apiUrl = isServer && rawApiUrl.startsWith("/") 
  ? `http://localhost:${process.env.PORT || 3000}${rawApiUrl}`
  : rawApiUrl;

export const env = {
  apiUrl,
  isProduction: process.env.NODE_ENV === "production",
  isDevelopment: process.env.NODE_ENV === "development",
} as const;
