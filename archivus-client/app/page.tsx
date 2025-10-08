'use client';
import Image from "next/image";
import { useAuth } from "./hooks/useAuth";
import FileExplorer from "./ui/fileexplorer_v2";
import { useSearchParams } from "next/navigation";
import { Suspense } from "react";

export default function Home() {
  useAuth();
  const searchParams = useSearchParams();
  const folder = searchParams.get('folder') || '';
  return (
    <>
    <Suspense fallback={<p>Loading...</p>}>
    <FileExplorer  />
    </Suspense>
    </>
  );
}
