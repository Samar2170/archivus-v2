import { apiFetch } from "../utils/fetcher";


const BASE_URL = "http://localhost:8000/"; // adjust to match your backend


export interface FileMetaData {
    ID: string;
    Name: string;
    IsDir: boolean;
    Extension: string;
    SignedUrl: string;
    Size: number;
    Path: string;
    NavigationPath: string;
}

export interface ListFileMetaData {
    ID: number;
    Name: string;
    IsImage: boolean;
    FilePath: string;
    CreatedAt: string;
    SizeInMb: number;
    UpdatedAt: string;
}

export interface getFilesByFolderResponse {
    files: FileMetaData[];
    size: number;
}

export interface filesListResponse {
    files: ListFileMetaData[];
}

export function getFilesByFolder(folder: string) {
    return apiFetch<getFilesByFolderResponse>(`${BASE_URL}files/get/?folder=${folder}`, {
        method: "GET",
    });
}

export function uploadFiles(folder: string, files: File[]) {
    const formData = new FormData();
    files.forEach(file => formData.append("file", file));
    formData.append("folder", folder);
    return apiFetch<{ success: boolean }>(`${BASE_URL}files/upload/`, {
        method: "POST",
        body: formData,
    },
    );
}

export function getFilesList(search: string) {
    return apiFetch<filesListResponse>(`${BASE_URL}files/list/?search=${search}`, {
        method: "GET",
    });
}