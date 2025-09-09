import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  images: {
    remotePatterns: [
      {
        protocol: 'http',
        hostname: '192.168.1.7',
        port: '8080',
        pathname: '/storage/**',
      },
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '8000',
        
      },
      {
        protocol: 'https',
        hostname: 'backend.barebasics.shop',
        pathname:'/api/**'
      },
      {
        protocol: 'https',
        hostname: 'assets.barebasics.shop',
      },
    ],
  },
  env: {
    NEXT_PUBLIC_BASE_URL: process.env.NEXT_PUBLIC_BASE_URL,
  },

};

export default nextConfig;
