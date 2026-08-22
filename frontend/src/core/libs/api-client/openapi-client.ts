import createClient from "openapi-fetch";
import type { paths } from "./api.types";
import { env } from "@/core/config/env.config";

/**
 * Fully typed OpenAPI fetch client instance.
 * Automatically injects the correct Base URL and ensures credentials (cookies) are sent.
 */
export const api = createClient<paths>({
  baseUrl: env.apiUrl,
  credentials: "include",
});
