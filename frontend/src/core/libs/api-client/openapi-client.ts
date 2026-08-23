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
  cache: "no-store",
});

// Middleware to automatically pass cookies during Next.js Server-Side Rendering
api.use({
  async onRequest({ request }) {
    if (typeof window === "undefined") {
      try {
        const { cookies } = require("next/headers");
        const cookieStore = await cookies();
        const cookieString = cookieStore
          .getAll()
          .map((c: any) => `${c.name}=${c.value}`)
          .join("; ");
        if (cookieString) {
          request.headers.set("cookie", cookieString);
        }
      } catch (e) {
        // Ignore errors if next/headers is unavailable or called outside Request scope
      }
    }
    return request;
  },
});
