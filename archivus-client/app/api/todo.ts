import { apiFetch } from "../utils/fetcher";
import { baseUrl } from "../data/constants";

export interface Project {
    id: number;
    title: string;
    description: string;
    projectId: number;
    userId: string;
}

export interface Todo {
    id: number;
    title: string;
    description: string;
    status: number;
    priority: number;
    projectId: number;
    userId: string;
}

export function getProjects() {
    return apiFetch<Project[]>(`${baseUrl}tempora/projects`, {
        method: "GET",
    });
}

export function createProject(title: string, description: string) {
    const body = { title, description };
    return apiFetch<Project>(`${baseUrl}tempora/projects`, {
        method: "POST",
        body: JSON.stringify(body),
        headers: {
            "Content-Type": "application/json",
        },
    });
}

export function getTodos(projectId?: number) {
    var url = `${baseUrl}tempora/todos`;
    if (projectId) {
        url += `?projectId=${projectId}`;
    }
    return apiFetch<Todo[]>(url, {
        method: "GET",
    });
}

export function createTodo(title: string, description: string, projectId?: number) {
    const body = [{ title, description, projectId }];
    return apiFetch<Todo>(`${baseUrl}tempora/todos`, {
        method: "POST",
        body: JSON.stringify(body),
        headers: {
            "Content-Type": "application/json",
        },
    });
}

export function updateTodo(id:number) {
    const body = [{
        'id':id,
        'status':2
    }];
    return apiFetch<Todo>(`${baseUrl}tempora/todos/update`, {
        method: "POST",
        body: JSON.stringify(body),
        headers: {
            "Content-Type": "application/json",
        },
    });
}