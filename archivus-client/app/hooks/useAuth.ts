'use client';
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "../store/auth";


export function useAuth() {
    const router = useRouter();
    const user = useAuthStore.getState().user;
    useEffect(() => {
        if (!user) {
            router.push("/login");
        }
    },[])
}