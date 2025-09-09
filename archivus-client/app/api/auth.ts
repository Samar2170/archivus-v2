import { apiFetch } from "../utils/fetcher";

const BASE_URL = "http://localhost:8000/"; // adjust to match your backend

export interface SigninResponse {
  token: string;
}

export function signin(username: string, password: string, pin: string) {
  return apiFetch<SigninResponse>(`${BASE_URL}login/`, {
    method: "POST",
    body: JSON.stringify({ username, password, pin }),
  });
}

