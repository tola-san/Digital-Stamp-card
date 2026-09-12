import type { NextConfig } from "next";

const apiProxyTarget = (
  process.env.API_PROXY_TARGET ||
  "https://digital-stamp-card.onrender.com/api"
).replace(/\/+$/, "");

const nextConfig: NextConfig = {
  /* Keep API requests same-origin so browser privacy rules do not block sessions. */
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${apiProxyTarget}/:path*`,
      },
    ];
  },
};

export default nextConfig;
