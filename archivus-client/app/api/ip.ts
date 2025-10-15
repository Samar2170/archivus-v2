import os from "os";
import type { NextApiRequest, NextApiResponse } from "next";

export default function handler(req: NextApiRequest, res: NextApiResponse) {
  const interfaces = os.networkInterfaces();
  let ip = "127.0.0.1";
  for (const name of Object.keys(interfaces)) {
    for (const net of interfaces[name] || []) {
      if (net.family === "IPv4" && !net.internal) ip = net.address;
    }
  }
  res.json({ ip });
}