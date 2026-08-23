import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactCompiler: true,
  async rewrites() {
    const apiProxyTarget =
      process.env.AKRITAS_API_PROXY_TARGET || "http://localhost:8080";

    return [
      {
        source: "/api/v1/:path*",
        destination: `${apiProxyTarget.replace(/\/$/, "")}/api/v1/:path*`,
      },
    ];
  },
};

export default nextConfig;
