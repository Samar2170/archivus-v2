import { apiFetch } from "../utils/fetcher";
const BASE_URL = "http://localhost:8000/"; // adjust to match your backend

export function addFolder(folder: string) {
    return apiFetch<{ success: boolean }>(`${BASE_URL}folder/add/`, {
        method: "POST",
        body: JSON.stringify({ folder }),
    });
}

export interface FolderResponse {
    Name: string;
    Path: string;
}

export function listFolders() {
    return apiFetch<FolderResponse[]>(`${BASE_URL}folder/list/`, {
        method: "GET",
    });
}