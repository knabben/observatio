import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "export",
  distDir: './output',
  images: { unoptimized: true }
};

export default nextConfig;
