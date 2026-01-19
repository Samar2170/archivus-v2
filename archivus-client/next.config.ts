import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  images: {
    unoptimized: true,
    remotePatterns: [
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '8000',
      },
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '8001',
      },
      {
        protocol: 'http',
        hostname: 'inspiron.local',
        port: '3000',
      },
      {
        protocol: 'http',
        hostname: 'inspiron.local',
        port: '8001',
      },
      {
        protocol: 'http',
        hostname: '0.0.0.0',
        port: '8001',
      }
    ],
  },
  env: {
    NEXT_PUBLIC_BASE_URL: process.env.NEXT_PUBLIC_BASE_URL,
  },

};

export default nextConfig;
