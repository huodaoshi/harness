/** @type {import('next').NextConfig} */
const apiProxyTarget =
  process.env.API_PROXY_TARGET ?? "http://127.0.0.1:8080";

const nextConfig = {
  eslint: {
    ignoreDuringBuilds: true,
  },
  typescript: {
    ignoreBuildErrors: true,
  },
  webpack(config) {
    config.module.rules.push({
      test: /\.svg$/,
      use: ["@svgr/webpack"],
    });
    config.resolve.fallback = {
      ...config.resolve.fallback,
      child_process: false,
    };
    return config;
  },
  async rewrites() {
    if (process.env.NEXT_PUBLIC_API_BASE_URL) {
      return [];
    }
    return [
      {
        source: "/api/:path*",
        destination: `${apiProxyTarget}/api/:path*`,
      },
      {
        source: "/v1/:path*",
        destination: `${apiProxyTarget}/v1/:path*`,
      },
    ];
  },
};

export default nextConfig;
