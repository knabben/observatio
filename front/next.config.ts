import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "export",
  distDir: './output',
  images: {
    unoptimized: true
  },
  experimental: {
    optimizePackageImports: ['@mantine/core', '@mantine/hooks']
  },
};

export default nextConfig;
