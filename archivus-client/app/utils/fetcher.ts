import { useAuthStore } from "../store/auth";

// src/utils/fetcher.ts
export async function apiFetch<T>(
    url: string,
    options: RequestInit = {},
  ): Promise<T> {

    const token = useAuthStore.getState().token;
    const headers: Record<string,string> = {
      // Authorization: token ? `Bearer ${token}` : "",
      // ...(options.headers || {}),
      ...(options.headers as Record<string, string>),
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    }
    if (!(options.body instanceof FormData)) {
      headers["Content-Type"] = "application/json";
    }


    const res = await fetch(url, {
      ...options,
      headers,
    });
  
    if (!res.ok) {
      const err = await res.text();
      throw new Error(err || res.statusText);
    }
  
    return res.json() as Promise<T>;
  }
  